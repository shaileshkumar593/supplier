package yanolja

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils/client"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// GetProduct get all product from yanolja
func (y *Yanolja) GetProduct(ctx context.Context, req yanolja.AllProduct) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "GetProduct")
	level.Info(logger).Log("info", "Service GetProduct")
	var param string
	if req.ProductStatusCode == "" {
		param = fmt.Sprintf("?pageNumber=%d&pageSize=%d", req.PageNumber, req.PageSize)
	} else {
		param = fmt.Sprintf("?pageNumber=%d&pageSize=%d&productStatusCode=%s", req.PageNumber, req.PageSize, req.ProductStatusCode)
	}
	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+"/v1/products"+param,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	res.Code = strconv.Itoa(response.Status)
	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "empty response", nil)
	}

	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}
	return res, nil
}

// GetProductByProductId get product by productId
func (y *Yanolja) GetProductByProductId(ctx context.Context, productid int64) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "GetProduct")
	level.Info(logger).Log("info", "Service GetProduct")

	urls := fmt.Sprintf(`/v1/products/%d`, productid)
	level.Info(logger).Log("info", "url for get request ", urls)

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "empty response", nil)
	}

	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}
	return res, nil
}

// GetProductOptionGroups get the  ProductOptionsGroup of a product
func (y *Yanolja) GetProductOptionGroups(ctx context.Context, productid int64) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "GetProduct")
	level.Info(logger).Log("info", "Service GetProduct")

	urls := fmt.Sprintf(`/v1/products/%d/option-groups`, productid)
	level.Info(logger).Log("info", "url for get request ", urls)

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "empty response", nil)
	}

	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}
	return res, nil
}

// GetProductsInventories get the  product inventory of a product
func (y *Yanolja) GetProductsInventories(ctx context.Context, req yanolja.ProductInventory) (res yanolja.Response, err error) {
	var urls string

	logger := log.With(y.Service.Logger, "method", "GetProductsInventories")
	level.Info(logger).Log("info", "Service GetProductsInventories")

	productId, err := strconv.ParseInt(req.ProductId, 10, 64)
	if err != nil {
		res.Code = "500"
		return res, customError.NewError(ctx, "Error converting ProductId to int64:", err.Error(), nil)
	}

	if req.InventoryDateStart != "" && req.InventoryDateEnd != "" {
		urls = fmt.Sprintf(`/v1/products/%d/inventories?inventoryDateStart=%s&inventoryDateEnd=%s`, productId, req.InventoryDateStart, req.InventoryDateEnd)
	} else {
		urls = fmt.Sprintf(`/v1/products/%d/inventories`, productId)
	}

	level.Info(logger).Log("info", "url for inventory request ", urls)
	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "empty response", nil)
	}

	//level.Info(logger).Log("response body from yanolja ", response.Body)
	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	return res, nil
}

// GetProductVariantInventory get product variant inventory based on a variantId
func (y *Yanolja) GetProductVariantInventory(ctx context.Context, req yanolja.VariantInventory) (res yanolja.Response, err error) {
	logger := log.With(y.Service.Logger, "method", "GetProduct")
	level.Info(logger).Log("info", "Service GetProduct")

	variantid, err := strconv.ParseInt(req.VariantId, 10, 64)
	if err != nil {
		return res, customError.NewError(ctx, "Error converting ProductId to int64:", err.Error(), nil)
	}

	urls := fmt.Sprintf(`/v1/products/-/variants/%d/inventory`, variantid)
	level.Info(logger).Log("info", "url for get request ", urls)

	response, err := y.Service.Send(
		ctx,
		ServiceName,
		y.Host+urls,
		http.MethodGet,
		client.ContentTypeJSON,
		nil,
	)
	res.Code = strconv.Itoa(response.Status)

	if err == nil && response == nil && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "empty response", nil)
	}

	r := response.Body
	if r == "" && response.Status != http.StatusOK {
		response, _ := json.Marshal(response)
		level.Info(logger).Log("err", string(response))
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}

	// Convert the string to a byte array
	bodyBytes := []byte(response.Body)

	// Unmarshal the byte array into the Response struct
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		// handle error
		return res, customError.NewError(y.Ctx, "external_processing_error", "Empty Body", nil)
	}
	return res, nil
}
