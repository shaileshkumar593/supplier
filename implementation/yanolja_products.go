package implementation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	model "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/yanolja"
	yanoljasvc "swallow-supplier/services/suppliers/yanolja"
	"swallow-supplier/utils"

	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetProducts View a list of products available for sale on your channel.
func (s *service) GetProducts(ctx context.Context, req yanolja.AllProduct) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "GetProducts",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetProduct(ctx, req)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error", err)
		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), "GetProduct"), nil)
		} else {
			err = fmt.Errorf("request to yanolja client raised error")
		}
		return resp, err
	}

	level.Info(logger).Log("response", resp)

	return resp, nil

}

// GetProductsById get a products from yanolja
func (s *service) GetProductsById(ctx context.Context, req yanolja.ProductsById) (resp yanolja.Response, err error) {

	var requestID = utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "GetProductsById",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)
	productId, err := strconv.ParseInt(req.ProductId, 10, 64)
	if err != nil {
		return resp, customError.NewError(ctx, "Error converting ProductId to int64:", err.Error(), nil)
	}

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetProductByProductId(ctx, productId)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error", err)
		resp.Code = resp.Code
		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), "GetProduct"), nil)
		} else {
			err = fmt.Errorf("request to yanolja client raised error")
		}
		return resp, err
	}

	level.Info(logger).Log("response", resp)

	return resp, nil

}

// GetProductsOptionGroups a product oprtion group  from yanolja
func (s *service) GetProductsOptionGroups(ctx context.Context, req yanolja.ProductsById) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "GetProductsOptionGroups",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	productId, err := strconv.ParseInt(req.ProductId, 10, 64)
	if err != nil {
		return resp, customError.NewError(ctx, "Error converting ProductId to int64:", err.Error(), nil)
	}

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetProductOptionGroups(ctx, productId)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error", err)
		resp.Code = resp.Code
		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), "GetProduct"), nil)
		} else {
			err = fmt.Errorf("request to yanolja client raised error")
		}
		return resp, err
	}

	level.Info(logger).Log("response", resp)

	return resp, nil

}

// GetProductsInventories get a product inventory from yanolja
func (s *service) GetProductsInventories(ctx context.Context, req yanolja.ProductInventory) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "GetProductsInventories",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetProductsInventories(ctx, req)
	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error", err)
		resp.Code = resp.Code
		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), "GetProduct"), nil)
		} else {
			err = fmt.Errorf("request to yanolja client raised error")
		}
		return resp, err
	}

	level.Info(logger).Log("response", resp)

	return resp, nil

}

// GetVariantInventory get a product all variant Inventory from yanolja
func (s *service) GetVariantInventory(ctx context.Context, req yanolja.VariantInventory) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "GetVariantInventory",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.GetProductVariantInventory(ctx, req)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error", err)
		resp.Code = resp.Code
		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), "GetProduct"), nil)
		} else {
			err = fmt.Errorf("request to yanolja client raised error")
		}
		return resp, err
	}

	level.Info(logger).Log("response", resp)

	return resp, nil

}

// GetProductByProductId
func (s *service) GetProductByProductId(ctx context.Context, productId int64) (resp yanolja.Response, err error) {

	requestID := utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "GetProductByProductId",
		"Request ID", requestID,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Info(logger).Log("info", "processing request went into panic mode")

		resp.Code = "500"

	}(ctx)

	rec, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, productId)
	if err != nil {
		resp.Code = "500"
		resp.Body = "Error during inserting product to database"
		return resp, err
	}

	resp.Code = "200"
	resp.Body = rec
	return resp, nil
}

// InsertAllProduct
func (s *service) InsertAllProduct(ctx context.Context) (resp yanolja.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "InsertAllProduct",
		"Request ID", requestID,
	)
	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Error(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	var pageDetail yanolja.AllProduct
	pageDetail.PageNumber = 0
	pageDetail.PageSize = 10
	var productAry []interface{}
	var pagecnt int = 1

	// Service call to fetch product data
	getsvc, _ := yanoljasvc.New(ctx)
	records, err := getsvc.GetProduct(ctx, pageDetail)
	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error", err)
		resp.Code = records.Code
		if resp.Code == "403" {
			err = customError.NewError(ctx, "leisure-api-0003", fmt.Sprintf(customError.ErrForbiddenClient.Error(), "GetProduct"), nil)
		} else {
			err = fmt.Errorf("request to yanolja client raised error")
		}
		return resp, err
	}
	productAry = append(productAry, records.Body)

	pagecnt = pagecnt + 1
	for pagecnt <= records.Page.TotalPageCount {
		pageDetail.PageNumber = pageDetail.PageNumber + 1
		records, err = getsvc.GetProduct(ctx, pageDetail)
		if err != nil {
			level.Error(logger).Log("error", "request to yanolja client raise error for pagenumber %d ", pageDetail.PageNumber, err)
			resp.Code = "500"
			return resp, fmt.Errorf("failed to fetch product for pagecount %d: %w", pageDetail.PageNumber, err)
		}
		pagecnt = pagecnt + 1
		productAry = append(productAry, records.Body)
	}

	for _, val := range productAry {
		// Assuming records.Body is not a string but directly a map or a struct
		doc, err := json.Marshal(val)
		if err != nil {
			resp.Code = "500"
			return resp, fmt.Errorf("marshal error ")
		}
		// Unmarshal JSON response
		var rec []model.Product
		err = json.Unmarshal([]byte(doc), &rec)
		if err != nil {
			level.Error(logger).Log("error", "failed to unmarshal JSON response", err)
			resp.Code = "500"
			return resp, fmt.Errorf("json unmarshal error: %w", err)
		}

		for _, product := range rec {

			product.OodoSyncStatus = false
			// Insert products into MongoDB
			err = s.mongoRepository[config.Instance().MongoDBName].UpsertProduct(ctx, product)
			if err != nil {
				level.Error(logger).Log("error", "failed to insert products into MongoDB", err)
				resp.Code = "500"
				return resp, fmt.Errorf("mongo insert error: %w", err)
			}
		}
	}

	resp.Code = "200"
	resp.Body = fmt.Sprint("All products are recorded into database")
	return resp, nil
}
