package travolution

import (
	"context"
	"fmt"
	"net/http"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetProducts get product based on the productId or get all products
func (t *Travolution) GetOptions(ctx context.Context, req travolution.OptionRequest) (res travolution.Response, err error) {
	logger := log.With(t.Service.Logger, "method", "GetOptions")

	var urlPath string

	// Handle URL based on presence of OptionUid
	switch v := req.OptionUid.(type) {
	case nil:
		// /products/:productUid/options
		urlPath = fmt.Sprintf("/api/partner/v1.1/products/%d/options/?lang=%s", req.ProductUid, req.Lang)

	case string:
		// /products/:productUid/options/:optionUid (string)
		urlPath = fmt.Sprintf("/api/partner/v1.1/products/%d/options/%s?lang=%s", req.ProductUid, v, req.Lang)

	case int, int64:
		// /products/:productUid/options/:optionUid (numeric)
		urlPath = fmt.Sprintf("/api/partner/v1.1/products/%d/options/%v?lang=%s", req.ProductUid, v, req.Lang)

	default:
		return res, fmt.Errorf("unsupported OptionUid type: %T", req.OptionUid)
	}

	// Append language query if present
	/* if req.Lang != "" {
		urlPath += fmt.Sprintf("?lang=%s", req.Lang)
	} */

	level.Info(logger).Log(
		"msg", "sending request to fetch product options",
		"url", urlPath,
	)

	response, err := t.Service.Send(
		ctx,
		ServiceName,
		t.Host+urlPath,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)

	return ResponseConvertor(ctx, response, logger, err)
}
