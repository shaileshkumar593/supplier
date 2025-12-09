package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	customContext "swallow-supplier/context"
	customError "swallow-supplier/error"
	svc "swallow-supplier/iface"
	"swallow-supplier/middleware"
	travolution_domain "swallow-supplier/mongo/domain/travolution"
	domain "swallow-supplier/mongo/domain/yanolja"

	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/transport"
	"swallow-supplier/utils"
)

// NewTransport set-up the router and initialize the http endpoints
func NewTransport(
	ctx context.Context,
	svcEndpoints transport.Endpoints,
	mongoRepo map[string]svc.MongoRepository,
) http.Handler {
	var (
		router = mux.NewRouter()

		// server options:
		errorEncoder = kithttp.ServerErrorEncoder(customError.EncodeErrorResponse)

		options = make([]kithttp.ServerOption, 0)
	)

	options = append(options,
		kithttp.ServerBefore(
			customContext.RequestURLExtractor,
			customContext.RequestPathExtractor,
			customContext.RequestPathTemplateExtractor,
			customContext.TraceIDExtractor,
			customContext.RequestIDHeaderExtractor,
			customContext.ChannelCodeHeaderExtractor,
			customContext.ApiKeyHeaderExtractor,
			customContext.ForwardedForHeaderExtractor,
			customContext.RemoteAddrExtractor,
			customContext.GGTAuthorizationExtractor,
			kithttp.PopulateRequestContext,
		),
		errorEncoder,
		kithttp.ServerAfter(
			customContext.TraceIDSetter,
		),
	)
	// list all endpoints handler
	heartbeatHandler := kithttp.NewServer(
		svcEndpoints.HeartBeat,
		decodeHeartbeatRequest,
		encodeResponse,
		options...,
	)

	// heartbeat endpoint
	router.Handle("/api/v1/heartbeat", heartbeatHandler).Methods("GET")

	// register call to yanolja API
	getYanoljaProducts := kithttp.NewServer(
		svcEndpoints.GetProducts,
		decodeGetProductsRequest,
		encodeResponse,
		options...,
	)

	getYanoljaProductsById := kithttp.NewServer(
		svcEndpoints.GetProductsById,
		decodeGetProductsByIdRequest,
		encodeResponse,
		options...,
	)

	getYanoljaProductsOptionGroups := kithttp.NewServer(
		svcEndpoints.GetProductsOptionGroups,
		decodeGetProductsOptionGroups,
		encodeResponse,
		options...,
	)

	getYanoljaProductsInventories := kithttp.NewServer(
		svcEndpoints.GetProductsInventories,
		decodeGetProductsInventories,
		encodeResponse,
		options...,
	)

	getYanoljaVariantInventory := kithttp.NewServer(
		svcEndpoints.GetVariantInventory,
		decodeGetVariantInventory,
		encodeResponse,
		options...,
	)

	getYanoljaCategories := kithttp.NewServer(
		svcEndpoints.GetCategories,
		decodeGetCategories,
		encodeResponse,
		options...,
	)

	getYanoljaRegions := kithttp.NewServer(
		svcEndpoints.GetRegions,
		decodeGetRegions,
		encodeResponse,
		options...,
	)

	postWaitForOrder := kithttp.NewServer(
		svcEndpoints.PostWaitForOrder,
		decodePostWaitForOrder,
		encodeResponse,
		options...,
	)

	postOrderCompletion := kithttp.NewServer(
		svcEndpoints.PostOrderConfirmation,
		decodePostOrderCompletion,
		encodeResponse,
		options...,
	)

	getOrderByOrderId := kithttp.NewServer(
		svcEndpoints.GetOrderByOrderId,
		decodeGetOrderByOrderId,
		encodeResponse,
		options...,
	)

	getOrderByPartnerOrderIdSuffix := kithttp.NewServer(
		svcEndpoints.GetOrderByPartnerOrderIdSuffix,
		decodeGetOrderByPartnerOrderIdSuffix,
		encodeCommonResponse,
		options...,
	)

	getEverLandOrders := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.GetEverlandOrders),
		decodeGetEverlandOrders,
		encodeCommonResponse,
		options...,
	)

	postFullOrderCancellation := kithttp.NewServer(
		svcEndpoints.PostFullOrderCancel,
		decodePostFullOrderCancellation,
		encodeResponse,
		options...,
	)

	postElFullOrderCancellation := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.PostFullOrderCancel),
		decodePostFullOrderCancellation,
		encodeResponse,
		options...,
	)

	postReqTimeoutOrderCancel := kithttp.NewServer(
		svcEndpoints.PostCancelOrderByTimeOut,
		decodePostReqTimeoutOrderCancel,
		encodeResponse,
		options...,
	)

	postForcedOrderCancellation := kithttp.NewServer(
		svcEndpoints.PostCancelOrderByForcely,
		decodePostForcedOrderCancellation,
		encodeResponse,
		options...,
	)

	getProductByProductIdHandler := kithttp.NewServer(
		svcEndpoints.GetProductByProductId,
		decodeGetProductByProductId,
		encodeResponse,
		options...,
	)

	postAllProductHandler := kithttp.NewServer(
		svcEndpoints.PostAllProducts,
		decodePostAllProducts,
		encodeResponse,
		options...,
	)

	postYanoljaCategories := kithttp.NewServer(
		svcEndpoints.PostCategories,
		decodePostCategories,
		encodeResponse,
		options...,
	)

	postYanoljaRegions := kithttp.NewServer(
		svcEndpoints.PostRegions,
		decodePostRegions,
		encodeResponse,
		options...,
	)
	postAllGgtTripYanoljaMapping := kithttp.NewServer(
		svcEndpoints.UpsertCategoryMapping,
		decodeUpsertCategoryMapping,
		encodeResponse,
		options...,
	)

	// register all handlers into an http routes
	router.Handle("/v1/products", getYanoljaProducts).Methods("GET")
	router.Handle("/v1/products/{productid}", getYanoljaProductsById).Methods("GET")
	router.Handle("/v1/products/{productid}/option-groups", getYanoljaProductsOptionGroups).Methods("GET")
	router.Handle("/v1/products/{productid}/inventories", getYanoljaProductsInventories).Methods("GET")
	router.Handle("/v1/products/-/variants/{variantid}/inventory", getYanoljaVariantInventory).Methods("GET")
	router.Handle("/v1/categories", getYanoljaCategories).Methods("GET")
	router.Handle("/v1/categories", postYanoljaCategories).Methods("POST")
	router.Handle("/v1/regions", getYanoljaRegions).Methods("GET")
	router.Handle("/v1/regions", postYanoljaRegions).Methods("POST")
	router.Handle("/v1/orders/prepare", postWaitForOrder).Methods("POST")
	router.Handle("/v1/orders/confirm", postOrderCompletion).Methods("POST")
	router.Handle("/v1/orders/{orderId}", getOrderByOrderId).Methods("GET")
	router.Handle("/v1/order/get/{partnerOrderId}/suffix", getOrderByPartnerOrderIdSuffix).Methods("GET")
	router.Handle("/v1/orders/{orderId}/full-cancel", postFullOrderCancellation).Methods("POST")
	router.Handle("/v1/everland/orders/{orderId}/full-cancel", postElFullOrderCancellation).Methods("POST") // for everland
	router.Handle("/v1/orders/timeout-cancel", postReqTimeoutOrderCancel).Methods("POST")
	router.Handle("/v1/orders/force-cancel", postForcedOrderCancellation).Methods("POST")
	router.Handle("/v1/product/{productId}/get", getProductByProductIdHandler).Methods("GET")
	router.Handle("/v1/product/all", postAllProductHandler).Methods("POST")
	router.Handle("/v1/category/mapping", postAllGgtTripYanoljaMapping).Methods("POST")
	router.Handle("/v1/orders/customeremail/{customeremail}/channelcode/{channelcode}", getEverLandOrders).Methods("GET")

	getOrderReconcilationDetailClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.GetOrderReconcilationDetail),
		decodeGetOrderReconcilationDetail,
		encodeResponse,
		options...,
	)

	postCancellationAckClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.CancellationAckClbk),
		decodepostCancellationAckClbk,
		encodeResponse,
		options...,
	)

	postRefusalToCancelClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.RefusalToCancelClbk),
		decodePostRefusalToCancelClbk,
		encodeResponse,
		options...,
	)

	getOrderStatusLookupClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.GetOrderStatusLookupClbk),
		decodeGetOrderStatusLookupClbk,
		encodeResponse,
		options...,
	)

	postForcedOrderCancellationClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.ForcedOrderCancellationClbk),
		decodePostForcedOrderCancellationClbk,
		encodeResponse,
		options...,
	)

	postIndividualVoucherUpdateClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.IndividualVoucherUpdateClbk),
		decodePostIndividualVoucherUpdateClbk,
		encodeResponse,
		options...,
	)
	postCombinedVoucherUpdateClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.CombinedVoucherUpdateClbk),
		decodePostCombinedVoucherUpdateClbk,
		encodeResponse,
		options...,
	)
	postProcessingOrRestoringClbkHandler := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.ProcessingOrRestoringClbk),
		decodePostProcessingOrRestoringClbk,
		encodeResponse,
		options...,
	)

	postProductCreationInGGT := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.PostProductCreation),
		decodePostProductCreationInGGT,
		encodeResponse,
		options...,
	)

	router.Handle("/v1/order/reconciliation/callback", getOrderReconcilationDetailClbkHandler).Methods("GET")
	router.Handle("/v1/order/cancellation/ack/callback", postCancellationAckClbkHandler).Methods("POST")
	router.Handle("/v1/order/refusalto/cancel/callback", postRefusalToCancelClbkHandler).Methods("POST")
	router.Handle("/v1/order/status/lookup/callback/{partnerOrderId}", getOrderStatusLookupClbkHandler).Methods("GET")
	router.Handle("/v1/force/order/cancellation/callback", postForcedOrderCancellationClbkHandler).Methods("POST")
	router.Handle("/v1/order/individual/voucher/update/callback", postIndividualVoucherUpdateClbkHandler).Methods("POST")
	router.Handle("/v1/order/combined/voucher/update/callback", postCombinedVoucherUpdateClbkHandler).Methods("POST")
	router.Handle("/v1/order/processing/restoring/callback", postProcessingOrRestoringClbkHandler).Methods("POST")
	router.Handle("/v1/product/create/callback", postProductCreationInGGT).Methods("POST")

	postRequestFromGGT := kithttp.NewServer(
		svcEndpoints.PostRequestFromGGT,
		decodePostRequestFromGGT,
		encodeResponse,
		options...,
	)

	getRedisDataHandler := kithttp.NewServer(
		svcEndpoints.GetRedisData,
		decodeGetRedisData,
		encodeResponse,
		options...,
	)

	getPluToRedisHandler := kithttp.NewServer(
		svcEndpoints.GetPluToRedis,
		decodeGetPluToRedis,
		encodeResponse,
		options...,
	)
	deleteIfNotEmpty := kithttp.NewServer(
		svcEndpoints.DeleteIfNotEmpty,
		decodeDeleteIfNotEmpty,
		encodeResponse,
		options...,
	)
	postMonitoProductUpdate := kithttp.NewServer(
		svcEndpoints.MonitorProductUpdate,
		decodePostMonitoProductUpdate,
		encodeResponse,
		options...,
	)
	postPluUpsertToRedis := kithttp.NewServer(
		svcEndpoints.PluUpsertToRedis,
		decodePostPluUpsertToRedis,
		encodeResponse,
		options...,
	)

	DeleteKeysFromRedis := kithttp.NewServer(
		svcEndpoints.DeleteKeyFromRedis,
		decodeDeleteRedisKeys,
		encodeResponse,
		options...,
	)

	getRedisKeyValue := kithttp.NewServer(
		svcEndpoints.FindRedisKeyVal,
		decodeFindKeysVal,
		encodeResponse,
		options...,
	)

	updatePluToRedis := kithttp.NewServer(
		svcEndpoints.UpdatePluToRedis,
		decodeUpdatePluToRedis,
		encodeResponse,
		options...,
	)

	getPluFromRedis := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.FindPluFromRedis),
		decodeFindPluFromRedis,
		encodeResponse,
		options...,
	)

	router.Handle("/v1/ggt/order", postRequestFromGGT).Methods("POST")
	router.Handle("/v1/ggt/redis/data", getRedisDataHandler).Methods("GET")
	router.Handle("/v1/ggt/plu/redis/data", getPluToRedisHandler).Methods("GET")
	router.Handle("/v1/delete/data", deleteIfNotEmpty).Methods("DELETE")
	router.Handle("/v1/monitor/product/update", postMonitoProductUpdate).Methods("POST")
	router.Handle("/v1/plu/upsert/to/redis", postPluUpsertToRedis).Methods("POST")
	router.Handle("/v1/plu/delete/keys", DeleteKeysFromRedis).Methods("DELETE")
	router.Handle("/v1/redis/get/{key}/value", getRedisKeyValue).Methods("GET")
	router.Handle("/v1/update/plu", updatePluToRedis).Methods("POST")
	router.Handle("/v1/get/plu/{productId}/productId/{variantId}/variantId/{productVersion}/productVersion", getPluFromRedis).Methods("GET")

	getOrderForOdoo := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.GetOrderForOdoo),
		decodeGetOrderForOdoo,
		encodeCommonResponse,
		options...,
	)

	getProductForOdoo := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.GetProductViewForOdoo),
		decodeGetProductViewForOdoo,
		encodeCommonResponse,
		options...,
	)

	getOrderInOdooByOrderId := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.GetOdooByOrderId),
		decodeGetOdooByOrderId,
		encodeCommonResponse,
		options...,
	)

	postOdooOrderSyncFlag := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.PostUpdateOdooSyncFlag),
		decodePostUpdateOdooSyncFlag,
		encodeCommonResponse,
		options...,
	)

	postOdooSyncStatusFlagFalse := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.PostUpdateOdooSyncFlagFalse),
		decodePostUpdateOdooSyncFlagFalse,
		encodeCommonResponse,
		options...,
	)

	getGGTInventorySync := kithttp.NewServer(
		svcEndpoints.GetGGTInventorySync,
		decodeGetGGTInventorySync,
		encodeResponse,
		options...,
	)

	putUpdateImageSyncStatusHandle := kithttp.NewServer(
		svcEndpoints.PostUpdateImageSyncStatus,
		decodepostUpdateImageSyncStatus,
		encodeResponse,
		options...,
	)

	getImageSyncToTrip := kithttp.NewServer(
		svcEndpoints.GetImageSyncToTrip,
		decodeGetImageSyncToTrip,
		encodeResponse,
		options...,
	)

	putImageSyncStatus := kithttp.NewServer(
		svcEndpoints.UpdateTripImageSyncStatus,
		decodeUpdateTripImageSyncStatus,
		encodeResponse,
		options...,
	)

	getProductContentData := kithttp.NewServer(
		svcEndpoints.GetProductContentData,
		decodeGetProductContentData,
		encodeResponse,
		options...,
	)

	getPackageContentData := kithttp.NewServer(
		svcEndpoints.GetPackageContentData,
		decodeGetPackageContentData,
		encodeResponse,
		options...,
	)

	putContentSyncStatus := kithttp.NewServer(
		svcEndpoints.PutContentSyncData,
		decodePutContentSyncData,
		encodeResponse,
		options...,
	)

	postContent := kithttp.NewServer(
		svcEndpoints.PostContent,
		decodePutContent,
		encodeResponse,
		options...,
	)

	getProducts := kithttp.NewServer(
		svcEndpoints.GetProductsFromGGT,
		decodeGetProductsFromGGT,
		encodeResponse,
		options...,
	)

	getProductByProductId := kithttp.NewServer(
		svcEndpoints.GetProductByIdFromGGT,
		decodeGetProductByIdFromGGT,
		encodeResponse,
		options...,
	)

	// odoo sync api
	router.Handle("/v1/odoo/order/sync", getOrderForOdoo).Methods("GET")
	router.Handle("/v1/odoo/product/sync", getProductForOdoo).Methods("GET")
	router.Handle("/v1/odoo/order/{orderId}/sync", getOrderInOdooByOrderId).Methods("GET")
	router.Handle("/v1/odoo/order/syncflag", postOdooOrderSyncFlag).Methods("POST")
	router.Handle("/v1/odoo/order/sync/flag/update", postOdooSyncStatusFlagFalse).Methods("POST")

	//All sync related api
	router.Handle("/v1/ggt/inventory/sync", getGGTInventorySync).Methods("GET")
	router.Handle("/v1/ggt/content/api", putUpdateImageSyncStatusHandle).Methods("PUT")
	router.Handle("/v1/ggt/image/sync", getImageSyncToTrip).Methods("GET")          // POST
	router.Handle("/v1/ggt/image/status/update", putImageSyncStatus).Methods("PUT") //-------------need to remove
	router.Handle("/v1/ggt/product/content/sync", getProductContentData).Methods("GET")
	router.Handle("/v1/ggt/package/content/sync", getPackageContentData).Methods("GET")
	router.Handle("/v1/ggt/update/content/sync/status", putContentSyncStatus).Methods("PUT")
	router.Handle("/v1/ggt/get/products", getProducts).Methods("GET")
	router.Handle("/v1/ggt/get/product/{productId}", getProductByProductId).Methods("GET")

	router.Handle("/v1/ggt/sync/content", postContent).Methods("POST")

	//*********************** Travolution  *************************************************

	getAllproducts := kithttp.NewServer(
		svcEndpoints.GetAllproducts,
		decodeGetProducts,
		encodeTravolutionResponse,
		options...,
	)

	getProductByUid := kithttp.NewServer(
		svcEndpoints.GetProductByUid,
		decodeGetProducts,
		encodeTravolutionResponse,
		options...,
	)

	getAllOptionsOfProduct := kithttp.NewServer(
		svcEndpoints.GetAllOptionsOfProduct,
		decodeOption,
		encodeTravolutionResponse,
		options...,
	)

	getOptionOfProductByOptionUid := kithttp.NewServer(
		svcEndpoints.GetOptionOfProductByOptionUid,
		decodeOption,
		encodeTravolutionResponse,
		options...,
	)

	getUnitSPriceByOptionUid := kithttp.NewServer(
		svcEndpoints.GetUnitSPriceByOptionUid,
		decodeGetUnitDetails,
		encodeTravolutionResponse,
		options...,
	)

	getUnitPriceByOptionUidAndUnitId := kithttp.NewServer(
		svcEndpoints.GetUnitPriceByOptionUid,
		decodeGetUnitDetails,
		encodeTravolutionResponse,
		options...,
	)

	getBookingSchedules := kithttp.NewServer(
		svcEndpoints.GetBookingSchedules,
		decodeBookingSchedule,
		encodeTravolutionResponse,
		options...,
	)

	getAdditionalInfos := kithttp.NewServer(
		svcEndpoints.GetAdditionalInfos,
		decodeBookingAdditionalInfo,
		encodeTravolutionResponse,
		options...,
	)

	getAdditionalInfoByUid := kithttp.NewServer(
		svcEndpoints.GetAdditionalInfoByUid,
		decodeBookingAdditionalInfo,
		encodeTravolutionResponse,
		options...,
	)

	postCreateAllProduct := kithttp.NewServer(
		svcEndpoints.PostCreateAllProduct,
		decodeCreateProduct,
		encodeTravolutionResponse,
		options...,
	)
	// product
	router.Handle("/v1/ggt/get/products", getAllproducts).Methods("GET")
	router.Handle("/v1/ggt/get/productuid/{productUid}", getProductByUid).Methods("GET")
	router.Handle("/v1/ggt/get/options/by/productUid/{productUid}/options/", getAllOptionsOfProduct).Methods("GET")
	router.Handle("/v1/ggt/get/option/by/productUid/{productUid}/options/{optionUid}", getOptionOfProductByOptionUid).Methods("GET")
	router.Handle("/v1/ggt/get/units/by/productUid/{productUid}/options/{optionUid}/units/", getUnitSPriceByOptionUid).Methods("GET")
	router.Handle("/v1/ggt/get/units/by/productUid/{productUid}/options/{optionUid}/units/{unitUid}", getUnitPriceByOptionUidAndUnitId).Methods("GET")
	router.Handle("/v1/ggt/get/schedule/by/productUid/{productUid}/options/{optionUid}/booking-schedules/", getBookingSchedules).Methods("GEt")
	router.Handle("/v1/ggt/get/additinalInfo/by/productUid/{productUid}/option/{optionUid}/booking-additional-infos/", getAdditionalInfos).Methods("GET")
	router.Handle("/v1/ggt/get/additinalInfo/by/productUid/{productUid}/option/{optionUid}/booking-additional-infos/{additionalInfoUid}", getAdditionalInfoByUid).Methods("GET")
	router.Handle("/v1/ggt/create/travolution/products", postCreateAllProduct).Methods("POST")

	// order

	postCreateTravolutionOrder := kithttp.NewServer(
		svcEndpoints.PostCreateTravolutionOrder,
		decodeCreateTravolutionOrder,
		encodeTravolutionResponse,
		options...,
	)

	getSearchTravolutionOrder := kithttp.NewServer(
		svcEndpoints.GetSearchTravolutionOrder,
		decodeSearchTravolutionOrder,
		encodeTravolutionResponse,
		options...,
	)

	deleteCancelTravolutionOrder := kithttp.NewServer(
		svcEndpoints.PostCancelTravolutionOrder,
		decodeCancelTravolutionOrder,
		encodeTravolutionResponse,
		options...,
	)

	PostTravolutionWebHook := kithttp.NewServer(
		middleware.AuthMiddleWare(mongoRepo)(svcEndpoints.PostTravolutionWebHook),
		decodeTravolutionWebHook,
		encodeTravolutionResponse,
		options...,
	)

	router.Handle("/v1/travolution/create/order", postCreateTravolutionOrder).Methods("POST")
	router.Handle("/v1/travolution/search/order/{orderNumber}", getSearchTravolutionOrder).Methods("GET")
	router.Handle("/v1/travolution/cancel/order/{orderNumber}", deleteCancelTravolutionOrder).Methods("DELETE")
	router.Handle("/v1/travolution/webhook", PostTravolutionWebHook).Methods("POST")

	// handling of 404 not found handler
	router.NotFoundHandler = http.HandlerFunc(DefaultNotFoundRouteHandler)
	return router
}

