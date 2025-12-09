package implementation

import (
	"fmt"
	"net/http"
	"strings"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	trip_domain "swallow-supplier/mongo/domain/trip"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/yanolja"
	tripservice "swallow-supplier/services/distributors/trip"
	"swallow-supplier/utils/constant"

	"swallow-supplier/utils"

	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// CancellationAck
func (s *service) CancellationAckClbk(ctx context.Context, ackreq yanolja.CancellationAck) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "CancellationAckClbk",
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

	record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, ackreq.OrderId)
	if err != nil {
		resp.Code = "500"
		level.Error(logger).Log("repository_error", "GetOrderbyOrderId  throws error ", err)
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyOrderId")
	}

	for _, variant := range record.OrderVariants {
		err = UpdateReconcilationDetail(ctx, s, logger, ackreq.OrderId, ackreq.PartnerOrderId, variant)
		if err != nil {
			resp.Code = "500"
			level.Error(logger).Log("error", "UpdateReconcilationDetail  throws error ", err)
			return resp, customError.NewError(ctx, "leisure-api-1006", fmt.Sprintf("error on reconcilation update, %v", err), "UpdateReconcilationDetail")
		}
	}

	err = s.mongoRepository[config.Instance().MongoDBName].UpdateOrderCancelAck(ctx, ackreq.OrderId, ackreq.PartnerOrderId, ackreq.OrderCancelTypeCode, "CANCELED")
	if err != nil {
		resp.Code = "500"
		level.Error(logger).Log("repository_error", "UpdateOrderCancelAck  throws error ", err)
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on update order by orderId and partnerOrderId, %v", err), "CancellationAckClbk")
	}
	level.Info(logger).Log("info", "document update with orderId ", ackreq.OrderId)

	record, err = s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, ackreq.OrderId)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "yanolja")
	}

	level.Info(logger).Log("info", "updated document with OrderId ", ackreq.OrderId)

	processTripNotification(ctx, logger, s, "CancelAckNotify", record, constant.TRIPFULLORDERCANCELREQUEST)

	resp.Body = record
	resp.Code = "200"

	return resp, nil
}

// RefusalToCancel
func (s *service) RefusalToCancelClbk(ctx context.Context, refusaltocancel yanolja.RefusalToCancel) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "RefusalToCancelClbk",
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

	err = s.mongoRepository[config.Instance().MongoDBName].UpdateRefusalToCancelInfo(ctx, refusaltocancel.OrderId, refusaltocancel.PartnerOrderId, refusaltocancel.OrderVariantID, refusaltocancel.CancelRejectTypeCode, refusaltocancel.Message)
	if err != nil {
		resp.Code = "500"
		level.Error(logger).Log("repository_error", "UpdateRefusalToCancelInfo  throws error ", err)
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on update order by orderId and partnerOrderId, %v", err), "RefusalToCancelClbk")
	}
	level.Info(logger).Log("info", "document update with orderId ", refusaltocancel.OrderId)

	record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, refusaltocancel.OrderId)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "yanolja")
	}
	level.Info(logger).Log("info", "updated document with OrderId ", refusaltocancel.OrderId)

	if IsAllVariantCancilationDone(record) {
		processTripNotification(ctx, logger, s, "CancelAckNotify", record, constant.TRIPFULLORDERCANCELREQUEST)
	}

	resp.Body = record
	resp.Code = "200"

	return resp, nil
}

