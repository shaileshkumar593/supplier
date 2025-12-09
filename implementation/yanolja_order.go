package implementation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/yanolja"
	yanoljasvc "swallow-supplier/services/suppliers/yanolja"
	"swallow-supplier/utils"
	"swallow-supplier/utils/constant"
	"swallow-supplier/utils/validator"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostWaitForOrder preorder order creation from yanolja
func (s *service) PostWaitForOrder(ctx context.Context, req yanolja.WaitingForOrder) (resp yanolja.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostWaitForOrder",
		"Request ID", requestID,
		"request", req,
	)

	// Defer when panic
	defer func(context.Context) {
		r := recover()
		if r == nil {
			return
		}

		level.Error(logger).Log("error", "processing request went into panic mode", "panic", r)
		resp.Code = "500"
		err = fmt.Errorf("panic occurred: %v", r)

	}(ctx)

	id, err := s.mongoRepository[config.Instance().MongoDBName].InsertPreOrder(ctx, req)
	if err != nil {
		resp.Code = "500"
		return resp, err
	}
	level.Info(logger).Log("info", "document inserted with id ", id)

	// validate condition for the preorder creation
	err = PreOrderCreationValidate(ctx, s, req)
	if err != nil {
		level.Error(logger).Log("error ", "preordercreation validation failed")
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}

	var ordersvc, _ = yanoljasvc.New(ctx)
	resp, err = ordersvc.WaitingForOrder(ctx, req)
	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		resp.Code = "500"
		return resp, err
	}

	level.Info(logger).Log("yanolja response", resp)
	orderResp := resp.Body.(map[string]interface{})

	update := map[string]any{
		"_id":             id,
		"orderId":         orderResp["orderId"],
		"orderStatusCode": orderResp["orderStatusCode"],
		"orderVariants":   orderResp["orderVariants"],
		"oodoSyncStatus":  false,
		"updatedAt":       time.Now().UTC().Format(time.RFC3339),
	}

	id, err = s.mongoRepository[config.Instance().MongoDBName].UpdatePreOrderById(ctx, update)
	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update by orderId, %v", err), "PostWaitForOrder")
	}
	level.Info(logger).Log("info", "document update with id ", id)

	// Fetch the updated record by orderId
	orderIdRaw, ok := orderResp["orderId"] // Assuming orderId is a string, adjust if needed
	if !ok {
		resp.Code = "500"
		return resp, fmt.Errorf("failed to parse orderId from yanolja response")
	}

	ordVrnts := orderResp["orderVariants"].([]interface{})

	// Step 6: Update reconciliation detail for each selected variant
	for i, val := range req.SelectVariants {
		ordrVrnt, _ := ordVrnts[i].(domain.OrderVariant)

		// Step 6.1: Try to fetch existing reconciliation details
		reconDetails, err := s.mongoRepository[config.Instance().MongoDBName].GetReconciliationDetailsByOrderAndVariant(
			ctx, int64(orderIdRaw.(float64)), ordrVrnt.OrderVariantID, val.VariantID, val.ProductID)

		// Step 6.2: Check if no document found, then initialize new reconciliation details
		if err == mongo.ErrNoDocuments || len(reconDetails) == 0 {
			// Create new reconciliation detail if no document found for the day
			newDetail := domain.ReconcilationDetail{
				ReconciliationDate:       time.Now().UTC().Format("2006-01-02"),
				ReconcileOrderStatusCode: constant.RECONCILESTATUSUNKNOWN,
			}
			reconDetails = []domain.ReconcilationDetail{newDetail} // initialize with new entry
		} else if err != nil {
			// If any other error occurred, return it
			resp.Code = "500"
			return resp, fmt.Errorf("failed to fetch reconciliation details: %w", err)
		} else {
			// Step 6.3: Update or append reconciliation detail for the given date
			reconUpdated := false
			for i, detail := range reconDetails {
				if detail.ReconciliationDate == time.Now().UTC().Format("2006-01-02") {
					// Update existing reconciliation detail for the day
					reconDetails[i].ReconcileOrderStatusCode = constant.RECONCILESTATUSUNKNOWN
					reconDetails[i].ReconciliationDate = time.Now().UTC().Format(time.RFC3339)
					reconUpdated = true
					break
				}
			}
			// If not updated, append a new reconciliation detail for today
			if !reconUpdated {
				newDetail := domain.ReconcilationDetail{
					ReconciliationDate:       time.Now().UTC().Format("2006-01-02"),
					ReconcileOrderStatusCode: constant.RECONCILESTATUSUNKNOWN,
				}
				reconDetails = append(reconDetails, newDetail)
			}
		}

		// Step 6.4: Update reconciliation details array in MongoDB
		filter := map[string]any{
			"orderId":        orderResp["orderId"],
			"partnerOrderId": req.PartnerOrderID,
			"productId":      val.ProductID,
			"variantId":      val.VariantID,
		}
		if err := s.mongoRepository[config.Instance().MongoDBName].UpdateReconciliationDetailByDayInsert(
			ctx, filter, reconDetails); err != nil {
			level.Error(logger).Log("error", "failed to update reconciliation details", err)
			resp.Code = "500"
			return resp, fmt.Errorf("failed to update reconciliation details in database: %w", err)
		}
	}

	record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, int64(orderIdRaw.(float64))) // getting float64 type in response for OrderId
	if err != nil {
		level.Error(logger).Log("error", "failed to fetch record by orderId", err)
		resp.Code = "500"
		return resp, fmt.Errorf("failed to fetch order by orderId: %w", err)
	}
	level.Info(logger).Log("response", record)

	resp.Code = "200"
	resp.Body = record
	return resp, nil

}

