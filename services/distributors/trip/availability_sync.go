package trip

import (
	"context"
	"fmt"
	"net/http"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"

	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetOrderbyId get order based on the orderId
func (tp *Trip) InventorySyncToTrip(ctx context.Context, inventorySyncReq yanolja.InventoryToTrip) (res yanolja.Response, err error) {

	logger := log.With(tp.Service.Logger, "method", "InventorySyncToTrip")

	level.Info(logger).Log("***************payload ****************** ", inventorySyncReq)
	level.Info(logger).Log("endpoint of trip :", tp.Host+"/out/availability/sync")

	response, err := tp.Service.Send(
		ctx,
		ServiceName,
		tp.Host+"/out/availability/sync",
		http.MethodPost,
		client.ContentTypeJSON,
		inventorySyncReq,
	)
	if err != nil {
		level.Error(logger).Log("error ", "trip service error")
		res.Body = response.Body
		res.Code = string(response.Status)
		return res, customError.NewError(ctx, "leisure-api-1022", fmt.Sprintf("trip request error %v", err), "NotifyToTrip")
	}

	return TripResponseConversion(tp.Ctx, response, logger, err)
}