// OrderStatusLookupClbk
func (s *service) OrderStatusLookupClbk(ctx context.Context, lookup yanolja.OrderStatusLookup) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "OrderStatusLookupClbk",
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

	record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderStatusLookup(ctx, lookup.OrderId,
		lookup.PartnerOrderId, lookup.OrderVariantID)
	if err != nil {
		resp.Code = "500"
		level.Error(logger).Log("repository_error", "GetOrderStatusLookup  throws error ", err)
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order status lookup, %v", err), "GetOrderStatusLookup")
	}

	level.Info(logger).Log("info", "document with orderId is  ", lookup.OrderId)

	var lookupresp yanolja.OrderStatusLookupResp
	var ordervarnt = make([]yanolja.OrderVariants, len(record))
	for i := 0; i < len(record); i++ {
		ordervarnt[i].OrderVariantID = record[i].OrderVariantID
		ordervarnt[i].OrderVariantStatusTypeCode = record[i].OrderVariantStatusTypeCode
	}

	lookupresp.OrderId = lookup.OrderId
	lookupresp.PartnerOrderId = lookup.PartnerOrderId
	lookupresp.OrderVariants = ordervarnt

	resp.Body = lookupresp
	resp.Code = "200"
	return resp, nil
}

// ForcedOrderCancellationClbk is used when customer requests a cancellation directly with Yanolja
// (through customer service, etc.), or when a cancellation request is made by the facility to Yanolja.
func (s *service) ForcedOrderCancellationClbk(ctx context.Context, forcecancellation yanolja.ForcedOrderCancellation) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "ForcedOrderCancellationClbk",
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

	record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, forcecancellation.OrderId)
	if err != nil {
		resp.Code = "500"
		level.Error(logger).Log("repository_error", "GetOrderbyOrderId  throws error ", err)
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyOrderId")
	}

	var incrmnt int8 = 0
	for _, ordervariant := range record.OrderVariants {
		if ordervariant.OrderVariantID == forcecancellation.ForceCancelVariants[incrmnt].OrderVariantID {
			err = UpdateReconcilationDetail(ctx, s, logger, forcecancellation.OrderId, forcecancellation.PartnerOrderId, ordervariant)
			if err != nil {
				resp.Code = "500"
				level.Error(logger).Log("error", "UpdateReconcilationDetail  throws error ", err)
				return resp, customError.NewError(ctx, "leisure-api-1006", fmt.Sprintf("error on reconcilation update, %v", err), "UpdateReconcilationDetail")
			}

			err = s.mongoRepository[config.Instance().MongoDBName].ForcedCancellationReasonUpdate(ctx, forcecancellation.OrderId,
				forcecancellation.PartnerOrderId, forcecancellation.ForceCancelVariants[incrmnt].OrderVariantID,
				forcecancellation.ForceCancelVariants[incrmnt].ForceCancelTypeCode)
			if err != nil {
				resp.Code = "500"
				level.Error(logger).Log("repository_error", "ForcedCancellationReasonUpdate  throws error ", err)
				return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on update order by orderId and partnerOrderId, %v", err), "CancellationAckClbk")
			}
			incrmnt = incrmnt + 1
			if len(forcecancellation.ForceCancelVariants) == int(incrmnt) {
				break
			}
		}

	}

	level.Info(logger).Log("info", "document update with orderId ", forcecancellation.OrderId)

	record, err = s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, forcecancellation.OrderId)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "yanolja")
	}

	level.Info(logger).Log("info", "updated document with OrderId ", forcecancellation.OrderId)
	processTripNotification(ctx, logger, s, "ForcedCancelNotify", record, "")

	resp.Body = record
	resp.Code = "200"
	return
}