// PostOrderCompletion confirm the order creation
func (s *service) PostOrderCompletion(ctx context.Context, req yanolja.OrderConfirmation) (resp yanolja.Response, err error) {

	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostOrderCompletion",
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
	resp, err = getsvc.OrderComplete(ctx, req)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		return resp, err
	}

	level.Info(logger).Log("yanolja response", resp)
	orderResp := resp.Body.(map[string]interface{})

	update := map[string]any{
		"orderId":         req.OrderId,
		"partnerOrderId":  req.PartnerOrderId,
		"orderStatusCode": orderResp["orderStatusCode"],
		"oodoSyncStatus":  false,
		"updatedAt":       time.Now().UTC().Format(time.RFC3339),
	}

	id, err := s.mongoRepository[config.Instance().MongoDBName].UpdateOrderByOrderId(ctx, req.OrderId, update)
	if err != nil {
		level.Error(logger).Log("error", "repository error on order update ", err)
		cancellationresp, err := s.PostForcelyCancelOrder(ctx, req.PartnerOrderId)
		if err != nil {
			level.Error(logger).Log("error", "yanolja forced cancellation error")
			cancellationresp.Code = "500"
			cancellationresp.Body = fmt.Errorf("forced cancellation failed")
			return cancellationresp, err
		}
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update, %v", err), "postOrderCompletion")
	}

	level.Info(logger).Log("info", "document update with id ", id)

	var orderVariants []interface{}

	orderVariants, ok := orderResp["orderVariants"].([]interface{})
	if !ok {
		level.Error(logger).Log("error", "error retrieving order variants", "")
		cancellationresp, err := s.PostForcelyCancelOrder(ctx, req.PartnerOrderId)
		if err != nil {
			level.Error(logger).Log("error", "yanolja forced cancellation error")
			cancellationresp.Code = "500"
			cancellationresp.Body = fmt.Errorf("forced cancellation failed")
			return cancellationresp, err
		}
		return resp, customError.NewError(ctx, "leisure-api-0005", fmt.Sprint("error retrieving order variants", ""), "postOrderCompletion")
	}

	level.Info(logger).Log("info", "orderVariants ", orderVariants)

	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	timestamp := now.Format(time.RFC3339)
	for _, ov := range orderVariants {
		orderVariant, ok := ov.(map[string]interface{})
		if !ok {
			level.Error(logger).Log("error", "data type conversion error")
			cancellationResp, err := s.PostForcelyCancelOrder(ctx, req.PartnerOrderId)
			if err != nil {
				level.Error(logger).Log("error", "yanolja forced cancellation error")
				cancellationResp.Code = "500"
				cancellationResp.Body = fmt.Errorf("forced cancellation failed")
				return cancellationResp, err
			}
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1017", fmt.Sprint("orderVariant data type mismatch err  ", ""), "PostOrderCompletion")
		}
		variantid, _ := toInt64(orderVariant["variantId"])
		productid, _ := toInt64(orderVariant["productId"])
		orderVariantId, _ := toInt64(orderVariant["orderVariantId"])

		level.Info(logger).Log("Info  GetReconciliationDetailsByOrderAndVariant", fmt.Sprintln("orderId : ", req.OrderId, " orderVariantId :", orderVariantId, "variantid :", variantid, " productid :", productid))
		// Step 6.1: Try to fetch existing reconciliation details
		reconDetails, err := s.mongoRepository[config.Instance().MongoDBName].GetReconciliationDetailsByOrderAndVariant(
			ctx, req.OrderId, orderVariantId, variantid, productid)

		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) || len(reconDetails) == 0 {
				level.Warn(logger).Log("error", "reconciliation details not found, for today")
				// initialize  slice to proceed with insert later
				newDetail := domain.ReconcilationDetail{
					ReconciliationDate:       today,
					ReconcileOrderStatusCode: constant.RECONCILESTATUSCREATED,
				}
				reconDetails = append(reconDetails, newDetail)
			} else {
				// If any other error occurred, return it
				//level.Error(logger).Log(" repository error", fmt.Sprintf("failed to fetch reconciliation details : %v", err))
				level.Error(logger).Log(
					"msg", "repository error while fetching reconciliation details",
					"error", fmt.Sprintf("%v", err),
				)

				cancellationResp, err := s.PostForcelyCancelOrder(ctx, req.PartnerOrderId)
				if err != nil {
					level.Error(logger).Log("error", "yanolja forced cancellation error", err)
					cancellationResp.Code = "500"
					cancellationResp.Body = fmt.Errorf("forced cancellation failed")
					return cancellationResp, err
				}
				resp.Code = "500"
				return resp, fmt.Errorf("failed to fetch reconciliation details: %w", err)
			}
		} else {
			// Step 6.3: Update or append reconciliation detail for the given date
			reconUpdated := false
			for i, detail := range reconDetails {
				if detail.ReconciliationDate == today {
					// Update existing reconciliation detail for the day
					reconDetails[i].ReconcileOrderStatusCode = constant.RECONCILESTATUSCREATED
					reconDetails[i].ReconciliationDate = timestamp
					reconUpdated = true
					break
				}
			}
			// If not updated, append a new reconciliation detail for today
			if !reconUpdated {
				newDetail := domain.ReconcilationDetail{
					ReconciliationDate:       today,
					ReconcileOrderStatusCode: constant.RECONCILESTATUSCREATED,
				}
				reconDetails = append(reconDetails, newDetail)
			}
		}

		// Step 6.4: Update reconciliation details array in MongoDB
		filter := map[string]any{
			"orderId":        req.OrderId,
			"partnerOrderId": req.PartnerOrderId,
			"productId":      orderVariant["productId"],
			"variantId":      orderVariant["variantId"],
		}
		if err := s.mongoRepository[config.Instance().MongoDBName].UpdateReconciliationDetailByDayInsert(
			ctx, filter, reconDetails); err != nil {
			level.Error(logger).Log("error", "failed to update reconciliation details", err)
			cancellationresp, err := s.PostForcelyCancelOrder(ctx, req.PartnerOrderId)
			if err != nil {
				level.Error(logger).Log("error", "yanolja forced cancellation error")
				cancellationresp.Code = "500"
				cancellationresp.Body = fmt.Errorf("forced cancellation failed")
				return cancellationresp, err
			}
			resp.Code = "500"
			return resp, fmt.Errorf("failed to update reconciliation details in database: %w", err)
		}
	}

	var record domain.Model
	record, err = s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, req.OrderId)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		cancellationresp, err := s.PostForcelyCancelOrder(ctx, req.PartnerOrderId)
		if err != nil {
			level.Error(logger).Log("error", "yanolja forced cancellation error")
			cancellationresp.Code = "500"
			cancellationresp.Body = fmt.Errorf("forced cancellation failed")
			return cancellationresp, err
		}
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "yanolja")
	}
	level.Info(logger).Log("response", record)

	/* //only for testing voucher call
	VoucherDataCreator(ctx, logger, s, record) */

	/* fmt.Println("************* 888888888888888888888888 ************************", "call to VoucherDataCreator")
	fmt.Println("staus ::::::::::::::::::::::::: ", IsPdfVoucher(ctx, logger, s, record)) */

	resp.Body = record
	resp.Code = "200"
	return resp, nil
}

