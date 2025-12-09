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
func (t *Travolution) GetProducts(ctx context.Context, req travolution.ProductReq) (res travolution.Response, err error) {

	logger := log.With(t.Service.Logger, "method", "GetProductsOfTravolution")

	var urls string
	if req.ProductUid == 0 && req.Take == 0 && req.Skip == 0 && req.Lang == "" {
		urls = `/api/partner/v1.1/products/`

	} else if req.ProductUid == 0 && req.Take > 0 && req.Skip == 0 && req.Lang == "" {
		urls = fmt.Sprintf(`/api/partner/v1.1/products/?take=%d`, req.Take)

	} else if req.ProductUid == 0 && req.Take > 0 && req.Skip > 0 && req.Lang == "" {
		urls = fmt.Sprintf(`/api/partner/v1.1/products/?take=%d&skip=%d`, req.Take, req.Skip)

	} else if req.ProductUid == 0 && req.Take > 0 && req.Skip > 0 && req.Lang != "" {
		urls = fmt.Sprintf(`/api/partner/v1.1/products/?take=%d&skip=%d&lang=%s`, req.Take, req.Skip, req.Lang)

	} else if req.ProductUid > 0 && req.Take == 0 && req.Skip == 0 && req.Lang != "" {
		urls = fmt.Sprintf(`/api/partner/v1.1/products/%d?lang=%s`, req.ProductUid, req.Lang)

	} else if req.ProductUid > 0 && req.Take == 0 && req.Skip == 0 && req.Lang == "" {
		urls = fmt.Sprintf(`/api/partner/v1.1/products/%d`, req.ProductUid)

	}

	level.Info(logger).Log(
		"info", "url in request ",
		"url", urls,
	)

	response, err := t.Service.Send(
		ctx,
		ServiceName,
		t.Host+urls,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)

	return ResponseConvertor(ctx, response, logger, err)
}
