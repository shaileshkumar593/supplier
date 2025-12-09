package trip

import (
	"context"
	"fmt"
	"net/http"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/trip"
	req_resp "swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ContentApiSyncToTrip sync ImageUrl for getting imageId
func (tp *Trip) ImageUrlSyncToTrip(ctx context.Context, imageUrlSyncReq trip.ImageSyncToTripRequest) (res yanolja.Response, err error) {

	logger := log.With(tp.Service.Logger, "method", "ContentApiSyncToTrip")

	level.Info(logger).Log("endpoint of trip :", tp.Host+"/out/content/sync")

	response, err := tp.Service.Send(
		ctx,
		ServiceName,
		tp.Host+"/out/content/sync",
		http.MethodPost,
		client.ContentTypeJSON,
		imageUrlSyncReq,
	)
	if err != nil {
		level.Error(logger).Log("error ", "trip service error")
		res.Body = response.Body
		res.Code = string(response.Status)
		return res, customError.NewError(ctx, "leisure-api-1022", fmt.Sprintf("trip request error %v", err), "ContentApiSyncToTrip")
	}

	return TripResponseConversion(tp.Ctx, response, logger, err)
}

// ProductContentApiSyncToTrip  sync the productContent to trip
func (tp *Trip) ProductContentApiSyncToTrip(ctx context.Context, productContents req_resp.ProductContentSync) (res yanolja.Response, err error) {

	logger := log.With(tp.Service.Logger, "method", "ProductContentApiSyncToTrip")

	level.Info(logger).Log("endpoint of trip :", tp.Host+"/out/content/sync")

	response, err := tp.Service.Send(
		ctx,
		ServiceName,
		tp.Host+"/out/content/sync",
		http.MethodPost,
		client.ContentTypeJSON,
		productContents,
	)
	if err != nil {
		level.Error(logger).Log("error ", "trip service error")
		res.Body = response.Body
		res.Code = string(response.Status)
		return res, customError.NewError(ctx, "leisure-api-1022", fmt.Sprintf("trip request error %v", err), "ProductContentApiSyncToTrip")
	}

	return TripResponseConversion(tp.Ctx, response, logger, err)
}

// PackageContentApiSyncToTrip  sync the productContent to trip
func (tp *Trip) PackageContentApiSyncToTrip(ctx context.Context, productContents req_resp.PackageContentSync) (res yanolja.Response, err error) {

	logger := log.With(tp.Service.Logger, "method", "PackageContentApiSyncToTrip")

	level.Info(logger).Log("endpoint of trip :", tp.Host+"/out/content/sync")

	response, err := tp.Service.Send(
		ctx,
		ServiceName,
		tp.Host+"/out/content/sync",
		http.MethodPost,
		client.ContentTypeJSON,
		productContents,
	)
	if err != nil {
		level.Error(logger).Log("error ", "trip service error")
		res.Body = response.Body
		res.Code = string(response.Status)
		return res, customError.NewError(ctx, "leisure-api-1022", fmt.Sprintf("trip request error %v", err), "PackageContentApiSyncToTrip")
	}

	return TripResponseConversion(tp.Ctx, response, logger, err)
}
