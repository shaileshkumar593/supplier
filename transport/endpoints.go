package transport

import (
	"context"
	"fmt"

	svc "swallow-supplier/iface"
	travolution_domain "swallow-supplier/mongo/domain/travolution"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/travolution"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints holds all Go kit endpoints for the service.
type Endpoints struct {
	HeartBeat endpoint.Endpoint

	// Yanolja
	GetProducts                    endpoint.Endpoint
	GetProductsById                endpoint.Endpoint
	GetProductsOptionGroups        endpoint.Endpoint
	GetProductsInventories         endpoint.Endpoint
	GetVariantInventory            endpoint.Endpoint
	GetCategories                  endpoint.Endpoint
	GetRegions                     endpoint.Endpoint
	PostWaitForOrder               endpoint.Endpoint
	PostOrderConfirmation          endpoint.Endpoint
	GetOrderByOrderId              endpoint.Endpoint
	PostFullOrderCancel            endpoint.Endpoint
	PostCancelOrderByTimeOut       endpoint.Endpoint
	PostCancelOrderByForcely       endpoint.Endpoint
	PostProductCreation            endpoint.Endpoint
	GetProductByProductId          endpoint.Endpoint
	PostAllProducts                endpoint.Endpoint
	GetOrderReconcilationDetail    endpoint.Endpoint
	CancellationAckClbk            endpoint.Endpoint
	RefusalToCancelClbk            endpoint.Endpoint
	GetOrderStatusLookupClbk       endpoint.Endpoint
	ForcedOrderCancellationClbk    endpoint.Endpoint
	IndividualVoucherUpdateClbk    endpoint.Endpoint
	CombinedVoucherUpdateClbk      endpoint.Endpoint
	ProcessingOrRestoringClbk      endpoint.Endpoint
	PostTranslateImages            endpoint.Endpoint
	PostRequestFromGGT             endpoint.Endpoint
	PostCategories                 endpoint.Endpoint
	PostRegions                    endpoint.Endpoint
	GetRedisData                   endpoint.Endpoint
	GetOrderForOdoo                endpoint.Endpoint
	GetProductViewForOdoo          endpoint.Endpoint
	GetGGTInventorySync            endpoint.Endpoint
	PostUpdateImageSyncStatus      endpoint.Endpoint
	GetImageSyncToTrip             endpoint.Endpoint
	UpdateTripImageSyncStatus      endpoint.Endpoint
	GetProductContentData          endpoint.Endpoint
	GetPackageContentData          endpoint.Endpoint
	PutContentSyncData             endpoint.Endpoint
	PostContent                    endpoint.Endpoint
	GetPluToRedis                  endpoint.Endpoint
	UpsertCategoryMapping          endpoint.Endpoint
	DeleteIfNotEmpty               endpoint.Endpoint
	MonitorProductUpdate           endpoint.Endpoint
	PluUpsertToRedis               endpoint.Endpoint
	DeleteKeyFromRedis             endpoint.Endpoint
	FindRedisKeyVal                endpoint.Endpoint
	UpdatePluToRedis               endpoint.Endpoint
	FindPluFromRedis               endpoint.Endpoint
	GetOdooByOrderId               endpoint.Endpoint
	PostUpdateOdooSyncFlag         endpoint.Endpoint
	GetEverlandOrders              endpoint.Endpoint
	GetOrderByPartnerOrderIdSuffix endpoint.Endpoint
	PostUpdateOdooSyncFlagFalse    endpoint.Endpoint

	// Travolution
	GetAllproducts                endpoint.Endpoint
	GetProductByUid               endpoint.Endpoint
	GetAllOptionsOfProduct        endpoint.Endpoint
	GetOptionOfProductByOptionUid endpoint.Endpoint
	GetUnitSPriceByOptionUid      endpoint.Endpoint
	GetUnitPriceByOptionUid       endpoint.Endpoint
	GetBookingSchedules           endpoint.Endpoint
	GetAdditionalInfos            endpoint.Endpoint
	GetAdditionalInfoByUid        endpoint.Endpoint
	PostCreateAllProduct          endpoint.Endpoint
	GetProductsFromGGT            endpoint.Endpoint
	GetProductByIdFromGGT         endpoint.Endpoint
	PostCreateTravolutionOrder    endpoint.Endpoint
	GetSearchTravolutionOrder     endpoint.Endpoint
	PostCancelTravolutionOrder    endpoint.Endpoint
	PostTravolutionWebHook        endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the boilerplate.
func MakeEndpoints(s svc.Service) Endpoints {
	return Endpoints{
		HeartBeat: makeHeartBeatEndpoint(s),

		GetProducts:                    makeGetProductsEndpoints(s),
		GetProductsById:                makeGetProductsByIdEndpoints(s),
		GetProductsOptionGroups:        makeGetProductsOptionGroupsEndpoints(s),
		GetProductsInventories:         makeGetProductsInventoriesEndpoints(s),
		GetVariantInventory:            makeGetVariantInventoryEndpoints(s),
		GetCategories:                  makeGetCategoriesEndpoints(s),
		GetRegions:                     makeGetRegionsEndpoints(s),
		PostWaitForOrder:               makePostWaitForOrderEndpoints(s),
		PostOrderConfirmation:          makePostOrderConfirmationEndpoints(s),
		GetOrderByOrderId:              makeGetOrderByOrderIdEndpoints(s),
		PostFullOrderCancel:            makePostFullOrderCancelEndpoints(s),
		PostCancelOrderByTimeOut:       makePostCancelOrderByTimeOutEndpoints(s),
		PostCancelOrderByForcely:       makePostCancelOrderByForcelyEndpoints(s),
		PostProductCreation:            makePostProductCreationEndpoints(s),
		GetProductByProductId:          makeGetProductByProductIdEndpoints(s),
		PostAllProducts:                makePostAllProductsEndpoints(s),
		GetOrderReconcilationDetail:    makeGetOrderReconcilationDetailEndpoints(s),
		CancellationAckClbk:            makeCancellationAckClbkEndpoints(s),
		RefusalToCancelClbk:            makeRefusalToCancelClbkEndpoints(s),
		GetOrderStatusLookupClbk:       makeGetOrderStatusLookupClbkEndpoints(s),
		ForcedOrderCancellationClbk:    makeForcedOrderCancellationClbkEndpoints(s),
		IndividualVoucherUpdateClbk:    makeIndividualVoucherUpdateClbkEndpoints(s),
		CombinedVoucherUpdateClbk:      makeCombinedVoucherUpdateClbkEndpoints(s),
		ProcessingOrRestoringClbk:      makeProcessingOrRestoringClbkEndpoints(s),
		PostRequestFromGGT:             makePostRequestFromGGTEndpoints(s),
		PostCategories:                 makePostCategoriesEndpoints(s),
		PostRegions:                    makePPostRegionsEndpoints(s),
		GetRedisData:                   makeGetRedisDataEndpoints(s),
		GetOrderForOdoo:                makeGetOrderForOdooEndpoints(s),
		GetProductViewForOdoo:          makeGetProductViewForOdooEndpoints(s),
		GetGGTInventorySync:            makeGetGGTInventorySyncEndpoints(s),
		PostUpdateImageSyncStatus:      makePostUpdateImageSyncStatusEndpoints(s),
		GetImageSyncToTrip:             makeGetImageSyncToTripEndPoints(s),
		UpdateTripImageSyncStatus:      makeUpdateTripImageSyncStatusEndpoints(s),
		GetProductContentData:          makeGetProductContentDataEndpoints(s),
		GetPackageContentData:          makeGetPackageContentDataEndpoints(s),
		PutContentSyncData:             makePutContentSyncDataEndpoints(s),
		PostContent:                    makePostContentEndpoints(s),
		GetPluToRedis:                  makeGetPluToRedisEndpoints(s),
		UpsertCategoryMapping:          makeUpsertCategoryMappingEndpoints(s),
		DeleteIfNotEmpty:               makeDeleteIfNotEmptyEndpoints(s),
		MonitorProductUpdate:           makeMonitorProductUpdateEndpoints(s),
		PluUpsertToRedis:               makePluUpsertToRedisEndpoints(s),
		DeleteKeyFromRedis:             makeDeleteKeyFromRedisEndpoints(s),
		FindRedisKeyVal:                makeFindRedisKeyValEndpoints(s),
		UpdatePluToRedis:               makeUpdatePluToRedisEndpoints(s),
		FindPluFromRedis:               makeFindPluFromRedisEndpoints(s),
		GetOdooByOrderId:               makeGetOdooByOrderIdEndpoints(s),
		PostUpdateOdooSyncFlag:         makePostUpdateOdooSyncFlagEndpoints(s),
		GetEverlandOrders:              makeGetEverlandOrdersEndpoints(s),
		PostUpdateOdooSyncFlagFalse:    makePostUpdateOdooSyncFlagFalseEndpoints(s),
		GetOrderByPartnerOrderIdSuffix: makeGetOrderByPartnerOrderIdSuffixEndpoints(s),

		// Travolution
		GetAllproducts:                makeGetAllproductsEndpoints(s),
		GetProductByUid:               makeGetProductByUidEndpoints(s),
		GetAllOptionsOfProduct:        makeGetAllOptionsOfProductEndpoints(s),
		GetOptionOfProductByOptionUid: makeGetOptionOfProductByOptionUidEndpoints(s),
		GetUnitSPriceByOptionUid:      makeGetUnitSPriceByOptionUidEndpoints(s),
		GetUnitPriceByOptionUid:       makeGetUnitPriceByOptionUidEndpoints(s),
		GetBookingSchedules:           makeGetBookingSchedulesEndpoints(s),
		GetAdditionalInfos:            makeGetAdditionalInfosEndpoints(s),
		GetAdditionalInfoByUid:        makeGetAdditionalInfoByUidEndpoints(s),
		PostCreateAllProduct:          makePostCreateAllProductEndpoints(s),
		GetProductsFromGGT:            makeGetProductsFromGGTEndpoints(s),
		GetProductByIdFromGGT:         makeGetProductByIdFromGGTEndpoints(s),
		PostCreateTravolutionOrder:    makePostCreateTravolutionOrderEndpoints(s),
		GetSearchTravolutionOrder:     makeGetSearchTravolutionOrderEndpoint(s),
		PostCancelTravolutionOrder:    makePostCancelTravolutionOrderEndpoint(s),
		PostTravolutionWebHook:        makePostTravolutionWebHookEndpoint(s),
	}

}

func makeHeartBeatEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.HeartBeat(ctx)
		return res, err
	}
}