// GetOrderByOrderId get the order by id
func (s *service) GetOrderByOrderId(ctx context.Context, req yanolja.OrderConfirmation) (resp yanolja.Response, err error) {

	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetOrderByOrderId",
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
	resp, err = getsvc.GetOrderbyId(ctx, req)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		resp.Code = "500"
		resp.Body = fmt.Sprintf("yanolja client error")
		return resp, err
	}

	// Unmarshal into a map
	orderResp := resp.Body.(map[string]interface{})
	// Extract `orderId`
	orderID, ok := orderResp["orderId"]
	if !ok {
		level.Error(logger).Log("orderId is not present in response: %v\n", orderID)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-0005", fmt.Sprintf("orderId is not present in response, %v", err), "GetOrderByOrderId")
	}

	orderid, typeok := orderID.(float64)
	if !typeok {
		level.Error(logger).Log("orderId expecting int64 data: %v\n", orderID)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-0005", fmt.Sprintf("orderId type not matched, %v", err), "GetOrderByOrderId")
	}
	rec, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, int64(orderid))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no order exist based on orderid  ", err)
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderbyOrderId")
		}
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyOrderId")
	}
	resp.Code = "200"
	resp.Body = rec
	level.Info(logger).Log("response")
	return resp, nil

}

// PostCancelOrderEntirly fully cancel the order
func (s *service) PostCancelOrderEntirly(ctx context.Context, orderId int64) (resp yanolja.Response, err error) {

	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostCancelOrderEntirly",
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
	var record domain.Model
	record, err = s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, orderId)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching order based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyOrderId")
	}

	if strings.ToUpper(record.OrderStatusCode) == "DONE" {
		for _, variant := range record.OrderVariants {
			if strings.ToUpper(variant.OrderVariantStatusTypeCode) != "NOT_USED" {
				level.Info(logger).Log("info", fmt.Sprintf("OrderVariantStatusTypeCode %s of order variantId %d", variant.OrderVariantStatusTypeCode, variant.VariantID))
				return resp, customError.NewError(ctx, "leisure-api-0006", fmt.Sprintf(" orderId %d with orderVariantId %d and orderVariantStatusTypeCode %s is not fit for cancellation", orderId, variant.VariantID, variant.OrderVariantStatusTypeCode), "postCancelOrderEntirly")

			}

		}
	}
	/*
			After discussing with yanolja on 10-june-2025 they ask to change cancelling state only if successfully returned from yanolja.
		Otherwise, we will leave the order.

			err = s.mongoRepository[config.Instance().MongoDBName].UpdateOrderVariantStatusByOrderId(ctx, orderId, "CANCELING")
			if err != nil {
				level.Error(logger).Log("repository error", "error in updaing order variant status ", err)
				return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("error in updaing order variant status by  orderId %d, %v", orderId, err), "UpdateOrderVariantStatusByOrderId")
			} */

	// product.isCancelPenalty == false then cancellation is allowed

	var getsvc, _ = yanoljasvc.New(ctx)
	resp, err = getsvc.CancelFullOrder(ctx, orderId)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error during full cancellation order ", err)
		return resp, err //customError.NewError(ctx, "leisure-api-0004", fmt.Sprintf(" full cancellation order yanolja server %v", err), "yanolja")
	}

	level.Info(logger).Log("response", resp)

	err = s.mongoRepository[config.Instance().MongoDBName].UpdateOrderVariantStatusByOrderId(ctx, orderId, "CANCELING")
	if err != nil {
		level.Error(logger).Log("repository error", "error in updaing order variant status ", err)
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("error in updaing order variant status by  orderId %d, %v", orderId, err), "UpdateOrderVariantStatusByOrderId")
	}

	orderResp := resp.Body.(map[string]interface{})
	cancelStatusCode := fmt.Sprintf("%v", orderResp["cancelStatusCode"])
	cancelFailReasonCode := fmt.Sprintf("%v", orderResp["cancelFailReasonCode"])

	if cancelStatusCode != "FAIL" && cancelStatusCode != "DIRECT" && cancelStatusCode != "ADMIN" {
		level.Error(logger).Log("error", "cancelStatusCode status is incorrect", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-0009", fmt.Sprintf("cancelStatusCode is wrong"), "PostCancelOrderEntirly")
	}

	level.Info(logger).Log("info", "cancelStatusCode", cancelStatusCode)
	if cancelStatusCode == "FAIL" {
		for _, ov := range record.OrderVariants {
			err = s.mongoRepository[config.Instance().MongoDBName].UpdateCancelDetailsForVariants(ctx, orderId, ov.ProductID, ov.OrderVariantID, cancelFailReasonCode, cancelStatusCode)
			if err != nil {
				resp.Code = "500"
				level.Error(logger).Log("error", "request to yanolja client raise error ", err)
				return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update by orderId, %v", err), "PostCancelOrderEntirly")
			}

		}
	} else if cancelStatusCode == "DIRECT" || cancelStatusCode == "ADMIN" {
		level.Info(logger).Log("response", resp)

		//partnerOrderId := record.PartnerOrderID

		for _, ov := range record.OrderVariants {
			err = s.mongoRepository[config.Instance().MongoDBName].UpdateCancelDetailsForVariants(ctx, orderId,
				ov.ProductID, ov.OrderVariantID, "", cancelStatusCode)

			if err != nil {
				resp.Code = "500"
				resp.Body = record
				return resp, fmt.Errorf("error in updating failed reason and status code")
			}
		}
	}

	return resp, nil
}

