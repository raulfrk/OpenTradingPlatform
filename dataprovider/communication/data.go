package communication

import (
	"context"
	"tradingplatform/dataprovider/handler"
	"tradingplatform/shared/communication"
	"tradingplatform/shared/logging"
	"tradingplatform/shared/requests"
	"tradingplatform/shared/types"
)

func (*Server) GetData(ctx context.Context, req *communication.DataRequestExternal) (*communication.DataResponseExternal, error) {

	// Generate stream request from flags
	dataRequest, err := requests.NewDataRequestFromRaw(req.Source,
		req.AssetClass,
		req.Symbols,
		req.Operation,
		req.DataTypes,
		req.Account,
		req.StartTime,
		req.EndTime,
		req.TimeFrame,
		req.NoConfirm)
	logging.Log().Info().
		RawJSON("dataRequest", dataRequest.JSON()).
		Msg("receiving data request")
	if err != nil {
		cmd.Print(types.NewDataError(err).Respond())
		return nil, err
	}
	var och chan types.DataResponse = make(chan types.DataResponse)
	go handler.HandleDataRequest(dataRequest, och)

	select {
	case response := <-och:
		cmd.Print(response.Respond())
		return
	case <-cmd.Context().Done():
		return
	}
}
