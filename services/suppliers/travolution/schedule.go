package travolution

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func (t *Travolution) GetBookingSchedules(ctx context.Context, req travolution.BookingScheduleReq) (res travolution.Response, err error) {
	logger := log.With(t.Service.Logger, "method", "GetBookingSchedules")

	// Validate required path fields
	if req.ProductUid <= 0 {
		level.Error(logger).Log("error", "invalid ProductUid", "value", req.ProductUid)
		return res, fmt.Errorf("invalid ProductUid: must be > 0")
	}
	if req.OptionUid == nil {
		level.Error(logger).Log("error", "missing OptionUid")
		return res, errors.New("OptionUid is required")
	}

	// Convert OptionUid to string
	var optionUIDSegment string
	switch v := req.OptionUid.(type) {
	case string:
		optionUIDSegment = v
	case int, int64:
		optionUIDSegment = fmt.Sprintf("%v", v)
	default:
		level.Error(logger).Log("error", "unsupported OptionUid type", "type", fmt.Sprintf("%T", req.OptionUid))
		return res, fmt.Errorf("unsupported OptionUid type: %T", req.OptionUid)
	}

	// Construct URL path
	urlPath := fmt.Sprintf("/api/partner/v1.1/products/%d/options/%s/booking-schedules", req.ProductUid, optionUIDSegment)

	// Build query parameters (both optional)
	query := url.Values{}
	if strings.TrimSpace(req.Date) != "" {
		query.Set("date", req.Date)
	}
	if strings.TrimSpace(req.Time) != "" {
		query.Set("time", req.Time)
	}

	// Compose full URL
	fullURL := t.Host + urlPath
	if encodedQuery := query.Encode(); encodedQuery != "" {
		fullURL += "?" + encodedQuery
	}

	level.Info(logger).Log("msg", "fetching booking schedules", "url", fullURL)

	// Send HTTP GET request
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