// PostCancelOrderByReqTimeOut cancel order when request timeout
func (s *service) PostCancelOrderByReqTimeOut(ctx context.Context, partnerOrderId string) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostCancelOrderByReqTimeOut",
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
	resp, err = getsvc.TimeoutAndForcedOrderCancellation(ctx, partnerOrderId)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		return resp, err //customError.NewError(ctx, "leisure-api-0004", fmt.Sprintf("error from yanolja server %v", err), "yanolja")
	}

	level.Info(logger).Log("response from yanolja", resp)

	_, err = s.mongoRepository[config.Instance().MongoDBName].UpdateForcedCancelOrderDetail(ctx, partnerOrderId, constant.ORDERVARIANTCANCELEDSTATUS)
	if err != nil {
		level.Error(logger).Log("error", "timeout cancel order update issue ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update by timeout cancellation, %v", err), "UpdateForcedCancelOrderDetail")
	}

	level.Info(logger).Log("response", resp)

	return resp, nil
}

// PostForcelyCancelOrder forcely cancel order during failure
func (s *service) PostForcelyCancelOrder(ctx context.Context, partnerOrderId string) (resp yanolja.Response, err error) {
	var requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostForcelyCancelOrder",
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
	resp, err = getsvc.ForcedOrderCancellation(ctx, partnerOrderId)

	if err != nil {
		level.Error(logger).Log("error", "request to yanolja client raise error ", err)
		return resp, err //customError.NewError(ctx, "leisure-api-0004", fmt.Sprintf("error from yanolja server %v", err), "yanolja")
	}

	level.Info(logger).Log("response from yanolja", resp)

	_, err = s.mongoRepository[config.Instance().MongoDBName].UpdateForcedCancelOrderDetail(ctx, partnerOrderId, constant.ORDERVARIANTCANCELEDSTATUS)
	if err != nil {
		level.Error(logger).Log("error", "forced cancel order update issue ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on order update by forced cancellation, %v", err), "UpdateForcedCancelOrderDetail")
	}

	level.Info(logger).Log("response", resp)

	return resp, nil
}

