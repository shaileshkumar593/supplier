package middleware

import (
	"context"
	"errors"
	"time"

	customContext "swallow-supplier/context"
	customError "swallow-supplier/error"
	svc "swallow-supplier/iface"
	travolution_domain "swallow-supplier/mongo/domain/travolution"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const (
	actionKey   = "action"
	durationKey = "duration"
	requestKey  = "request"
	responseKey = "response"
	errorKey    = "err"
)

type loggingMiddleware struct {
	logger kitlog.Logger
	next   svc.Service
}

// NewLoggingMiddleware for handling logging across all routes
func NewLoggingMiddleware(logger kitlog.Logger) ServiceMiddleware {
	return func(next svc.Service) svc.Service {
		return loggingMiddleware{logger, next}
	}
}

// HeartBeat
// Logs the request and response of the /heartbeat endpoint
func (mw loggingMiddleware) HeartBeat(ctx context.Context) (res yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, res, err)
	}(time.Now())

	return mw.next.HeartBeat(ctx)
}

// GetProducts get all products from yanolja
func (mw loggingMiddleware) GetProducts(ctx context.Context, req yanolja.AllProduct) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.GetProducts(ctx, req)
}

// GetProductsById GetProducts get products by id from yanolja
func (mw loggingMiddleware) GetProductsById(ctx context.Context, req yanolja.ProductsById) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetProductsById(ctx, req)
}

// GetProductsOptionGroups get products option group from yanolja
func (mw loggingMiddleware) GetProductsOptionGroups(ctx context.Context, req yanolja.ProductsById) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetProductsOptionGroups(ctx, req)
}

// GetProductsInventories get  products inventories  from yanolja
func (mw loggingMiddleware) GetProductsInventories(ctx context.Context, req yanolja.ProductInventory) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetProductsInventories(ctx, req)
}

// GetVariantInventory get product variant inventory from yanolja
func (mw loggingMiddleware) GetVariantInventory(ctx context.Context, req yanolja.VariantInventory) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetVariantInventory(ctx, req)
}

// GetCategories get all categories from yanolja
func (mw loggingMiddleware) GetCategories(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.GetCategories(ctx)
}

func (mw loggingMiddleware) InsertAllCategories(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.InsertAllCategories(ctx)
}

// GetRegions get all regional information from yanolja
func (mw loggingMiddleware) GetRegions(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.GetRegions(ctx)
}

// InsertAllRegions for storing regions detail
func (mw loggingMiddleware) InsertAllRegions(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.InsertAllRegions(ctx)
}

// PostWaitForOrder preorder order creation from yanolja
func (mw loggingMiddleware) PostWaitForOrder(ctx context.Context, req yanolja.WaitingForOrder) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.PostWaitForOrder(ctx, req)
}

// PostOrderCompletion confirm the order creation
func (mw loggingMiddleware) PostOrderCompletion(ctx context.Context, req yanolja.OrderConfirmation) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.PostOrderCompletion(ctx, req)
}

// GetOrderByOrderId get the order by id
func (mw loggingMiddleware) GetOrderByOrderId(ctx context.Context, req yanolja.OrderConfirmation) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetOrderByOrderId(ctx, req)
}

// PostCancelOrderEntirly fully cancel the order
func (mw loggingMiddleware) PostCancelOrderEntirly(ctx context.Context, orderId int64) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, orderId, resp, err)
	}(time.Now())

	return mw.next.PostCancelOrderEntirly(ctx, orderId)
}

// PostCancelOrderByReqTimeOut cancel order when request timeout
func (mw loggingMiddleware) PostCancelOrderByReqTimeOut(ctx context.Context, partnerOrderId string) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, partnerOrderId, resp, err)
	}(time.Now())

	return mw.next.PostCancelOrderByReqTimeOut(ctx, partnerOrderId)
}

// PostForcelyCancelOrder forcely cancel order during failure
func (mw loggingMiddleware) PostForcelyCancelOrder(ctx context.Context, partnerOrderId string) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, partnerOrderId, resp, err)
	}(time.Now())

	return mw.next.PostForcelyCancelOrder(ctx, partnerOrderId)
}