// IndividualVoucherUpdate  one by one voucher update
func (s *service) IndividualVoucherUpdateClbk(ctx context.Context, voucher yanolja.IndividualVoucherUpdate) (resp yanolja.Response, err error) {

	requestID := utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "IndividualVoucherUpdateClbk",
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

	productId, err := s.mongoRepository[config.Instance().MongoDBName].GetProductIdFromOrder(ctx, voucher.OrderId, voucher.PartnerOrderId, voucher.OrderVariantID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no product exist  ", err)
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "FindOrderByOrderIdAndPartnerOrderId")
		}
		level.Error(logger).Log("error", "request to yanolja client raise error from database of order ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching productId , %v", err), "GetProductId")
	}

	record, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, productId)
	if err != nil {
		level.Error(logger).Log("database error", "retriving product by productId", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching  product by  productId , %v", err), "FetchProductByProductId")
	}

	level.Info(logger).Log("info", fmt.Sprintf("IsIntegratedVoucher status is %v", record.IsIntegratedVoucher))

	var update map[string]any
	if !record.IsIntegratedVoucher {
		update = map[string]any{
			"voucherDisplayTypeCode": voucher.VoucherDisplayTypeCode,
			"voucherCode":            voucher.VoucherCode,
		}

		level.Info(logger).Log("Info ", "record to update ", " update ", update)

		err = s.mongoRepository[config.Instance().MongoDBName].UpdateOrderVoucherIndividually(ctx, voucher.OrderId, voucher.PartnerOrderId, voucher.OrderVariantID, voucher.OrderVariantItemId, update)
		if err != nil {
			level.Error(logger).Log("database error", " individual voucher update", err)
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on individual voucher update of variantItem of order, %v", err), "UpdateOrderVoucherIndividually")
		}
	}

	order, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, voucher.OrderId)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "yanolja")
	}

	level.Info(logger).Log("info", "document update with product id ", record.ProductID)
	status := IsAllVoucherDataAvailable(order)

	level.Info(logger).Log("info", "IsAllVoucherDataAvailable status  ", status)

	if status == true {
		processTripNotification(ctx, logger, s, "VoucherUpdateNotify", order, constant.TRIPPAYMENTREQUEST)
	}

	resp.Body = order
	resp.Code = "200"

	return resp, nil
}

// CombinedVoucherUpdate multiple voucher upadate in single call
func (s *service) CombinedVoucherUpdateClbk(ctx context.Context, voucher yanolja.CombinedVoucherUpdate) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "CombinedVoucherUpdateClbk",
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

	level.Info(logger).Log("Info ", "sequence Id Generation")

	record, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, voucher.ProductID)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching product based on ProductId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching product by productId, %v", err), "FetchProductByProductId")
	}

	level.Info(logger).Log("info", fmt.Sprintf("IsIntegratedVoucher status is %v", record.IsIntegratedVoucher))

	var update map[string]any
	update = map[string]any{
		"voucherDisplayTypeCode": voucher.VoucherDisplayTypeCode,
		"voucherCode":            voucher.VoucherCode,
	}

	level.Info(logger).Log("Info ", "record to update ", " update ", update)

	if record.IsIntegratedVoucher {

		for _, val := range voucher.OrderVariantIds {
			for _, itemIds := range val.OrderVariantItemIds {

				err = s.mongoRepository[config.Instance().MongoDBName].UpdateOrderVoucherIndividually(ctx, voucher.OrderId, voucher.PartnerOrderId, val.OrderVariantID, itemIds, update)
				if err != nil {
					level.Error(logger).Log("repository error", "individual voucher update", err)
					resp.Code = "500"
					return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on updating voucher, %v", err), "UpdateOrderVoucherIndividually")

				}
			}

		}

	}

	combinedupdate, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, voucher.OrderId)
	if err != nil {
		level.Error(logger).Log("error", "error due accesing order by orderId", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "yanolja")
	}

	// Process Trip Notification asynchronously
	if IsAllVoucherDataAvailable(combinedupdate) {
		processTripNotification(ctx, logger, s, "BulkVoucherUpdateNotify", combinedupdate, constant.TRIPPAYMENTREQUEST)
	}
	resp.Body = combinedupdate
	resp.Code = "200"

	return
}