// GetOrderReconcilationDetail get the order by id
func (s *service) GetOrderReconcilationDetail(ctx context.Context, req yanolja.OrderReconcileReq) (resp yanolja.Response, err error) {

	record := make([]yanolja.OrderReconcilation, 0)
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetOrderReconcilationDetail",
		"Request ID", requestID,
		"Trace ID", resp.TraceID,
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

	// as per yanolja document Page number (min:1), Number of items per page. 1000 requests (min: 1, max:1000)
	if req.PageNumber < 1 || req.PageSize < 1 || req.PageSize > 1000 {
		return resp, customError.NewError(ctx, "leisure-api-0001", " Page number (min:1), Number of items per page. 1000 requests (min:1, max:1000)", nil)
	}

	statusCode := strings.ToUpper(req.ReconcileOrderStatusCode)
	if statusCode != constant.RECONCILESTATUSCANCELED && statusCode != constant.RECONCILESTATUSCREATED &&
		statusCode != constant.RECONCILESTATUSRESTORED && statusCode != constant.RECONCILESTATUSUSED {

		level.Error(logger).Log("error", "reconcilation status is incorrectly entered", err)
		resp.Code = "400"
		return resp, customError.NewError(ctx, "leisure-api-0001", fmt.Sprintf("reconcilation status wrongly entered"), "GetOrderReconcilationDetail")
	}

	record, err = s.mongoRepository[config.Instance().MongoDBName].GetReconciliationDetailsByDateAndStatus(ctx, req.ReconciliationDate, statusCode)
	if err != nil {
		level.Error(logger).Log("error", "request to GetOrderReconcilationDetail raised error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error from GetReconciliationDetailsByDateAndStatus, %v", err), "GetOrderReconcilationDetail")
	}

	lengthOfRecord := len(record)

	// Initialize response struct with default values.
	resp = yanolja.Response{
		Code:        "",
		Collection:  true,
		ContentType: nil,
		Page: yanolja.NumberOfPage{
			Number:            req.PageNumber,
			Size:              req.PageSize,
			TotalElementCount: lengthOfRecord,
			TotalPageCount:    (lengthOfRecord + req.PageSize - 1) / req.PageSize,
		},
		Body: map[string]interface{}{
			"orders": []yanolja.OrderReconcilation{},
		},
	}

	// Validate and adjust pagination bounds.
	if lengthOfRecord > 0 {
		start := (req.PageNumber - 1) * req.PageSize
		end := req.PageNumber * req.PageSize

		// Ensure start is within range.
		if start < lengthOfRecord {
			// Adjust end to avoid slicing errors.
			if end > lengthOfRecord {
				end = lengthOfRecord
			}

			// Slice the records for the current page.
			paginatedOrders := record[start:end]

			// Populate the response body with the paginated data.
			resp.Body = map[string]interface{}{
				"orders": paginatedOrders,
			}
		}
	}

	// Convert response to JSON if needed (for HTTP response).
	/* responseJSON, err := json.Marshal(resp)
	if err != nil {
		logger.Log("msg", "Error marshalling response", "err", err)
		return resp, err
	} */

	level.Info(logger).Log("response of reconcilation", resp)
	resp.Code = "200"

	return resp, nil

}