func (mw loggingMiddleware) GetProductByProductId(ctx context.Context, productId int64) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, productId, resp, err)
	}(time.Now())

	return mw.next.GetProductByProductId(ctx, productId)
}

func (mw loggingMiddleware) InsertAllProduct(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.InsertAllProduct(ctx)
}

// GetOrderReconcilationDetail OrderReconcilationDetail
func (mw loggingMiddleware) GetOrderReconcilationDetail(ctx context.Context, req yanolja.OrderReconcileReq) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetOrderReconcilationDetail(ctx, req)
}

// InsertProductClbk ---------------------------------callback-----------------------------------------------------------------
func (mw loggingMiddleware) InsertProductClbk(ctx context.Context, product yanolja.Upsert_Product) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, product, resp, err)
	}(time.Now())

	return mw.next.InsertProductClbk(ctx, product)
}

// CancellationAckClbk CancellationAck
func (mw loggingMiddleware) CancellationAckClbk(ctx context.Context, ackReq yanolja.CancellationAck) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, ackReq, resp, err)
	}(time.Now())

	return mw.next.CancellationAckClbk(ctx, ackReq)
}

// RefusalToCancelClbk RefusalToCancel
func (mw loggingMiddleware) RefusalToCancelClbk(ctx context.Context, refusalToCancel yanolja.RefusalToCancel) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, refusalToCancel, resp, err)
	}(time.Now())

	return mw.next.RefusalToCancelClbk(ctx, refusalToCancel)
}

// OrderStatusLookupClbk OrderStatusLookup
func (mw loggingMiddleware) OrderStatusLookupClbk(ctx context.Context, lookup yanolja.OrderStatusLookup) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, lookup, resp, err)
	}(time.Now())

	return mw.next.OrderStatusLookupClbk(ctx, lookup)
}

// ForcedOrderCancellationClbk ForcedOrderCancellation
func (mw loggingMiddleware) ForcedOrderCancellationClbk(ctx context.Context, cancellation yanolja.ForcedOrderCancellation) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, cancellation, resp, err)
	}(time.Now())

	return mw.next.ForcedOrderCancellationClbk(ctx, cancellation)
}

// IndividualVoucherUpdateClbk IndividualVoucherUpdate
func (mw loggingMiddleware) IndividualVoucherUpdateClbk(ctx context.Context, voucher yanolja.IndividualVoucherUpdate) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, voucher, resp, err)
	}(time.Now())

	return mw.next.IndividualVoucherUpdateClbk(ctx, voucher)
}

// CombinedVoucherUpdateClbk CombinedVoucherUpdate
func (mw loggingMiddleware) CombinedVoucherUpdateClbk(ctx context.Context, voucher yanolja.CombinedVoucherUpdate) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, voucher, resp, err)
	}(time.Now())

	return mw.next.CombinedVoucherUpdateClbk(ctx, voucher)
}

func (mw loggingMiddleware) ProcessingOrRestoringClbk(ctx context.Context, req yanolja.ProcessingOrRestoringReq) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.ProcessingOrRestoringClbk(ctx, req)
}

// GetProductSync -----------------------odoo related -------------------------------------------------------------------
func (mw loggingMiddleware) GetProductSync(ctx context.Context) (resp common.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.GetProductSync(ctx)
}

func (mw loggingMiddleware) GetOrderSync(ctx context.Context) (resp common.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.GetOrderSync(ctx)
}

// InventorySync ----------------------------------------------------------------------------------------------------------------
func (mw loggingMiddleware) InventorySync(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.InventorySync(ctx)
}

func (mw loggingMiddleware) UpdateImageSyncStatus(ctx context.Context, req []domain.ImageUrlForProcessing) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.UpdateImageSyncStatus(ctx, req)
}

func (mw loggingMiddleware) UpdateTripImageSyncStatus(ctx context.Context, req []trip.ImageSyncToTrip) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.UpdateTripImageSyncStatus(ctx, req)
}

