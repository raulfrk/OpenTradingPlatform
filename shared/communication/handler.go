package communication

import (
	"google.golang.org/grpc"
	"net"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/utils"
)

var communicationHandler *utils.Handler[string]

func StartCommunicationHandler(server *grpc.Server) {
	if communicationHandler != nil {
		logging.Log().Warn().Msg("Communication handler already started")
		return
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logging.Log().Debug().Err(err).Msg("failed to listen")
	}
	communicationHandler = utils.NewHandler[string]()
	logging.Log().Debug().Msg("Starting communication handler")
	if err := server.Serve(lis); err != nil {
		logging.Log().Debug().Err(err).Msg("failed to serve")
	}
}

func GetCommunicationHandler() *utils.Handler[string] {
	return communicationHandler
}
