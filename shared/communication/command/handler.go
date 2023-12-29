package command

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/types"
	"tradingplatform/shared/utils"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var commandHandler *utils.Handler[string]

// Initialize command handling
func StartCommandHandler(component types.Component, cliHandler func() *cobra.Command, jsonHandler func(string) string) {
	commandHandler = utils.NewHandler[string]()
	go handleCommand(commandHandler, component, cliHandler, jsonHandler)
}

func GetCommandHandler() *utils.Handler[string] {
	return commandHandler
}

func handleCommandContent(m *nats.Msg, handler *utils.Handler[string], cliHandler func() *cobra.Command, jsonHandler func(string) string) {
	cmd := string(m.Data)
	if len(cmd) > 4 && cmd[:4] == "json" {
		handleJSONCommand(m, jsonHandler)
	} else {
		handleCLICommand(m, handler, cliHandler)
	}

}

func handleJSONCommand(m *nats.Msg, jsonHandler func(string) string) {
	response := jsonHandler(string(m.Data[4:]))
	if len(response) == 0 {
		m.Respond([]byte(types.NewError(
			fmt.Errorf("no response from command due to invalid command, or component might have quit"),
		).Respond()))
		return
	}
	m.Respond([]byte(response))
}

func handleCLICommand(m *nats.Msg, handler *utils.Handler[string], cliHandler func() *cobra.Command) {
	str := strings.Split(string(m.Data), " ")
	rootCmd := cliHandler()
	rootCmd.SetArgs(str)

	buf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(errBuf)
	childContext, cancel := context.WithCancel(handler.Ctx())
	defer cancel()
	err := rootCmd.ExecuteContext(childContext)
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
			fmt.Errorf("no response from command due to invalid command or component might have quit"),
		).Respond()))
		return
	}
	m.Respond(buf.Bytes())
}

func handleCommand(handler *utils.Handler[string], component types.Component, cliHandler func() *cobra.Command, jsonHandler func(string) string) {
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