// PostRequestFromGGT -----------------------------------------GGT------------------------------------------------------------------------------
func (mw loggingMiddleware) PostRequestFromGGT(ctx context.Context, req trip.SwallowRequest) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.PostRequestFromGGT(ctx, req)
}

func (mw loggingMiddleware) GetRedisData(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.GetRedisData(ctx)
}

func (mw loggingMiddleware) GetImageSyncToTrip(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.GetImageSyncToTrip(ctx)
}

func (mw loggingMiddleware) ProductContentSync(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.ProductContentSync(ctx)
}

func (mw loggingMiddleware) PackageContentSync(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.PackageContentSync(ctx)
}

func (mw loggingMiddleware) UpdateContentSyncStatus(ctx context.Context, contentSyncStatus trip.TripMessageForSync) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, contentSyncStatus, resp, err)
	}(time.Now())

	return mw.next.UpdateContentSyncStatus(ctx, contentSyncStatus)
}

func (mw loggingMiddleware) ProductSyncToTrip(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.ProductSyncToTrip(ctx)
}

func (mw loggingMiddleware) GetPluToRedis(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.GetPluToRedis(ctx)
}

func (mw loggingMiddleware) UpsertCategoryMapping(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.UpsertCategoryMapping(ctx)
}

func (mw loggingMiddleware) DeleteRecoRdIfNotEmpty(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.DeleteRecoRdIfNotEmpty(ctx)
}

func (mw loggingMiddleware) MonitorProductUpdateSvc(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.MonitorProductUpdateSvc(ctx)
}

func (mw loggingMiddleware) SyncAllPluToRedis(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.SyncAllPluToRedis(ctx)
}

func (mw loggingMiddleware) DeleteRedisData(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.DeleteRedisData(ctx)
}

func (mw loggingMiddleware) FindRedisKeyValue(ctx context.Context, key string) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, key, resp, err)
	}(time.Now())

	return mw.next.FindRedisKeyValue(ctx, key)
}

func (mw loggingMiddleware) UpdatePluToRedis(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.UpdatePluToRedis(ctx)
}

func (mw loggingMiddleware) GetPluFromRedis(ctx context.Context, req yanolja.PluRequest) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetPluFromRedis(ctx, req)
}

func (mw loggingMiddleware) GetOrderSyncByOrderIdInOdoo(ctx context.Context, req yanolja.Order) (resp common.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetOrderSyncByOrderIdInOdoo(ctx, req)
}

func (mw loggingMiddleware) PostOdooSyncStatus(ctx context.Context, req yanolja.OrderList) (resp common.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.PostOdooSyncStatus(ctx, req)
}

func (mw loggingMiddleware) GetEverlandOrders(ctx context.Context, req yanolja.EverlandGetRequest) (resp common.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetEverlandOrders(ctx, req)
}

// Travolution

func (mw loggingMiddleware) GetAllproducts(ctx context.Context, req travolution.ProductReq) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetAllproducts(ctx, req)
}

func (mw loggingMiddleware) GetProductByUid(ctx context.Context, req travolution.ProductReq) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetProductByUid(ctx, req)
}

func (mw loggingMiddleware) GetAllOptionsOfProduct(ctx context.Context, req travolution.OptionRequest) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetAllOptionsOfProduct(ctx, req)
}

func (mw loggingMiddleware) GetOptionOfProductByOptionUid(ctx context.Context, req travolution.OptionRequest) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetOptionOfProductByOptionUid(ctx, req)
}

func (mw loggingMiddleware) GetUnitSPriceByOptionUid(ctx context.Context, req travolution.UnitPriceRequest) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetUnitSPriceByOptionUid(ctx, req)
}

func (mw loggingMiddleware) GetUnitPriceByOptionUid(ctx context.Context, req travolution.UnitPriceRequest) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetUnitPriceByOptionUid(ctx, req)
}

func (mw loggingMiddleware) GetBookingSchedules(ctx context.Context, req travolution.BookingScheduleReq) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetBookingSchedules(ctx, req)
}

func (mw loggingMiddleware) GetAdditionalInfos(ctx context.Context, req travolution.BookingAdditionalInfoRequest) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetAdditionalInfos(ctx, req)
}

