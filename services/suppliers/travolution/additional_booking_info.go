package travolution

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func (t *Travolution) GetBookingAdditionalInfo(ctx context.Context, req travolution.BookingAdditionalInfoRequest) (res travolution.Response, err error) {
	logger := log.With(t.Service.Logger, "method", "GetBookingAdditionalInfo")

	// --- Validate required productUid
	if req.ProductUID <= 0 {
		level.Error(logger).Log("error", "invalid ProductUid", "value", req.ProductUID)
		return res, errors.New("invalid ProductUid: must be greater than zero")
	}

	// --- Validate required optionUid
	if req.OptionUID == nil {
		level.Error(logger).Log("error", "missing OptionUid")
		return res, errors.New("optionUid is required")
	}

	// --- Convert OptionUid to string
	var optionUIDSegment string
	switch v := req.OptionUID.(type) {
	case string:
		optionUIDSegment = v
	case int, int64:
		optionUIDSegment = fmt.Sprintf("%v", v)
	default:
		level.Error(logger).Log("error", "unsupported OptionUid type", "type", fmt.Sprintf("%T", req.ProductUID))
		return res, fmt.Errorf("unsupported OptionUid type: %T", req.OptionUID)
	}

	// --- Optional additionalInfoUid
	var additionalInfoUIDSegment string
	if req.AdditionalInfoUID != nil {
		switch v := req.AdditionalInfoUID.(type) {
		case string:
			additionalInfoUIDSegment = "/" + v
		case int, int64:
			additionalInfoUIDSegment = fmt.Sprintf("/%v", v)
		default:
			level.Error(logger).Log("error", "unsupported AdditionalInfoUid type", "type", fmt.Sprintf("%T", req.AdditionalInfoUID))
			return res, fmt.Errorf("unsupported AdditionalInfoUid type: %T", req.AdditionalInfoUID)
		}
	}

	// --- Build final URL
	urlPath := fmt.Sprintf(
		"/api/partner/v1.1/products/%d/options/%s/booking-additional-info%s",
		req.ProductUID,
		optionUIDSegment,
		additionalInfoUIDSegment,
	)

	fullURL := fmt.Sprintf("%s%s", t.Host, urlPath)

	level.Info(logger).Log("msg", "fetching booking additional info", "url", fullURL)

	// --- Make HTTP GET request
	response, err := t.Service.Send(
		ctx,
		ServiceName,
		fullURL,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	return ResponseConvertor(ctx, response, logger, err)
}