// ProcessingOrRestoringClbk
func (s *service) ProcessingOrRestoringClbk(ctx context.Context, req yanolja.ProcessingOrRestoringReq) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)
	logger := log.With(
		s.logger,
		"method", "ProcessingOrRestoringClbk",
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

	if req.EventType != constant.CONSTANTCONSUME && req.EventType != constant.CONSTANTRESTORE {
		resp.Body = "eventType value is not CONSUME|RESTORE "
		resp.Body = "500"
		return resp, customError.NewError(ctx, "leisure-api-0006", fmt.Sprintf("eventType value is not CONSUME|RESTORED, %v", err), "ProcessingOrRestoringClbk")
	}
	level.Info(logger).Log("event_type : ", req.EventType)
	update := make(map[string]any)
	if req.EventType == constant.CONSTANTCONSUME {
		update = map[string]any{
			"partnerOrderId":             req.PartnerOrderId,
			"orderVariantId":             req.OrderVariantId,
			"dateTime":                   req.DateTime,
			"dateTimeTimeZone":           req.DateTimeTimezone,
			"dateTimeOffset":             req.DateTimeOffset,
			"orderVariantStatusTypeCode": constant.ORDERVARIANTUSEDSTATUS,
		}
	} else if req.EventType == constant.CONSTANTRESTORE {
		update = map[string]any{
			"partnerOrderId":             req.PartnerOrderId,
			"orderVariantId":             req.OrderVariantId,
			"dateTime":                   req.DateTime,
			"dateTimeTimeZone":           req.DateTimeTimezone,
			"dateTimeOffset":             req.DateTimeOffset,
			"orderVariantStatusTypeCode": constant.ORDERVARIANTNOTUSEDSTATUS,
		}
	}

	level.Info(logger).Log("Info ", "record to update ", " update ", update)
	err = s.mongoRepository[config.Instance().MongoDBName].UpdateProcessingRestoringOfOrder(ctx, req.OrderId, update)
	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error from database ", err)
		resp.Body = "Database error in updating used or restoring detail"
		resp.Body = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by partnerOrderId, %v", err), "UpdateProcessingRestoringOfOrder")
	}

	record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyPartnerOrderId(ctx, req.PartnerOrderId)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching order based on partnerOrderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by partnerOrderId, %v", err), "GetOrderbyPartnerOrderId")
	}

	level.Info(logger).Log("info", "Document updated with orderId", "orderId :", record.OrderId)

	level.Info(logger).Log("info", " ---- ConsumeRestoreNotify -----", "orderVariantId : ", update["orderVariantId"])
	// Process Trip Notification asynchronously
	processTripNotification(ctx, logger, s, "ConsumeRestoreNotify", record, "")

	// Immediately return response to the caller
	resp.Body = record
	resp.Code = "200"

	return resp, nil
}

// use for testing
func SayHello(logger log.Logger, record domain.Model) {
	level.Info(logger).Log("method", fmt.Sprintf("SayHello for message %s", "Hello Call "))
	fmt.Println("+++++++++++++++++++++++++++++++  Hello Trip +++++++++++++++++++++++ ", record)

}

