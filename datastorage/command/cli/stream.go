package cli

import (
	"tradingplatform/datastorage/handler"
	"tradingplatform/datastorage/subscriber"
	"tradingplatform/shared/requests"

	"tradingplatform/shared/types"

	"github.com/spf13/cobra"
)

func NewStreamCommand() *cobra.Command {
	streamAddCmd := cobra.Command{
		Use:   "stream",
		Short: "Stream management for DataStorage",
	}

	streamAddCmd.AddCommand(NewStreamAddCommand())
	streamAddCmd.AddCommand(NewStreamDeleteCommand())

	return &streamAddCmd
}

func NewStreamDeleteCommand() *cobra.Command {
	streamAddCmd := cobra.Command{
		Use:   "delete",
		Short: "Deletes a stream from the DataStorage (the data on stream will not be stored anymore)",
		Run: func(cmd *cobra.Command, args []string) {
			topic, _ := cmd.Flags().GetStringArray("topic")
			req, err := requests.NewStreamSubscribeRequestFromRaw(topic, types.StreamRemoveOp)
			if err != nil {
				cmd.Print(subscriber.NewStreamErrorResponseTopic(err).Respond())
				return
			}
			res := handler.HandleStreamRequest(req)

			cmd.Print(res.Respond())
		},
	}

	streamAddCmd.Flags().StringArrayP("topic", "t", []string{}, "Topics to delete")

	streamAddCmd.MarkFlagRequired("topic")

	return &streamAddCmd
}

func NewStreamAddCommand() *cobra.Command {
	streamAddCmd := cobra.Command{
		Use:   "add",
		Short: "Adds a stream to the DataStorage",
		Run: func(cmd *cobra.Command, args []string) {
			topicAgents, _ := cmd.Flags().GetStringArray("topic-agents")
			req, err := requests.NewStreamSubscribeRequestFromRaw(topicAgents, types.StreamAddOp)
			if err != nil {
				cmd.Print(subscriber.NewStreamErrorResponseTopic(err).Respond())
				return
			}
			res := handler.HandleStreamRequest(req)

			cmd.Print(res.Respond())
		},
	}

	streamAddCmd.Flags().StringArrayP("topic-agents", "t", []string{},
		"Comma separated pair with topic and agents to assign to it. If number of agents is not specified, 5 will be used.")

	streamAddCmd.MarkFlagRequired("topic-agents")

	return &streamAddCmd
}