// decodeHeartbeatRequest decodes the request from the heartbeatHandler
func decodeHeartbeatRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

// Yanolja handler
// decodeGetProductsRequest handler for GetProduts
func decodeGetProductsRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.AllProduct

	// Extract query parameters
	pno := r.URL.Query().Get("pageNumber")
	n, e := strconv.ParseInt(pno, 10, 32)
	if e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return req, e
	}
	req.PageNumber = int32(n)
	psz := r.URL.Query().Get("pageSize")
	p, e := strconv.ParseInt(psz, 10, 32)
	if e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return req, e
	}
	req.PageSize = int32(p)

	statusCode := r.URL.Query().Get("productStatusCode")
	if statusCode != "" {
		req.ProductStatusCode = statusCode
	}
	if req.ProductStatusCode != "WAITING_FOR_SALE" && req.ProductStatusCode != "IN_SALE" &&
		req.ProductStatusCode != "SOLD_OUT" && req.ProductStatusCode != "END_OF_SALE" {
		e = customError.NewError(ctx, "leisure-api-0001", "passed productstatuscode is not valid", nil)
		return req, e
	}

	/* validate := validator.New()
	err1 := validate.Struct(req)
	if err1 != nil {
		e := customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), err1)
		return req, e
	} */

	return req, nil
}

func decodeGetProductsByIdRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.ProductsById

	vars := mux.Vars(r)

	_, ok := vars["productid"]
	if ok {
		req.ProductId = vars["productid"]
	}
	return req, nil
}

func decodeGetProductsOptionGroups(ctx context.Context, r *http.Request) (request interface{}, err error) {

	var req yanolja.ProductsById

	vars := mux.Vars(r)
	_, ok := vars["productid"]
	if ok {
		req.ProductId = vars["productid"]
	}

	return req, nil
}

func decodeGetProductsInventories(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.ProductInventory

	vars := mux.Vars(r)
	_, ok := vars["productid"]
	if ok {
		req.ProductId = vars["productid"]
	}

	// Get the query parameter
	s := r.URL.Query()
	if len(s) > 0 {
		req.InventoryDateStart = s.Get("inventoryDateStart")
		req.InventoryDateEnd = s.Get("inventoryDateEnd")
	}

	return req, nil
}

func decodeGetVariantInventory(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.VariantInventory

	vars := mux.Vars(r)
	_, ok := vars["variantid"]
	if ok {
		req.VariantId = vars["variantid"]
	}

	s := r.URL.Query()
	if len(s) > 0 {
		req.Date = s.Get("date")
		req.Time = s.Get("time")
	}

	return req, nil

}

