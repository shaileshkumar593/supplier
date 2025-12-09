package iface

import (
	"context"
	"swallow-supplier/mongo/domain/odoo"
	travolution_domain "swallow-supplier/mongo/domain/travolution"
	trip_domain "swallow-supplier/mongo/domain/trip"
	"swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/heartbeat"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/request_response/trip"
	req_resp "swallow-supplier/request_response/yanolja"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoRepository interface defines the methods for the repository
type MongoRepository interface {
	GetDbTx(ctx context.Context) (mongo.Session, error)
	GetMongoDb(ctx context.Context) *mongo.Database
	GetMongoClient(ctx context.Context) *mongo.Client
	CommitTransaction(ctx context.Context) error
	AbortTransaction(ctx context.Context) error

	Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, collection string, documents []interface{}) (*mongo.InsertManyResult, error)
	Find(ctx context.Context, collection string, filter interface{}) (bson.M, error)
	FindMany(ctx context.Context, collection string, filter interface{}) ([]bson.M, error)
	Update(ctx context.Context, collection string, filter, update interface{}) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, collection string, filter, update interface{}) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error)

	GetHeartBeatFromMongo(ctx context.Context) (res heartbeat.MongoResponse, err error)

	// InsertPreOrder insert pre-order to the order datadase
	InsertPreOrder(ctx context.Context, rec req_resp.WaitingForOrder) (id string, err error)

	//UpdatePreOrderById
	UpdatePreOrderById(ctx context.Context, update map[string]any) (id string, err error)

	// GetOrderbyOrderId
	GetOrderbyOrderId(ctx context.Context, orderid int64) (record yanolja.Model, err error)

	//GetOrderbyPartnerOrderId
	GetOrderbyPartnerOrderId(ctx context.Context, partnerOrderId string) (record yanolja.Model, err error)

	// UpdateOrderByOrderId
	UpdateOrderByOrderId(ctx context.Context, orderid int64, update map[string]any) (id string, err error)

	// DeleteOrderByOrderId
	DeleteOrderByOrderId(ctx context.Context, orderid int64) (id string, err error)

	// InsertCategories
	InsertCategories(ctx context.Context, categories []yanolja.Category) (err error)

	// UpdateCategoriesByCategoryId
	UpdateCategoriesByCategoryId(ctx context.Context, categoryId int64, update map[string]any) (Id string, err error)

	//GetTotalQuantityPurchasedByPerson
	GetTotalQuantityPurchasedByPersonToday(ctx context.Context, name, email, tel string, productId int64, variantId int64) (totalQuantity int32, err error)

	// FindCategory
	FindCategory(ctx context.Context) (record []yanolja.Category, err error)

	//FindCategoryByCategoryId
	FindCategoryByCategoryId(ctx context.Context, categoryId int64) (record yanolja.Category, err error)

	//DeleteCategoryByCategoryId
	DeleteCategoryByCategoryId(ctx context.Context, categoryId int64) (id string, err error)

	//InsertRegions
	InsertRegions(ctx context.Context, regions []yanolja.Region) (err error)

	//UpdateRegionsByRegionId
	UpdateRegionsByRegionId(ctx context.Context, regiondId int64, update map[string]any) (Id string, err error)

	// FindRegion
	FindRegion(ctx context.Context) (records []yanolja.Region, err error)

	//FindRecordByRegionId
	FindRecordByRegionId(ctx context.Context, regionId int64) (record yanolja.Region, err error)

	//DeleteRegionByRegionId
	DeleteRegionByRegionId(ctx context.Context, regionId int64) (id string, err error)

	//InsertProducts insert all products
	InsertProducts(ctx context.Context, products []yanolja.Product) (err error)

	//InsertOneProduct insert single product to the product database
	InsertOneProduct(ctx context.Context, product yanolja.Product) (err error)

	// UpdateProductsByProductID
	UpdateProductsByProductID(ctx context.Context, productId int64, update map[string]any) (id string, err error)

	//FetchProductByProductId
	FetchProductByProductId(ctx context.Context, productId int64) (record yanolja.Product, err error)

	//FetchProducts
	FetchProducts(ctx context.Context) (records []yanolja.Product, err error)

	//UpsertProduct
	UpsertProduct(ctx context.Context, product yanolja.Product) error

	//FetchAllProductUpdatedToday
	FetchAllProductsWithinDateRange(ctx context.Context) (products []yanolja.Product, err error)

	//DeleteProductByProductId
	DeleteProductByProductId(ctx context.Context, productId int64) (id string, err error)

	//UpdateOrderDueToRefusalToCancel
	UpdateOrderDueToRefusalToCancel(ctx context.Context, orderid int64, partnerOrderId string, orderVariantId int64, update map[string]any) (id string, err error)

	//GetOrderStatusLookup
	GetOrderStatusLookup(ctx context.Context, orderid int64, partnerorderid string, ordervariantid int64) (record []yanolja.OrderVariant, err error)

	//GetProductId
	GetProductIdFromOrder(ctx context.Context, orderid int64, partnerorderid string, ordervariantid int64) (productId int64, err error)

	//UpdateOrderVoucherIndividually
	UpdateOrderVoucherIndividually(ctx context.Context, orderid int64, partnerOrderId string, orderVariantId, orderVariantItemId int64, update map[string]any) (err error)

	//ForcedCancellationReasonUpdate
	ForcedCancellationReasonUpdate(ctx context.Context, orderid int64, partnerOrderId string, orderVariantId int64, forceCancelTypeCode string) (err error)

	//UpdateReconcilationDetail
	//UpdateReconciliationDetailByDay(ctx context.Context, req, update map[string]any) (err error)
	UpdateReconciliationDetailByDay(ctx context.Context, req map[string]any, updates []yanolja.ReconcilationDetail) error

	//UpdateReconciliationDetailByDayInsert
	UpdateReconciliationDetailByDayInsert(ctx context.Context, req map[string]any, updates []yanolja.ReconcilationDetail) error
	//GetReconciliationDetailsByDateAndStatus
	GetReconciliationDetailsByDateAndStatus(ctx context.Context, reconciliationDate string, statusCode string) (results []req_resp.OrderReconcilation, err error)

	//GetReconciliationDetailsByOrderAndVariant
	GetReconciliationDetailsByOrderAndVariant(ctx context.Context, orderId int64, orderVariantId int64, variantId int64, productId int64) ([]yanolja.ReconcilationDetail, error)

	//UpdateCancelDetailsForVariants
	UpdateCancelDetailsForVariants(ctx context.Context, orderId int64, productId int64, orderVariantId int64, cancelFailReasonCode string, cancelStatusCode string) error

	//UpdateProcessingRestoringOfOrder
	UpdateProcessingRestoringOfOrder(ctx context.Context, orderId int64, updaterec map[string]any) (err error)

	//UpdateForcedOrderDetail
	UpdateForcedCancelOrderDetail(ctx context.Context, partnerOrderId string, newStatus string) (id string, err error)

	//UpdateOrderVariantStatusByOrderId
	UpdateOrderVariantStatusByOrderId(ctx context.Context, orderid int64, variantStatus string) (err error)

	//UpdateOrderCancelAck
	UpdateOrderCancelAck(ctx context.Context, orderid int64, partnerOrderId string, canceltypecode string, orderstatus string) (err error)

	//UpdateRefusalToCancelInfo
	UpdateRefusalToCancelInfo(ctx context.Context, orderid int64, partnerOrderId string, ovariantId int64, cancelRejectTypeCode, message string) (err error)

	//UpdateForceCancelVariants
	UpdateForceCancelVariants(ctx context.Context, orderid int64, partnerOrderId string, forceCancelVariants []req_resp.CancelledVariants) (err error)

	//FindOrderByOrderIdAndPartnerOrderId
	FindOrderByOrderIdAndPartnerOrderId(ctx context.Context, orederId int64, partnerOrderId string) (record yanolja.Model, err error)
	//------------------------------------------------------------------------------------------------------------------------------
	//updateOrInsertProductView
	UpdateOrInsertProductView(ctx context.Context, products []yanolja.Product) error

	//GetRecentOrders
	GetRecentOrders(ctx context.Context) ([]yanolja.Model, error)

	//GetRecentProducts
	GetRecentProducts(ctx context.Context) ([]yanolja.ProductView, error)

	//UpsertTripSequencdId
	UpsertTripSequencdId(ctx context.Context, req yanolja.SequenceIdDetail) error

	//UpsertItemIdDetails
	UpsertItemIdDetails(ctx context.Context, itemIdDetails yanolja.ItemIdDetails) error

	//GetAllItemIdDetail
	GetAllItemIdDetail(ctx context.Context) (itemIdDetail []yanolja.ItemIdDetails, err error)

	//FetchProductImagesForLast12Hours
	FetchProductImagesForLast12Hours(ctx context.Context) (results []req_resp.ProductImages, err error)

	//BulkInsertProductImagesForProcessing
	BulkInsertProductImagesUrl(ctx context.Context, images []yanolja.ImageUrlForProcessing) error

	//BulkUpdateProductImageStatusAndImageId
	BulkUpdateProductImageStatusAndImageId(ctx context.Context, images []yanolja.ImageUrlForProcessing) error

	//GetUnsyncedImagesForTrip
	GetUnsyncedImagesForTrip(ctx context.Context) (images []trip.ImageSyncToTrip, err error)

	//BulkUpdateImageSyncStatus
	BulkUpdateImageSyncStatus(ctx context.Context, images []trip.ImageSyncToTrip) error

	//FetchProductsUpdatedToday get all the product updated today
	FetchProductsWithContentScheduleStatusFalse(ctx context.Context) ([]yanolja.Product, error)

	//GetTripsByCategory  get the trip category based on yanolja_category_code and yanolja_category_level
	GetTripsByCategory(ctx context.Context, filters []trip.CategoryFilter) ([]string, error)

	//FetchAllPluHashes  use for syncing the plu to redis using scheduler
	FetchAllPluHashes(ctx context.Context) (map[string]string, error)

	//FetchTripImageIdsByProductID  TripImageId for syncking to content api
	FetchTripImageIdsByProductID(ctx context.Context, productId int64) ([]string, error)

	//BulkUpsertPackageContent
	BulkUpsertPackageContent(ctx context.Context, packageContents []trip_domain.PackageContent) error

	//BulkUpsertProductContent
	BulkUpsertProductContent(ctx context.Context, productContents []trip_domain.ProuctContent) error

	//GetProductContentNotSync
	GetProductContentNotSync(ctx context.Context) ([]trip_domain.ProuctContent, error)

	//GetPackageContentNotSync
	GetPackageContentNotSync(ctx context.Context) ([]trip_domain.PackageContent, error)

	//BulkUpdateSyncStatus
	BulkUpdateSyncStatus(ctx context.Context, updates trip.TripMessageForSync) error

	//BulkUpsertGooglePlaceIdOfProduct
	BulkUpsertGooglePlaceIdOfProduct(ctx context.Context, docs []trip_domain.GooglePlaceIdOfProduct) error

	//InsertCallBackDetail
	InsertCallBackDetail(ctx context.Context, callbackDetail trip_domain.CallBackDetail) (id string, err error)

	//UpdateCallBackStatus
	UpdateCallBackStatus(ctx context.Context, id, status string) error

	//FetchPluByKey
	FetchPluByKey(ctx context.Context, Key string) (string, error)

	//GetDocumentID
	GetDocumentID(ctx context.Context) (primitive.ObjectID, error)

	//FindPluHashValue
	FindPluHashValue(ctx context.Context, key string) (string, error)

	//UpsertAllPlu
	UpsertAllPlu(ctx context.Context, pluHash map[string]string) error

	//FetchPluHashesByProductID
	FetchPluHashesByProductID(ctx context.Context, productId int64) (map[string]string, error)

	//BulkUpsertCategoryMapping
	BulkUpsertCategoryMapping(ctx context.Context, records []map[string]any) (map[string]int64, error)

	//DeleteAllIfNotEmpty
	DeleteAllIfNotEmpty(ctx context.Context) (map[string]int64, error)

	//Trip
	//InsertPreorderRequestFromTrip
	InsertPreorderRequestFromTrip(ctx context.Context, preorder trip.PreorderRequest) (err error)

	//InsertPaymentRequestFromTrip
	InsertPaymentRequestFromTrip(ctx context.Context, payment trip.PreOrderPaymentRequest) (err error)

	//InsertFullCancelOrderRequestFromTrip
	InsertFullCancelOrderRequestFromTrip(ctx context.Context, cancelRequest trip.CancellationRequest) (err error)

	//FetchTripRequests
	FetchTripRequests(ctx context.Context) (triporders [][]trip.ResponseForTripRequest, err error)

	//GetOrdersByOdooSyncStatus
	GetOrdersByOdooSyncStatus(ctx context.Context) ([]yanolja.Model, error)

	// FetchProductsBasedOnOdooSyncStatus
	FetchProductsBasedOnOdooSyncStatus(ctx context.Context) ([]yanolja.Product, error)

	// GetSequenceIDByOtaOrderIDAndRequestCategoryForPreorder
	GetSequenceIDByOtaOrderIDAndRequestCategory(ctx context.Context, otaOrderID, serviceName string) (bool, error)

	// UpsertOdooOrder
	UpsertOdooOrder(ctx context.Context, orders []odoo.Order) ([]odoo.Order, error)

	// UpsertOdooProduct
	UpsertOdooProduct(ctx context.Context, products []odoo.Product) ([]odoo.Product, error)

	// GetSequenceIDByKey
	GetSequenceIDByKey(ctx context.Context, key string, typeOfcollection string) (string, error)

	// GetOrdersByChannelCodeAndCustomerEmail
	GetOrdersByChannelCodeAndCustomerEmail(ctx context.Context, channelCode string, cutomerEmail string) ([]yanolja.Model, error)

	// GetOdooOrderbyOrderId
	GetOdooOrderbyOrderId(ctx context.Context, orderid int64) (record odoo.Order, err error)

	// GetPLUDetails
	GetPLUDetails(ctx context.Context, productID int64) (productview yanolja.ProductView, err error)

	// GetProductIDVersionAndSalePeriod
	GetProductIDVersionAndSalePeriod(ctx context.Context) (data []common.ProductValidityAndVersion, err error)

	//Travolution

	//InsertTravolutionProduct
	UpsertTravolutionProduct(ctx context.Context, product travolution.RawProduct) (id string, err error)
	// UpdateOdooSyncStatusToFalseIfTrue
	UpdateOdooSyncStatusToFalseIfTrue(ctx context.Context) error

	// GetOrderByPartnerIdSuffix
	GetOrderByPartnerIdSuffix(ctx context.Context, suffix string) (yanolja.Model, error)

	// GetProductByProductId
	GetProductViewByProductId(ctx context.Context, productId int64) (product yanolja.ProductView, err error)

	// GetProductByProductId
	GetAllProductViews(ctx context.Context) (products []yanolja.ProductView, err error)

	// GetProductNameByProductID
	GetProductNameByProductID(ctx context.Context, productID int64) (string, error)

	// GetProductByProductUid
	GetProductByProductUid(ctx context.Context, productUid int) (product travolution_domain.Product, err error)

	// UpsertOrder
	UpsertTravolutionOrder(ctx context.Context, payload interface{}, requestType string, status string) (string, error)

	// UpdateTravolutionOrderById
	UpdateTravolutionOrderById(ctx context.Context, update map[string]any) (id string, err error)

	// GetOrderByOrderNumber
	GetOrderByOrderNumber(ctx context.Context, orderNumber string) (order travolution_domain.Order, err error)

	// InsertWebhook
	UpsertTravolutionWebhook(ctx context.Context, payload travolution_domain.Webhook) (string, error)

	// UpsertWebhookToOrder
	UpsertWebhookToOrder(ctx context.Context, payload travolution_domain.Webhook, update map[string]any) (string, error)

	// UpdateOrderByOrderNumber
	UpdateOrderByOrderNumber(ctx context.Context, orderNumber string, update map[string]any) error
}