func (mw loggingMiddleware) GetAdditionalInfoByUid(ctx context.Context, req travolution.BookingAdditionalInfoRequest) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.GetAdditionalInfoByUid(ctx, req)
}

func (mw loggingMiddleware) PostCreateAllProduct(ctx context.Context) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, nil, resp, err)
	}(time.Now())

	return mw.next.PostCreateAllProduct(ctx)
}

func (mw loggingMiddleware) GetOderByPartialPartnerOrderIdSuffix(ctx context.Context, partialPartnerId string) (resp common.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, partialPartnerId, resp, err)
	}(time.Now())

	return mw.next.GetOderByPartialPartnerOrderIdSuffix(ctx, partialPartnerId)
}

func (mw loggingMiddleware) PostChangeOdooSyncStatusToFalse(ctx context.Context) (resp common.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.PostChangeOdooSyncStatusToFalse(ctx)
}

func (mw loggingMiddleware) GetProductsFromGGT(ctx context.Context) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, "", resp, err)
	}(time.Now())

	return mw.next.GetProductsFromGGT(ctx)
}

func (mw loggingMiddleware) GetProductByIdFromGGT(ctx context.Context, productId int64) (resp yanolja.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, productId, resp, err)
	}(time.Now())

	return mw.next.GetProductByIdFromGGT(ctx, productId)
}

func (mw loggingMiddleware) CreateTravolutionOrder(ctx context.Context, req travolution.OrderRequest) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, req, resp, err)
	}(time.Now())

	return mw.next.CreateTravolutionOrder(ctx, req)
}

func (mw loggingMiddleware) SearchTravolutionOrder(ctx context.Context, orderNumber string) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, orderNumber, resp, err)
	}(time.Now())

	return mw.next.SearchTravolutionOrder(ctx, orderNumber)
}

func (mw loggingMiddleware) CancelTravolutionOrder(ctx context.Context, orderNumber string) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, orderNumber, resp, err)
	}(time.Now())

	return mw.next.CancelTravolutionOrder(ctx, orderNumber)
}

func (mw loggingMiddleware) OrderWebhookUpdate(ctx context.Context, payload travolution_domain.Webhook) (resp travolution.Response, err error) {
	defer func(startTime time.Time) {
		logRequest(ctx, mw.logger, startTime, payload, resp, err)
	}(time.Now())

	return mw.next.OrderWebhookUpdate(ctx, payload)
}

// logRequest log the requests
func logRequest(ctx context.Context, logger kitlog.Logger, startTime time.Time, req interface{}, res interface{}, err error) {
	if err == nil {
		level.Info(logger).Log(
			actionKey, customContext.CtxRequestPath(ctx, customContext.CtxLabelRequestPath),
			durationKey, duration(startTime),
			customContext.CtxLabelTraceID, customContext.CtxTraceID(ctx),
			customContext.CtxLabelRequestID, customContext.GetCtxHeader(ctx, customContext.CtxLabelRequestID),
			customContext.CtxLabelChannelCode, customContext.GetCtxHeader(ctx, customContext.CtxLabelChannelCode),
			requestKey, req,
			responseKey, res)
		return
	}

	// log the error response
	var re *customError.ResponseError
	if errors.As(err, &re) {
		// If res is customError.ResponseError (value):
		res = *re
	}

	logger.Log(
		actionKey, customContext.CtxRequestPath(ctx, customContext.CtxLabelRequestPath),
		durationKey, duration(startTime),
		customContext.CtxLabelTraceID, customContext.CtxTraceID(ctx),
		customContext.CtxLabelRequestID, customContext.GetCtxHeader(ctx, customContext.CtxLabelRequestID),
		customContext.CtxLabelChannelCode, customContext.GetCtxHeader(ctx, customContext.CtxLabelChannelCode),
		requestKey, req,
		responseKey, res,
		errorKey, err)
}

// duration gets the elapsed time of the calls
func duration(startTime time.Time) time.Duration {
	return time.Since(startTime)
}