// decodeGetOrderByPartnerOrderIdSuffix
func decodeGetOrderByPartnerOrderIdSuffix(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var partnerOrdeIdSuffix string

	vars := mux.Vars(r)
	val, ok := vars["partnerOrderId"]
	if ok {
		partnerOrdeIdSuffix = val
	}

	return partnerOrdeIdSuffix, nil

}

// decodeGetCategories
func decodeGetCategories(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

// decodePostCategories
func decodePostCategories(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

// decodeGetRegions
func decodeGetRegions(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

// decodePostRegions
func decodePostRegions(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodePostWaitForOrder(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.WaitingForOrder
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	return req, nil
}

func decodePostOrderCompletion(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.OrderConfirmation
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	return req, nil
}

func decodeGetOrderByOrderId(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.OrderConfirmation

	vars := mux.Vars(r)
	val, ok := vars["orderId"]
	if ok {
		n, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		}
		req.OrderId = n
	}
	s := r.URL.Query()
	if len(s) > 0 {
		req.PartnerOrderId = s.Get("partnerOrderId")
	}
	return req, nil
}

func decodeGetEverlandOrders(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.EverlandGetRequest

	vars := mux.Vars(r)
	email, ok := vars["customeremail"]
	if !ok {
		e := customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	channelcode, ok := vars["channelcode"]
	if !ok {
		e := customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	req.ChannelCode = channelcode
	req.CustomerEmail = email
	return req, nil
}

func decodePostFullOrderCancellation(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.Order
	vars := mux.Vars(r)
	val, ok := vars["orderId"]
	if ok {
		n, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		}
		req.OrderId = n
	}

	return req, nil

}

func decodeGetOdooByOrderId(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.Order
	vars := mux.Vars(r)
	val, ok := vars["orderId"]
	if ok {
		n, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		}
		req.OrderId = n
	}

	return req, nil

}

func decodePostUpdateOdooSyncFlagFalse(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodePostUpdateOdooSyncFlag(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.OrderList
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	return req, nil
}

func decodePostReqTimeoutOrderCancel(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.PartnerOrder
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	return req, nil
}
func decodePostForcedOrderCancellation(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.PartnerOrder
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	return req, nil
}

func decodePostProductCreationInGGT(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.Upsert_Product
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

func decodeGetProductByProductId(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.GetProduct

	vars := mux.Vars(r)
	_, ok := vars["productId"]
	if ok {
		req.ProductId, err = strconv.ParseInt(vars["productId"], 10, 64)
		if err != nil {
			fmt.Printf("Error converting string to int64: %v\n", err)
			return
		}
		fmt.Println("productId :", req.ProductId)
	}

	return req, nil
}

func decodePostAllProducts(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeGetOrderReconcilationDetail(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.OrderReconcileReq

	// Extract query parameters
	req.ReconciliationDate = r.URL.Query().Get("reconciliationDate")
	req.ReconcileOrderStatusCode = r.URL.Query().Get("reconcileOrderStatusCode")

	pno := r.URL.Query().Get("pageNumber")
	n, e := strconv.Atoi(pno)
	if e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return req, e
	}
	req.PageNumber = n

	psz := r.URL.Query().Get("pageSize")
	p, e := strconv.Atoi(psz)
	if e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return req, e
	}
	req.PageSize = p

	return req, nil
}

func decodepostCancellationAckClbk(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.CancellationAck
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	return req, nil
}

func decodePostRefusalToCancelClbk(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.RefusalToCancel
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

func decodeGetOrderStatusLookupClbk(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.OrderStatusLookup
	vars := mux.Vars(r)
	val, ok := vars["partnerOrderId"]
	if !ok {
		return req, customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
	}
	req.PartnerOrderId = val

	// Extract query parameters
	orderIdStr := r.URL.Query().Get("orderId")
	orderVariantIdStr := r.URL.Query().Get("orderVariantId")

	// Convert orderId and orderVariantId to integers
	req.OrderId, err = strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		e := customError.NewError(ctx, "leisure-api-0001", "Invalid orderId", nil)
		return req, e
	}
	if orderVariantIdStr != "" {
		req.OrderVariantID, err = strconv.ParseInt(orderVariantIdStr, 10, 64)
		if err != nil {
			e := customError.NewError(ctx, "leisure-api-0001", "Invalid orderVariantId", nil)
			return req, e
		}
	}

	return req, nil
}

func decodePostForcedOrderCancellationClbk(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.ForcedOrderCancellation
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

func decodePostIndividualVoucherUpdateClbk(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.IndividualVoucherUpdate
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

func decodePostCombinedVoucherUpdateClbk(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.CombinedVoucherUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e := customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), err)
		return nil, e
	}
	return req, nil
}

func decodePostProcessingOrRestoringClbk(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.ProcessingOrRestoringReq
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

func decodePostRequestFromGGT(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req trip.SwallowRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

// decodeGetCategories
func decodeGetRedisData(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeGetOrderForOdoo(_ context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeGetProductViewForOdoo(_ context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}
func decodeGetGGTInventorySync(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodePutContent(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}
func decodePostMonitoProductUpdate(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodePostPluUpsertToRedis(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodepostUpdateImageSyncStatus(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req []domain.ImageUrlForProcessing
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

func decodeGetImageSyncToTrip(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

// check for use
func decodeUpdateTripImageSyncStatus(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req []domain.ImageDetailForSync
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}
	return req, nil
}

func decodeGetProductContentData(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeGetPackageContentData(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeGetPluToRedis(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeDeleteRedisKeys(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeUpsertCategoryMapping(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeDeleteIfNotEmpty(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodePutContentSyncData(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req trip.TripMessageForSync
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	return request, nil
}

func decodeFindKeysVal(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req common.RedisKey
	vars := mux.Vars(r)
	key, ok := vars["key"]
	if ok {
		req.Key = key
	}
	return req, nil
}

func decodeUpdatePluToRedis(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeFindPluFromRedis(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req yanolja.PluRequest

	vars := mux.Vars(r)
	_, ok := vars["productId"]
	if ok {
		req.ProductId, err = strconv.ParseInt(vars["productId"], 10, 64)
		if err != nil {
			fmt.Printf("Error converting string to int64: %v\n", err)
			return
		}
	}

	_, ok = vars["variantId"]
	if ok {
		req.VariantId, err = strconv.ParseInt(vars["variantId"], 10, 64)
		if err != nil {
			fmt.Printf("Error converting string to int64: %v\n", err)
			return
		}
	}

	var version int64
	_, ok = vars["productVersion"]
	if ok {
		version, err = strconv.ParseInt(vars["productVersion"], 10, 64)
		if err != nil {
			fmt.Printf("Error converting string to int64: %v\n", err)
			return
		}
		req.ProductVersion = int32(version)

	}
	return req, nil
}

// Travolution
// decodeGetProducts
func decodeGetProducts(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req travolution.ProductReq

	// --- Handle optional path param: /products/{productUid}
	vars := mux.Vars(r)

	if productUidStr, ok := vars["productUid"]; ok {
		productUid, err := strconv.Atoi(productUidStr)
		if err != nil {
			return nil, fmt.Errorf("invalid productUid: %v", err)
		}
		req.ProductUid = productUid
	}
	// --- Handle query params: take, skip, lang
	q := r.URL.Query()

	if takeStr := q.Get("take"); takeStr != "" {
		takeStr = strings.TrimSpace(takeStr)
		take, err := strconv.Atoi(takeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid query parameter take: %v", err)
		}
		req.Take = take
	}

	if skipStr := q.Get("skip"); skipStr != "" {
		skipStr = strings.TrimSpace(skipStr)
		skip, err := strconv.Atoi(skipStr)
		if err != nil {
			return nil, fmt.Errorf("invalid skip: %v", err)
		}
		req.Skip = skip
	}

	if lang := q.Get("lang"); lang != "" {
		lang = strings.TrimSpace(lang)
		req.Lang = lang
	}

	fmt.Println(" request received ", req)
	return req, nil
}

// decodeOption
func decodeOption(ctx context.Context, r *http.Request) (request interface{}, err error) {

	var req travolution.OptionRequest
	vars := mux.Vars(r)

	// Parse productUid
	if productUidStr, ok := vars["productUid"]; ok {
		productUid, err := strconv.Atoi(productUidStr)
		if err != nil {
			return nil, fmt.Errorf("invalid productUid: %v", err)
		}
		req.ProductUid = productUid
	}

	// Parse optionUid
	req.OptionUid, _ = vars["optionUid"]

	q := r.URL.Query()

	if lang := q.Get("lang"); lang != "" {
		req.Lang = lang
	}

	return req, nil
}

// decodeGetUnitDetails
func decodeGetUnitDetails(ctx context.Context, r *http.Request) (request interface{}, err error) {
	fmt.Println("<<<<<<<<<<<<<<<<< unit >>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	var req travolution.UnitPriceRequest

	vars := mux.Vars(r)

	// --- ProductUid (Required)
	productUidStr, ok := vars["productUid"]
	if !ok {
		return nil, fmt.Errorf("missing productUid")
	}
	productUid, err := strconv.Atoi(productUidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid productUid: %w", err)
	}
	req.ProductUid = productUid

	// --- OptionUid (Required: int or PKG/PAS)
	optionUidStr, ok := vars["optionUid"]
	if !ok {
		return nil, fmt.Errorf("missing optionUid")
	}
	if optInt, err := strconv.Atoi(optionUidStr); err == nil {
		req.OptionUid = optInt
	} else if optionUidStr == "PKG" || optionUidStr == "PAS" {
		req.OptionUid = optionUidStr
	} else {
		return nil, fmt.Errorf("invalid optionUid: must be int or PKG/PAS")
	}

	// --- UnitUid (Optional: int or PKG/PAS)
	unitUidStr, ok := vars["unitUid"]
	if ok && unitUidStr != "" {
		if unitInt, err := strconv.Atoi(unitUidStr); err == nil {
			req.UnitUid = unitInt
		} else if unitUidStr == "PKG" || unitUidStr == "PAS" {
			req.UnitUid = unitUidStr
		} else {
			req.UnitUid = nil
		}
	}

	fmt.Println("::::::::::::::::::::::::: req ::::::::::::::::::::::::::: ", req)

	return req, nil
}

// decodeBookingSchedule
func decodeBookingSchedule(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)

	// --- Extract productUid (required)
	productUidStr, ok := vars["productUid"]
	if !ok {
		return nil, errors.New("missing productUid in path")
	}
	productUid, err := strconv.Atoi(productUidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid productUid: %w", err)
	}

	// --- Extract optionUid (required, string or int as string)
	optionUid, ok := vars["optionUid"]
	if !ok {
		return nil, errors.New("missing optionUid in path")
	}

	// --- Extract optional query parameters
	query := r.URL.Query()
	date := query.Get("date") // e.g., 20250215
	time := query.Get("time") // e.g., 0900

	req := travolution.BookingScheduleReq{
		ProductUid: productUid,
		OptionUid:  optionUid,
		Date:       date,
		Time:       time,
	}

	return req, nil
}

// decodeBookingAdditionalInfo
func decodeBookingAdditionalInfo(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)

	// Extract and validate productUid (required integer)
	productUidStr, ok := vars["productUid"]
	if !ok {
		return nil, errors.New("missing productUid in path")
	}
	productUid, err := strconv.Atoi(productUidStr)
	if err != nil {
		return nil, fmt.Errorf("invalid productUid: %v", err)
	}

	// Extract and validate optionUid (required, type interface{})
	optionUid, ok := vars["optionUid"]
	if !ok || optionUid == "" {
		return nil, errors.New("missing optionUid in path")
	}

	// Extract additionalInfoUid (optional, type interface{})
	var additionalInfoUid interface{}
	if val, ok := vars["additionalInfoUid"]; ok && val != "" {
		additionalInfoUid = val
	}

	// Return the decoded request object
	return travolution.BookingAdditionalInfoRequest{
		ProductUID:        productUid,
		OptionUID:         optionUid,         // interface{}
		AdditionalInfoUID: additionalInfoUid, // interface{}
	}, nil
}

// decodeCreateProduct
func decodeCreateProduct(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

// decodeGetProductByIdFromGGT
func decodeGetProductByIdFromGGT(ctx context.Context, r *http.Request) (request interface{}, err error) {

	vars := mux.Vars(r)
	_, ok := vars["productId"]
	if !ok {
		err = customError.NewError(ctx, "leisure-api-0001", "productId not passed as  path parameter ", nil)
		return nil, err
	}

	productId, err := strconv.ParseInt(vars["productId"], 10, 64)
	if err != nil {
		fmt.Printf("Error converting string to int64: %v\n", err)
		return
	}

	return productId, nil
}

func decodeGetProductsFromGGT(ctx context.Context, r *http.Request) (request interface{}, err error) {
	return request, nil
}

func decodeCreateTravolutionOrder(ctx context.Context, r *http.Request) (request interface{}, err error) {

	// ---- Dump body ----
	/*fmt.Println("=== Body ===")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	// Important: restore Body so it can be re-read later
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	fmt.Println("///////////////////////////////////////////////////////// ", string(bodyBytes))
	*/
	var req travolution.OrderRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	fmt.Println(" <<<<<<<<<<<< request received >>>>>>>>>>> ", req)
	return req, nil
}

func decodeSearchTravolutionOrder(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	orderNumber, ok := vars["orderNumber"]
	if !ok {
		err = customError.NewError(ctx, "leisure-api-0001", "productId not passed as  path parameter ", nil)
		return nil, err
	}

	return orderNumber, nil

}

func decodeCancelTravolutionOrder(ctx context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	orderNumber, ok := vars["orderNumber"]
	if !ok {
		err = customError.NewError(ctx, "leisure-api-0001", "productId not passed as  path parameter ", nil)
		return nil, err
	}

	return orderNumber, nil
}

func decodeTravolutionWebHook(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req travolution_domain.Webhook
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		e = customError.NewError(ctx, "leisure-api-0001", customError.ErrInvalidBody.Error(), nil)
		return nil, e
	}

	fmt.Println(" <<<<<<<<<<<< request received >>>>>>>>>>> ", req)
	return req, nil

}

// encodeResponse encodes final response as json
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rs, _ := response.(yanolja.Response)
	rs.TraceID = utils.GenerateUUID("", true)
	status, _ := strconv.Atoi(rs.Code)
	if status < 200 || status >= 300 {
		/* var errorMessage string

		switch v := rs.Body.(type) {
		case string:
			errorMessage = v
		case []byte:
			errorMessage = string(v)
		case io.ReadCloser:
			defer v.Close()
			b, err := io.ReadAll(v)
			if err != nil {
				errorMessage = "failed to read body"
			} else {
				errorMessage = string(b)
			}
		default:
			errorMessage = fmt.Sprintf("unsupported body type: %T", rs.Body)
		} */

		rs.Body = yanolja.ErrorMsg{
			ErrorCode:    rs.Code,
			ErrorMessage: rs.Body.(string),
		}
	}
	var res1 yanolja.SupplierResponse
	res1.TraceID = rs.TraceID
	res1.Body = rs.Body
	res1.Page = rs.Page
	res1.Collection = rs.Collection
	res1.ContentType = rs.ContentType
	//w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res1)
}

// encodeResponse encodes final response as json
func encodeCommonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rs, _ := response.(common.Response)
	rs.TraceID = utils.GenerateUUID("", true)
	status, _ := strconv.Atoi(rs.Code)
	if status < 200 || status >= 300 {
		/* var errorMessage string

		switch v := rs.Body.(type) {
		case string:
			errorMessage = v
		case []byte:
			errorMessage = string(v)
		case io.ReadCloser:
			defer v.Close()
			b, err := io.ReadAll(v)
			if err != nil {
				errorMessage = "failed to read body"
			} else {
				errorMessage = string(b)
			}
		default:
			errorMessage = fmt.Sprintf("unsupported body type: %T", rs.Body)
		} */

		rs.Body = common.ErrorResponse{
			Code:    rs.Code,
			Message: rs.Body.(string),
			Status:  rs.Status,
		}
	}
	var res1 common.Response
	res1.TraceID = rs.TraceID
	res1.Body = rs.Body
	res1.Code = rs.Code
	res1.Status = rs.Status

	//w.WriteHeader(status)
	return json.NewEncoder(w).Encode(res1)
}

func encodeTravolutionResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	rs, _ := response.(travolution.Response)
	status, _ := strconv.Atoi(rs.Code)
	if status < 200 || status >= 300 {
		rs.Body = travolution.ErrorMsg{
			ErrorCode:    rs.Code,
			ErrorMessage: rs.Body.(string),
		}
	}

	rs.HtmlTypeContent = nil
	//w.WriteHeader(status)
	return json.NewEncoder(w).Encode(rs)
}

// DefaultNotFoundRouteHandler handler for 404 resource not found
func DefaultNotFoundRouteHandler(w http.ResponseWriter, req *http.Request) {
	logger := log.NewLogfmtLogger(os.Stdout)
	level.Info(logger).Log(" request ", req)

	var responseError = customError.NewError(context.Background(), "leisure-api-0006", customError.ErrResourceNotFound.Error(), nil)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(responseError.Status)
	json.NewEncoder(w).Encode(responseError)
}
