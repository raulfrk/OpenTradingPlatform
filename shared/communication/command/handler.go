package command

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var commandHandler *utils.Handler[string]

var cancelKeys = make(map[string]context.CancelFunc)
var cancelKeyMutex sync.RWMutex

// CancelFunc is a function that can be used to get the cancel function for a command
func GetCancelFunc(key string) (context.CancelFunc, bool) {
	cancelKeyMutex.RLock()
	defer cancelKeyMutex.RUnlock()
	cancelFunc, ok := cancelKeys[key]
	return cancelFunc, ok
}

// AddCancelFunc adds a cancel function for a command, returns an error if the key already exists
func AddCancelFunc(key string, cancelFunc context.CancelFunc) error {
	cancelKeyMutex.Lock()
	defer cancelKeyMutex.Unlock()
	if _, ok := cancelKeys[key]; ok {
		return fmt.Errorf("cancel key %s already exists", key)
	}
	cancelKeys[key] = cancelFunc
	return nil
}

// RemoveCancelFunc removes a cancel function for a command
func RemoveCancelFunc(key string) {
	cancelKeyMutex.Lock()
	defer cancelKeyMutex.Unlock()
	delete(cancelKeys, key)
}

// Initialize command handling
func StartCommandHandler(component types.Component, cliHandler func() *cobra.Command, jsonHandler func(context.Context, string) string) {
	commandHandler = utils.NewHandler[string]()
	go handleCommand(commandHandler, component, cliHandler, jsonHandler)
}

// GetCommandHandler returns the command handler
func GetCommandHandler() *utils.Handler[string] {
	return commandHandler
}

func handleCommandContent(m *nats.Msg, handler *utils.Handler[string], cliHandler func() *cobra.Command, jsonHandler func(context.Context, string) string) {
	cmd := string(m.Data)
	if len(cmd) > 4 && cmd[:4] == "json" {
		handleJSONCommand(m, handler, jsonHandler)
	} else {
		handleCLICommand(m, handler, cliHandler)
	}

}

func handleJSONCommand(m *nats.Msg, handler *utils.Handler[string], jsonHandler func(context.Context, string) string) {
	childContext, cancel := context.WithCancel(handler.Ctx())
	defer cancel()
	response := jsonHandler(context.WithValue(childContext, CancelKey{}, cancel), string(m.Data[4:]))
	if len(response) == 0 {
		m.Respond([]byte(types.NewError(
			fmt.Errorf("no response provided, either component quit or command was cancelled"),
		).Respond()))
		return
	}
	m.Respond([]byte(response))
}

func splitArgs(input string) []string {
	// Split the string by spaces
	words := strings.Fields(input)

	var result []string
	var inQuote bool
	var currentWord string

	// Iterate through each word to handle quotes
	for _, word := range words {

		if strings.HasPrefix(word, `"`) {
			// Start of a double-quoted section
			inQuote = true
			currentWord = word[1:]
		} else if strings.HasSuffix(word, `"`) {
			// End of a double-quoted section
			inQuote = false
			currentWord += " " + strings.TrimSuffix(word, `"`)
			result = append(result, currentWord)
		} else if inQuote {
			// Inside a double-quoted section
			currentWord += " " + word
		} else {
			// Not inside a double-quoted section
			result = append(result, word)
		}
	}

	return result
}

type CancelKey struct{}

func handleCLICommand(m *nats.Msg, handler *utils.Handler[string], cliHandler func() *cobra.Command) {
	iString := string(m.Data)

	rootCmd := cliHandler()
	rootCmd.SetArgs(splitArgs(iString))

	buf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(errBuf)
	childContext, cancel := context.WithCancel(handler.Ctx())
	defer cancel()
	err := rootCmd.ExecuteContext(context.WithValue(childContext, CancelKey{}, cancel))
	// If cobra produces an error (e.g. unknown command)
	// we want to send it back to the caller
	if errBuf.Len() > 0 {
		logging.Log().Debug().Err(err)
		m.Respond([]byte(types.NewError(
			fmt.Errorf("%s\n%s", errBuf.String(), buf.String()),
		).Respond()))
		return
	}
	// If command produced no response, the command is considered to have failed
	if buf.Len() == 0 {
		m.Respond([]byte(types.NewError(
			fmt.Errorf("no response provided, either component quit or command was cancelled"),
		).Respond()))
		return
	}
	m.Respond(buf.Bytes())
}

func handleCommand(handler *utils.Handler[string], component types.Component, cliHandler func() *cobra.Command, jsonHandler func(context.Context, string) string) {
	nc, err := nats.Connect(nats.DefaultURL, nats.FlusherTimeout(0))
	ctx := handler.Ctx()
	if err != nil {
		logging.Log().Fatal().Err(err).Msg("NATS client could not connect to handle command")
	}
	defer nc.Close()
	commandTopic := utils.NewCommandTopic(component).Generate()
	logging.Log().Debug().Str("topic", commandTopic).Msg("subscribing to topic")
	nc.QueueSubscribe(commandTopic, "command", func(m *nats.Msg) {
		handler.Wg.Add(1)
		go func() {
			handleCommandContent(m, handler, cliHandler, jsonHandler)
			handler.Wg.Done()
		}()
	})
	<-ctx.Done()
	handler.Wg.Wait()
	logging.Log().Debug().Str("topic", commandTopic).Msg("unsubscribing from topic")
}