// ------------------------------------------------------------------------------------------------------------------
// Yanolja all endpoint listed here
func makeGetProductsEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.AllProduct)
		res, err := s.GetProducts(ctx, req)
		return res, err
	}
}

func makeGetProductsByIdEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.ProductsById)
		fmt.Println("1")
		res, err := s.GetProductsById(ctx, req)
		return res, err
	}
}

func makeGetProductsOptionGroupsEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.ProductsById)
		res, err := s.GetProductsOptionGroups(ctx, req)
		return res, err
	}
}

func makeGetProductsInventoriesEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.ProductInventory)
		res, err := s.GetProductsInventories(ctx, req)
		return res, err
	}
}

func makeGetVariantInventoryEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.VariantInventory)
		res, err := s.GetVariantInventory(ctx, req)
		return res, err
	}
}

func makeGetCategoriesEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetCategories(ctx)
		return res, err
	}
}

func makeGetRegionsEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetRegions(ctx)
		return res, err
	}
}

func makePostWaitForOrderEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.WaitingForOrder)
		res, err := s.PostWaitForOrder(ctx, req)
		return res, err
	}
}

func makePostOrderConfirmationEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.OrderConfirmation)
		res, err := s.PostOrderCompletion(ctx, req)
		return res, err
	}
}

func makeGetOrderByOrderIdEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.OrderConfirmation)
		res, err := s.GetOrderByOrderId(ctx, req)
		return res, err
	}
}
func makePostFullOrderCancelEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.Order)
		res, err := s.PostCancelOrderEntirly(ctx, req.OrderId)
		return res, err
	}

}
func makePostCancelOrderByTimeOutEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.PartnerOrder)
		res, err := s.PostCancelOrderByReqTimeOut(ctx, req.PartnerOrderId)
		return res, err
	}
}
func makePostCancelOrderByForcelyEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.PartnerOrder)
		res, err := s.PostForcelyCancelOrder(ctx, req.PartnerOrderId)
		return res, err
	}
}

func makePostProductCreationEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.Upsert_Product)
		res, err := s.InsertProductClbk(ctx, req)
		return res, err
	}
}

func makeGetProductByProductIdEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.GetProduct)
		res, err := s.GetProductByProductId(ctx, req.ProductId)
		return res, err
	}
}

func makePostAllProductsEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.InsertAllProduct(ctx)
		return res, err
	}
}

func makeGetOrderReconcilationDetailEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.OrderReconcileReq)
		res, err := s.GetOrderReconcilationDetail(ctx, req)
		return res, err
	}
}

func makeCancellationAckClbkEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.CancellationAck)
		res, err := s.CancellationAckClbk(ctx, req)
		return res, err
	}
}

func makeRefusalToCancelClbkEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.RefusalToCancel)
		res, err := s.RefusalToCancelClbk(ctx, req)
		return res, err
	}
}

func makeGetOrderStatusLookupClbkEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.OrderStatusLookup)
		res, err := s.OrderStatusLookupClbk(ctx, req)
		return res, err
	}
}

func makeForcedOrderCancellationClbkEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.ForcedOrderCancellation)
		res, err := s.ForcedOrderCancellationClbk(ctx, req)
		return res, err
	}
}

func makeIndividualVoucherUpdateClbkEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.IndividualVoucherUpdate)
		res, err := s.IndividualVoucherUpdateClbk(ctx, req)
		return res, err
	}
}

func makeCombinedVoucherUpdateClbkEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.CombinedVoucherUpdate)
		res, err := s.CombinedVoucherUpdateClbk(ctx, req)
		return res, err
	}
}

func makeProcessingOrRestoringClbkEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.ProcessingOrRestoringReq)
		res, err := s.ProcessingOrRestoringClbk(ctx, req)
		return res, err
	}
}

func makePostRequestFromGGTEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(trip.SwallowRequest)
		res, err := s.PostRequestFromGGT(ctx, req)
		return res, err
	}
}

func makePostCategoriesEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.InsertAllCategories(ctx)
		return res, err
	}
}

func makePPostRegionsEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.InsertAllRegions(ctx)
		return res, err
	}
}

func makeGetRedisDataEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetRedisData(ctx)
		return res, err
	}
}

func makeGetOrderForOdooEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetOrderSync(ctx)
		return res, err
	}
}

func makeGetProductViewForOdooEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetProductSync(ctx)
		return res, err
	}
}

func makeGetGGTInventorySyncEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.InventorySync(ctx)
		return res, err
	}
}

func makePostUpdateImageSyncStatusEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.([]domain.ImageUrlForProcessing)
		res, err := s.UpdateImageSyncStatus(ctx, req)
		return res, err
	}
}

func makeGetImageSyncToTripEndPoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetImageSyncToTrip(ctx)
		return res, err
	}
}
func makeUpdateTripImageSyncStatusEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.([]trip.ImageSyncToTrip)
		res, err := s.UpdateTripImageSyncStatus(ctx, req)
		return res, err
	}
}

func makeGetProductContentDataEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.ProductContentSync(ctx)
		return res, err
	}
}
func makeGetPackageContentDataEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.PackageContentSync(ctx)
		return res, err
	}
}

func makePutContentSyncDataEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(trip.TripMessageForSync)
		res, err := s.UpdateContentSyncStatus(ctx, req)
		return res, err
	}
}

func makePostContentEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.ProductSyncToTrip(ctx)
		return res, err
	}
}

func makeGetPluToRedisEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetPluToRedis(ctx)
		return res, err
	}
}

func makeUpsertCategoryMappingEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.UpsertCategoryMapping(ctx)
		return res, err
	}
}

func makeDeleteIfNotEmptyEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.DeleteRecoRdIfNotEmpty(ctx)
		return res, err
	}
}

func makeMonitorProductUpdateEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.MonitorProductUpdateSvc(ctx)
		return res, err
	}
}

func makePluUpsertToRedisEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.SyncAllPluToRedis(ctx)
		return res, err
	}
}

func makeDeleteKeyFromRedisEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.DeleteRedisData(ctx)
		return res, err
	}
}

func makeFindRedisKeyValEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(common.RedisKey)
		res, err := s.FindRedisKeyValue(ctx, req.Key)
		return res, err
	}
}

func makeUpdatePluToRedisEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.UpdatePluToRedis(ctx)
		return res, err
	}
}

func makeFindPluFromRedisEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.PluRequest)
		res, err := s.GetPluFromRedis(ctx, req)
		return res, err
	}
}

func makeGetOdooByOrderIdEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.Order)
		res, err := s.GetOrderSyncByOrderIdInOdoo(ctx, req)
		return res, err
	}

}

func makePostUpdateOdooSyncFlagEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.OrderList)
		res, err := s.PostOdooSyncStatus(ctx, req)
		return res, err
	}

}

func makeGetEverlandOrdersEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(yanolja.EverlandGetRequest)
		res, err := s.GetEverlandOrders(ctx, req)
		return res, err
	}

}

// Travolution
func makeGetAllproductsEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.ProductReq)
		res, err := s.GetAllproducts(ctx, req)
		return res, err
	}

}

func makeGetProductByUidEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.ProductReq)
		res, err := s.GetProductByUid(ctx, req)
		return res, err
	}

}

func makeGetAllOptionsOfProductEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.OptionRequest)
		res, err := s.GetAllOptionsOfProduct(ctx, req)
		return res, err
	}

}

func makeGetOptionOfProductByOptionUidEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.OptionRequest)
		res, err := s.GetOptionOfProductByOptionUid(ctx, req)
		return res, err
	}

}

func makeGetUnitSPriceByOptionUidEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.UnitPriceRequest)
		res, err := s.GetUnitSPriceByOptionUid(ctx, req)
		return res, err
	}

}

func makeGetUnitPriceByOptionUidEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.UnitPriceRequest)
		res, err := s.GetUnitPriceByOptionUid(ctx, req)
		return res, err
	}

}

func makeGetBookingSchedulesEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.BookingScheduleReq)
		res, err := s.GetBookingSchedules(ctx, req)
		return res, err
	}

}

func makeGetAdditionalInfosEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.BookingAdditionalInfoRequest)
		res, err := s.GetAdditionalInfos(ctx, req)
		return res, err
	}

}

func makeGetAdditionalInfoByUidEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.BookingAdditionalInfoRequest)
		res, err := s.GetAdditionalInfos(ctx, req)
		return res, err
	}

}

func makePostCreateAllProductEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.PostCreateAllProduct(ctx)
		return res, err
	}
}
func makePostUpdateOdooSyncFlagFalseEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.PostChangeOdooSyncStatusToFalse(ctx)
		return res, err
	}
}

func makeGetOrderByPartnerOrderIdSuffixEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		res, err := s.GetOderByPartialPartnerOrderIdSuffix(ctx, req)
		return res, err
	}

}

func makeGetProductsFromGGTEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res, err := s.GetProductsFromGGT(ctx)
		return res, err
	}

}

func makeGetProductByIdFromGGTEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		productId := request.(int64)
		res, err := s.GetProductByIdFromGGT(ctx, productId)
		return res, err
	}

}

func makePostCreateTravolutionOrderEndpoints(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution.OrderRequest)
		res, err := s.CreateTravolutionOrder(ctx, req)
		return res, err
	}
}

func makeGetSearchTravolutionOrderEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		res, err := s.SearchTravolutionOrder(ctx, req)
		return res, err
	}
}

func makePostCancelTravolutionOrderEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		res, err := s.CancelTravolutionOrder(ctx, req)
		return res, err
	}
}

func makePostTravolutionWebHookEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(travolution_domain.Webhook)
		res, err := s.OrderWebhookUpdate(ctx, req)
		return res, err
	}
}