// processTripNotification handles Trip response in the background with retries.
func processTripNotification(ctx context.Context, logger log.Logger, svc *service, message string, record domain.Model, requestType string) {
	level.Info(logger).Log("method", fmt.Sprintf("processTripNotification for message %s", message))

	var sequenceId string
	defer func() {
		if r := recover(); r != nil {
			level.Error(logger).Log("error", "Panic in processTripNotification operation", "panic", r)
		}
	}()

	tripresp, err := MakeTripResponse(ctx, logger, record)
	if err != nil {
		level.Error(logger).Log("error", "MakeTripResponse failed", "orderId", record.OrderId, "error", err)
		return
	}

	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		return
	}
	for i, item := range tripresp.Items.Items {
		itemstr, err := cacheLayer.Get(ctx, item.PLU)
		if err != nil {
			level.Error(logger).Log("cache error in getting plu ", err)
			return
		}
		tripresp.Items.Items[i].PLU = itemstr
	}

	level.Info(logger).Log("requestType ", requestType)
	// for  order api its required to read from cache
	if requestType != "" {
		level.Info(logger).Log("info", "retriving sequence ID from redis ")

		key := record.PartnerOrderID + "-" + requestType
		level.Info(logger).Log("key : ", key)

		// counter added to check how many times loop will continue
		count := 10
		for key != "" && sequenceId == "" && count != 0 {
			sequenceId, err = cacheLayer.Get(ctx, key)
			level.Info(logger).Log("key value ", sequenceId)
			if sequenceId != "" {
				break
			}
			level.Info(logger).Log("counter ", count)
			count = count - 1
		}

		if err != nil {
			level.Error(logger).Log("cache error in getting sequenceId", err)
			return
		}

		if sequenceId == "" {
			// call repository trip_payment_request
			sequenceId, _ = svc.mongoRepository[config.Instance().MongoDBName].GetSequenceIDByKey(ctx, record.PartnerOrderID, requestType)
			level.Info(logger).Log("sequenceId ", sequenceId)

		}
	} else {
		// for content api its required to generate
		level.Info(logger).Log("info", "Generating sequence ID")
		sequenceId, err = utils.GetSequenceID()
		if err != nil {
			level.Error(logger).Log("error", "Error generating sequenceId", "error", err)
			return
		}
	}

	//  to resolve cycle error
	tresp := trip_domain.ChannelReqest{
		Items: tripresp.Items,
		Order: tripresp.Order,
	}
	level.Info(logger).Log("info", "Notifying trip service")
	tripRequest := trip_domain.CallBackRequest{
		ChannelReq:    tresp,
		Message:       message,
		SequenceIdGen: sequenceId,
	}

	level.Info(logger).Log("payload sent to trip ", tripRequest)

	caser := cases.Title(language.English)

	detail := trip_domain.CallBackDetail{
		ChannelCallBackInfo: tripRequest,
		Supplier:            caser.String(strings.ToLower(record.Suppliers)),
		Distributor:         record.PartnerOrderChannelCode,
	}
	id, err := svc.mongoRepository[config.Instance().MongoDBName].InsertCallBackDetail(ctx, detail)
	if err != nil {
		level.Error(logger).Log("repository error", "callback insertion error  ", err)
		return
	}

	level.Info(logger).Log("info ", "call to trip service", record.OrderId)
	tripsvc, _ := tripservice.New(ctx)
	success := false // Track success status

	resp, err := tripsvc.NotifyToTrip(ctx, tripRequest)
	if err == nil {
		level.Info(logger).Log("info", "Successfully notified trip service", "orderId", record.OrderId)
		success = true
	}

	if !success && resp.Code != "200" {
		level.Error(logger).Log("error", "Unsuccessful trip notification after retries", "orderId", record.OrderId)
		// add repository call "FAILED"
		err = svc.mongoRepository[config.Instance().MongoDBName].UpdateCallBackStatus(ctx, id, "FAILED")
		if err != nil {
			return
		}
	} else {
		level.Info(logger).Log("info", "Trip notification was successful", "orderId", record.OrderId)
		// Add repository call for SUCCESS
		err = svc.mongoRepository[config.Instance().MongoDBName].UpdateCallBackStatus(ctx, id, "SUCCESS")
		if err != nil {
			return
		}
	}

}

// checking for all voucher available for a order
func IsAllVoucherDataAvailable(order domain.Model) bool {
	if len(order.OrderVariants) == 0 {
		return false
	}

	for _, variant := range order.OrderVariants {
		if len(variant.OrderVariantItems) == 0 {
			return false
		}

		for _, item := range variant.OrderVariantItems {
			v := item.Voucher
			if v.VoucherCode == "" || v.VoucherDisplayTypeCode == "" {
				return false
			}
		}
	}

	return true
}

// check for cancelling order
func IsAllVariantCancilationDone(order domain.Model) bool {
	if len(order.OrderVariants) == 0 {
		return false
	}

	for _, variant := range order.OrderVariants {
		if variant.CancelRejectTypeCode == "" && variant.OrderVariantStatusTypeCode == "CANCELING" {
			return false
		}
	}

	return true
}

func IsAllVariantUsed(order domain.Model) bool {
	if len(order.OrderVariants) == 0 {
		return false
	}

	for _, variant := range order.OrderVariants {
		if variant.OrderVariantStatusTypeCode != "USED" {
			return false
		}
	}

	return true
}