func PreOrderCreationValidate(ctx context.Context, s *service, req yanolja.WaitingForOrder) (err error) {
	s.logger.Log("info", "PreOrderCreationValidate")

	var totalPurchaseQuantity int32
	//value := constant.PRODUCTVARIANTSTATUSCODE

	level.Info(s.logger).Log(" selectVariant ", req.SelectVariants)
	for _, val := range req.SelectVariants {
		totalPurchaseQuantity += val.Quantity
		product, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, val.ProductID)
		if err != nil {
			level.Error(s.logger).Log("error", "product doesnot exist ", err)
			return customError.NewError(ctx, "leisure-api-0006", fmt.Sprintf("repository error from PreOrderCreationValidate, %v", err), "PreOrderCreationValidate")
		}
		if product.ProductVersion != val.ProductVersion {
			e := customError.NewErrorCustom(ctx, "leisure-api-1016", fmt.Sprintf("PreOrderCreationValidate failed with requested and available product version  %d --- %d mismatch", val.ProductVersion, product.ProductVersion), "validation error", 400, "PreOrderCreationValidate")
			return e
		}
		for _, productoptiongroup := range product.ProductOptionGroups {
			for _, variant := range productoptiongroup.Variants {
				if variant.VariantID == val.VariantID && variant.ProductID == val.ProductID {
					if !validator.ValidateProductVariantStatusCode(ctx, product.ProductStatusCode) || !validator.ValidateProductVariantStatusCode(ctx, variant.VariantStatusCode) {
						e := customError.NewErrorCustom(ctx, "leisure-api-1016", fmt.Sprintln("enum validation error for product/variant statuscode"), "validation error", 500, "PreOrderCreationValidate")
						return e
					}

					if product.ProductStatusCode != "IN_SALE" || variant.VariantStatusCode != "IN_SALE" {
						e := customError.NewError(ctx, "leisure-api-1016", fmt.Sprintln("End of sale for product/variant  ", "PreOrderCreationValidate"), nil)
						return e
					}
					if variant.QuantityPerPersonValidityDays == 0 {
						if variant.QuantityPerPerson != -1 && variant.QuantityPerPerson < val.Quantity {
							e := customError.NewError(ctx, "leisure-api-0016", fmt.Sprintf("error in validating %s", "quantity ordered exceeds per person purchase for  per person validity days value of 0"), PreOrderCreationValidate)
							return e
						}
					} else if variant.QuantityPerPersonValidityDays > 0 {

						today := time.Now().UTC()

						startDateStr := variant.SalePeriod.StartDateTime

						validEndDatestr := variant.SalePeriod.EndDateTime
						ValidEndDate, err := time.Parse(time.RFC3339, validEndDatestr)
						if err != nil {
							e := customError.NewError(ctx, "leisure-api-1016", fmt.Sprintf("error due to startDate datetime_conversion %s", "salesPeriod_startDate"), nil)
							return e
						}
						// Parse the string to time.Time
						startDate, err := time.Parse(time.RFC3339, startDateStr)

						if err != nil {
							e := customError.NewError(ctx, "leisure-api-1016", fmt.Sprintf("error due to startDate datetime_conversion %s", "salesPeriod_startDate"), nil)
							return e
						}

						if today.Before(startDate) {
							e := customError.NewError(ctx, "leisure-api-0016", fmt.Sprintf(" todays date %s is not a valid date  date for purchase ", today), PreOrderCreationValidate)
							return e
						}
						// Add 10 days to the parsed date
						endDate := today.AddDate(0, 0, int(variant.QuantityPerPersonValidityDays))

						// Convert the new date back to string format (RFC3339)
						//endDate, err := time.Parse (time.RFC3339, newDate)

						// Compare if currentDate is outside startDate and endDate
						if endDate.After(ValidEndDate) {
							e := customError.NewError(ctx, "leisure-api-1016", fmt.Sprintf("person %s exceeds its booking days %d limit for productid %d and variantid %d during saleperiod period of variant", req.Customer.Name, variant.QuantityPerPersonValidityDays, variant.ProductID, variant.VariantID), PreOrderCreationValidate)
							return e
						}

						//  add count check for days a person allowd to book for particular variant during saleperiod

						totaQuantityPurchasedByPersonToday, err := s.mongoRepository[config.Instance().MongoDBName].
							GetTotalQuantityPurchasedByPersonToday(ctx, req.Customer.Name, req.Customer.Email, req.Customer.Tel, val.ProductID, val.VariantID)
						if err != nil {
							e := customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("error in fetching total_quantity_purchase by person today %s", "GetTotalQuantityPurchasedByPerson"), PreOrderCreationValidate)
							return e
						}

						if variant.QuantityPerPerson != -1 && variant.QuantityPerPerson < (totaQuantityPurchasedByPersonToday+val.Quantity) {
							e := customError.NewError(ctx, "leisure-api-0016", fmt.Sprintf("error in validating %s", "total_purchase_by_person exceeds the per day limit"), PreOrderCreationValidate)
							return e
						}
					}

					if variant.QuantityPerPurchase < val.Quantity && variant.QuantityPerPurchase != -1 {
						e := customError.NewError(ctx, "leisure-api-0016", fmt.Sprintf("exceed the limit of  %s", "quantity_per_purchase"), PreOrderCreationValidate)
						return e
					}

				}
			}
		}
	}

	if totalPurchaseQuantity > constant.MAXORDERQUANTITYONETIME {
		e := customError.NewError(ctx, "leisure-api-0016", fmt.Sprintf("selectvariants order quantity %d exceeds maximum one time total_order_quantity value of 30 %s", totalPurchaseQuantity, ""), PreOrderCreationValidate)
		return e
	}

	return nil
}

