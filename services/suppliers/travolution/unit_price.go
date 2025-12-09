package travolution

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func (t *Travolution) GetUnits(ctx context.Context, req travolution.UnitPriceRequest) (res travolution.Response, err error) {
	logger := log.With(t.Service.Logger, "method", "GetUnits")

	// Validate input
	if req.ProductUid <= 0 {
		level.Error(logger).Log("error", "invalid ProductUid: must be > 0")
		return res, fmt.Errorf("invalid ProductUid: must be > 0")
	}

	if req.OptionUid == nil {
		level.Error(logger).Log("error", "OptionUid is required")
		return res, fmt.Errorf("OptionUid is required")
	}

	// Format OptionUid
	var optionUIDSegment string
	switch v := req.OptionUid.(type) {
	case string:
		optionUIDSegment = v
	case int:
		optionUIDSegment = strconv.Itoa(v)
	case int64:
		optionUIDSegment = strconv.FormatInt(v, 10)
	default:
		level.Error(logger).Log("error", "unsupported OptionUid type", "type", fmt.Sprintf("%T", req.OptionUid))
		return res, fmt.Errorf("unsupported OptionUid type: %T", req.OptionUid)
	}

	// Build base URL path
	urlPath := fmt.Sprintf("api/partner/v1.1/products/%d/options/%s/units", req.ProductUid, optionUIDSegment)

	// Format UnitUid if present
	if req.UnitUid != nil {
		switch v := req.UnitUid.(type) {
		case string:
			urlPath = fmt.Sprintf("%s/%s", urlPath, v)
		case int:
			urlPath = fmt.Sprintf("%s/%d", urlPath, v)
		case int64:
			urlPath = fmt.Sprintf("%s/%d", urlPath, v)
		default:
			level.Error(logger).Log("error", "unsupported UnitUid type", "type", fmt.Sprintf("%T", req.UnitUid))
			return res, fmt.Errorf("unsupported UnitUid type: %T", req.UnitUid)
		}
	} else {
		urlPath = fmt.Sprintf("%s/", urlPath)
	}

	fullURL := t.Host + "/" + urlPath
	level.Info(logger).Log("msg", "Sending GET request to fetch unit(s)", "url", fullURL)

	// Perform GET request
	response, err := t.Service.Send(
		ctx,
		ServiceName,
		fullURL,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	if err != nil {
		level.Error(logger).Log("msg", "failed to send request", "error", err)
		return res, err
	}

	return ResponseConvertor(ctx, response, logger, err)
}
