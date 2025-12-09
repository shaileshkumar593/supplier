package iface

import (
	"context"

	travolution_domain "swallow-supplier/mongo/domain/travolution"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"
	req_resp "swallow-supplier/request_response/yanolja"
)

// Service will hold all the available methods on the Service
type Service interface {
	// HeartBeat heartbeat implementation
	HeartBeat(ctx context.Context) (res yanolja.Response, err error)

	// -------------------------Yanolja-----------------------------------------------------------------------
	// GetProducts get all products from yanolja
	GetProducts(ctx context.Context, req yanolja.AllProduct) (resp yanolja.Response, err error)

	// GetProductsById get a products from yanolja
	GetProductsById(ctx context.Context, req yanolja.ProductsById) (resp yanolja.Response, err error)

	// GetProductsOptionGroups get all products from yanolja
	GetProductsOptionGroups(ctx context.Context, req yanolja.ProductsById) (resp yanolja.Response, err error)

	// GetProductsInventories get all products from yanolja
	GetProductsInventories(ctx context.Context, req yanolja.ProductInventory) (resp yanolja.Response, err error)

	// GetVariantInventory get all products from yanolja
	GetVariantInventory(ctx context.Context, req yanolja.VariantInventory) (resp yanolja.Response, err error)

	// GetCategories get all products from yanolja
	GetCategories(ctx context.Context) (resp yanolja.Response, err error)

	// InsertAllCategories insert all product categories to ggt
	InsertAllCategories(ctx context.Context) (resp yanolja.Response, err error)

	// GetRegions get all products from yanolja
	GetRegions(ctx context.Context) (resp yanolja.Response, err error)

	// InsertAllRegions  insert all regions into ggt
	InsertAllRegions(ctx context.Context) (resp yanolja.Response, err error)

	// postWaitForOrder preorder order creation from yanolja
	PostWaitForOrder(ctx context.Context, req yanolja.WaitingForOrder) (resp yanolja.Response, err error)

	// postOrderCompletion confirm the order creation
	PostOrderCompletion(ctx context.Context, req yanolja.OrderConfirmation) (resp yanolja.Response, err error)

	// GetOrderByOrderId get the order by id
	GetOrderByOrderId(ctx context.Context, req yanolja.OrderConfirmation) (resp yanolja.Response, err error)

	// PostCancelOrderEntirly fully cancel the order
	PostCancelOrderEntirly(ctx context.Context, orderId int64) (resp yanolja.Response, err error)

	// PostCancelOrderByReqTimeOut cancel order when request timeout
	PostCancelOrderByReqTimeOut(ctx context.Context, partnerOrderId string) (resp yanolja.Response, err error)

	// PostForcelyCancelOrder forcely cancel order during failure
	PostForcelyCancelOrder(ctx context.Context, partnerOrderId string) (resp yanolja.Response, err error)

	// InsertAllProduct
	InsertAllProduct(ctx context.Context) (resp yanolja.Response, err error)

	// GetProductByProductId
	GetProductByProductId(ctx context.Context, productId int64) (resp yanolja.Response, err error)

	// OrderReconcilationDetail
	GetOrderReconcilationDetail(ctx context.Context, req yanolja.OrderReconcileReq) (resp yanolja.Response, err error)

	// ---------------------------------callback-----------------------------------------------------------------
	//InsertProducts  insert the product received in request
	InsertProductClbk(ctx context.Context, product yanolja.Upsert_Product) (resp yanolja.Response, err error)

	// CancellationAck
	CancellationAckClbk(ctx context.Context, ackreq yanolja.CancellationAck) (resp yanolja.Response, err error)

	// RefusalToCancel
	RefusalToCancelClbk(ctx context.Context, refusaltocancel yanolja.RefusalToCancel) (resp yanolja.Response, err error)

	// OrderStatusLookup
	OrderStatusLookupClbk(crx context.Context, lookup yanolja.OrderStatusLookup) (resp yanolja.Response, err error)

	// ForcedOrderCancellation
	ForcedOrderCancellationClbk(ctx context.Context, cancellation yanolja.ForcedOrderCancellation) (resp yanolja.Response, err error)

	// IndividualVoucherUpdate
	IndividualVoucherUpdateClbk(ctx context.Context, voucher yanolja.IndividualVoucherUpdate) (resp yanolja.Response, err error)

	// CombinedVoucherUpdate
	CombinedVoucherUpdateClbk(ctx context.Context, voucher yanolja.CombinedVoucherUpdate) (resp yanolja.Response, err error)

	// ProcessingOrRestoringClbk
	ProcessingOrRestoringClbk(ctx context.Context, req yanolja.ProcessingOrRestoringReq) (resp yanolja.Response, err error)

	//TranslateImages
	//TranslateImagesToText(ctx context.Context, req yanolja.ImageUrl) (resp yanolja.Response, err error)

	// InventorySync
	InventorySync(ctx context.Context) (resp yanolja.Response, err error)

	// UpdateImageSyncStatus
	UpdateImageSyncStatus(ctx context.Context, req []domain.ImageUrlForProcessing) (resp yanolja.Response, err error)

	//--------------------------------------------------------GGT Services--------------------------------------------------------------------------------
	// PostRequestFromGGT
	PostRequestFromGGT(ctx context.Context, req trip.SwallowRequest) (resp yanolja.Response, err error)

	// GetRedisData
	GetRedisData(ctx context.Context) (resp yanolja.Response, err error)

	// GetImageSyncStatusToTrip
	GetImageSyncToTrip(ctx context.Context) (resp yanolja.Response, err error)

	// UpdateTripImageSyncStatus
	UpdateTripImageSyncStatus(ctx context.Context, req []trip.ImageSyncToTrip) (resp yanolja.Response, err error)

	// ProductContentSync
	ProductContentSync(ctx context.Context) (resp yanolja.Response, err error)

	// PackageContentSync
	PackageContentSync(ctx context.Context) (resp yanolja.Response, err error)

	// UpdateContentSyncStatus
	UpdateContentSyncStatus(ctx context.Context, contentSyncStatus trip.TripMessageForSync) (resp yanolja.Response, err error)

	// ProductSyncToTrip
	ProductSyncToTrip(ctx context.Context) (resp yanolja.Response, err error)

	// GetPluToRedis
	GetPluToRedis(ctx context.Context) (resp yanolja.Response, err error)

	// UpsertCategoryMapping
	UpsertCategoryMapping(ctx context.Context) (resp req_resp.Response, err error)

	// DeleteRecoRdIfNotEmpty
	DeleteRecoRdIfNotEmpty(ctx context.Context) (resp req_resp.Response, err error)

	// MonitorProductUpdateSvc
	MonitorProductUpdateSvc(ctx context.Context) (resp yanolja.Response, err error)

	// manuall testing of plu
	SyncAllPluToRedis(ctx context.Context) (resp yanolja.Response, err error)

	// DeleteRedisData
	DeleteRedisData(ctx context.Context) (resp yanolja.Response, err error)

	// FindRedisKeyValue
	FindRedisKeyValue(ctx context.Context, key string) (resp yanolja.Response, err error)

	// UpdatePluToRedis
	UpdatePluToRedis(ctx context.Context) (resp yanolja.Response, err error)

	// GetPluFromRedis
	GetPluFromRedis(ctx context.Context, req yanolja.PluRequest) (resp yanolja.Response, err error)

	// GetOrderSyncByOrderIdInOdoo
	GetOrderSyncByOrderIdInOdoo(ctx context.Context, req yanolja.Order) (resp common.Response, err error)

	// PostOdooSyncStatus
	PostOdooSyncStatus(ctx context.Context, req yanolja.OrderList) (resp common.Response, err error)

	// GetProductViewSync   odoo product sync
	GetProductSync(ctx context.Context) (resp common.Response, err error)

	// GetOrderSync  odoo order sync
	GetOrderSync(ctx context.Context) (resp common.Response, err error)

	// GetEverlandOrders
	GetEverlandOrders(ctx context.Context, req yanolja.EverlandGetRequest) (resp common.Response, err error)

	// ::::::::::::::::::::::::::::::::::::::::Travolution:::::::::::::::::::::::::::::::::::::::::::::::::::

	// GetAllproducts
	GetAllproducts(ctx context.Context, req travolution.ProductReq) (resp travolution.Response, err error)

	// GetProductByUid
	GetProductByUid(ctx context.Context, req travolution.ProductReq) (resp travolution.Response, err error)

	// GetAllOptionsOfProduct
	GetAllOptionsOfProduct(ctx context.Context, req travolution.OptionRequest) (resp travolution.Response, err error)

	// GetOptionOfProductByOptionUid
	GetOptionOfProductByOptionUid(ctx context.Context, req travolution.OptionRequest) (resp travolution.Response, err error)

	// GetUnitSPriceByOptionUid
	GetUnitSPriceByOptionUid(ctx context.Context, req travolution.UnitPriceRequest) (resp travolution.Response, err error)

	// GetUnitPriceByOptionUid
	GetUnitPriceByOptionUid(ctx context.Context, req travolution.UnitPriceRequest) (resp travolution.Response, err error)

	// GetBookingSchedules
	GetBookingSchedules(ctx context.Context, req travolution.BookingScheduleReq) (resp travolution.Response, err error)

	// GetAdditionalInfos
	GetAdditionalInfos(ctx context.Context, req travolution.BookingAdditionalInfoRequest) (resp travolution.Response, err error)

	// GetAdditionalInfoByUid
	GetAdditionalInfoByUid(ctx context.Context, req travolution.BookingAdditionalInfoRequest) (resp travolution.Response, err error)

	// PostCreateAllProduct
	PostCreateAllProduct(ctx context.Context) (resp travolution.Response, err error)

	// PostChangeOdooSyncStatusToFalse
	PostChangeOdooSyncStatusToFalse(ctx context.Context) (resp common.Response, err error)

	// GetOderByPartialPartnerOrderIdSuffix
	GetOderByPartialPartnerOrderIdSuffix(ctx context.Context, partialPartnerId string) (resp common.Response, err error)

	// GetProductsFromGGT
	GetProductsFromGGT(ctx context.Context) (resp yanolja.Response, err error)

	// GetProductByIdFromGGT
	GetProductByIdFromGGT(ctx context.Context, productId int64) (resp yanolja.Response, err error)

	// CreateTravolutionOrder
	CreateTravolutionOrder(ctx context.Context, req travolution.OrderRequest) (resp travolution.Response, err error)

	// SearchTravolutionOrder
	SearchTravolutionOrder(ctx context.Context, orderNumber string) (resp travolution.Response, err error)

	// CancelTravolutionOrder
	CancelTravolutionOrder(ctx context.Context, orderNumber string) (resp travolution.Response, err error)

	// OrderWebhookUpdate
	OrderWebhookUpdate(ctx context.Context, payload travolution_domain.Webhook) (resp travolution.Response, err error)
}