func toInt64(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case string:
		// Parse string to int64
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unsupported type: %v", reflect.TypeOf(value))
	}
}

// UpdateReconcilationDetail  func to update individual variant per day basis
func UpdateReconcilationDetail(ctx context.Context, s *service, logger log.Logger, orderId int64, partnerOrderId string, ov domain.OrderVariant) (err error) {
	level.Info(logger).Log("function name", "UpdateReconciliationDetail")

	// Update MongoDB with the modified reconciliation details array
	filter := map[string]any{
		"orderId":        orderId,
		"partnerOrderId": partnerOrderId,
		"productId":      ov.ProductID,
		"variantId":      ov.VariantID,
	}
	// Fetch existing reconciliation details
	reconDetails, err := s.mongoRepository[config.Instance().MongoDBName].GetReconciliationDetailsByOrderAndVariant(
		ctx, orderId, ov.OrderVariantID, ov.VariantID, ov.ProductID)

	// Define today's date
	todayDate := time.Now().UTC().Format("2006-01-02")
	if err == mongo.ErrNoDocuments && len(reconDetails) == 0 {
		// Create new entry if no documents are found
		level.Error(logger).Log("repository error", "no record exists on particular day")
		newDetail := domain.ReconcilationDetail{
			ReconciliationDate:       todayDate,
			ReconcileOrderStatusCode: constant.RECONCILESTATUSCANCELED,
		}
		reconDetails = []domain.ReconcilationDetail{newDetail}

		if err := s.mongoRepository[config.Instance().MongoDBName].UpdateReconciliationDetailByDay(ctx, filter, reconDetails); err != nil {
			level.Error(logger).Log("error", "failed to update reconciliation details", err)
			return fmt.Errorf("failed to update reconciliation details in database: %w", err)
		}
	} else if err != nil && len(reconDetails) == 0 {
		// Handle any other error
		level.Error(logger).Log("error", "repository error from GetReconciliationDetailsByOrderAndVariant")
		return fmt.Errorf("failed to fetch reconciliation details: %w", err)
	} else if err == nil && len(reconDetails) != 0 {

		// Loop through each reconciliation detail to check for today's date
		for i, detail := range reconDetails {
			if detail.ReconciliationDate == todayDate && detail.ReconcileOrderStatusCode != constant.RECONCILESTATUSCANCELED {
				reconDetails[i].ReconcileOrderStatusCode = constant.RECONCILESTATUSCANCELED // Update the existing entry for today

				if err := s.mongoRepository[config.Instance().MongoDBName].UpdateReconciliationDetailByDay(ctx, filter, reconDetails); err != nil {
					level.Error(logger).Log("error", "failed to update reconciliation details", err)
					return fmt.Errorf("failed to update reconciliation details in database: %w", err)
				}
			}
		}
	}
	return nil
}

func (s *service) GetEverlandOrders(ctx context.Context, req yanolja.EverlandGetRequest) (resp common.Response, err error) {

	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetEverlandOrders",
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

	rec, err := s.mongoRepository[config.Instance().MongoDBName].GetOrdersByChannelCodeAndCustomerEmail(ctx, req.ChannelCode, req.CustomerEmail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no record exist based on cuatomer email and channelcode  ", err)
			resp.Code = "404"
			resp.Status = http.StatusBadRequest
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetEverlandOrders")
		}

		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyOrderId")
	}

	resp.Code = "200"
	resp.Body = rec

	return resp, nil

}

// GetOderByPartialPartnerOrderIdSuffix
func (s *service) GetOderByPartialPartnerOrderIdSuffix(ctx context.Context, partialPartnerId string) (resp common.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetOderByPartialPartnerOrderId",
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

	rec, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderByPartnerIdSuffix(ctx, partialPartnerId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no record exist based partnerOrderId suffix  ", err)
			resp.Code = "404"
			resp.Status = http.StatusBadRequest
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderByPartnerIdSuffix")
		}

		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by partnerOrderId suffix, %v", err), "GetOrderbyOrderId")
	}

	resp.Code = "200"
	resp.Body = rec

	return resp, nil

}
