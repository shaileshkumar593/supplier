package implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/request_response/travolution"
	travolutionSvc "swallow-supplier/services/suppliers/travolution"

	domain "swallow-supplier/mongo/domain/travolution"
	"swallow-supplier/utils/constant"
	"time"

	"swallow-supplier/utils"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateTravolutionOrder  create order for booking ticket
func (s *service) CreateTravolutionOrder(ctx context.Context, req travolution.OrderRequest) (resp travolution.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "CreateTravolutionOrder",
		"Request ID", requestID,
	)
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<=======req===================>>>>>>>>>>>>>>>>.", req)
	fmt.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<=======req.Option===================>>>>>>>>>>>>>>>>.", reflect.ValueOf(req.Option).Kind())
	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	level.Info(logger).Log(" info ", "travolution service call")
	fmt.Println("---------1----------")
	product, err := s.mongoRepository[config.Instance().MongoDBName].GetProductByProductUid(ctx, req.Product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no product exist based on productUid  ", err)
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetProductByProductUid")
		}
		level.Error(logger).Log("repository error", "fetching record based on productUid ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching product by productUid, %v", err), "GetProductByProductUid")

	}
	fmt.Println("---------2----------")

	// AV: Available
	// AP: Approved(Used)
	// CR: Cancel Request
	// CL: Canceled
	// EP: Expired

	// add request to mongo  pass product.Type
	id, err := s.mongoRepository[config.Instance().MongoDBName].UpsertTravolutionOrder(ctx, req, product.Type, "AV")
	if id == "" && err != nil {
		level.Error(logger).Log(" error ", err)
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}
	fmt.Println("---------3----------")

	var tsvc, _ = travolutionSvc.New(ctx)
	resp, err = tsvc.TravolutionOrder(ctx, req, product.Type)
	if err != nil {
		level.Error(logger).Log(" travolution error ", err)
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}
	fmt.Println("---------4----------")

	now := time.Now().UTC().Format(time.RFC3339)

	ticketResp := resp.Body.(map[string]interface{})

	updates := make(map[string]interface{}, 0)

	if ticketResp["type"] == "BK" {
		updates = map[string]interface{}{
			"_id":             id,
			"referenceNumber": ticketResp["referenceNumber"],
			"voucherType":     ticketResp["voucherType"],
			"voucherInfo":     ticketResp["voucherInfo"],
			"orderNumber":     ticketResp["orderNumber"],
			"status":          constant.BOOKINGPENDING,
			"bookingStatus":   constant.BOOKINGPENDING,
			"approvedAt":      now,
			"expiredAt":       ticketResp["expiredAt"],
		}
	} else {
		updates = map[string]interface{}{
			"_id":             id,
			"referenceNumber": ticketResp["referenceNumber"],
			"voucherType":     ticketResp["voucherType"],
			"voucherInfo":     ticketResp["voucherInfo"],
			"orderNumber":     ticketResp["orderNumber"],
			"status":          constant.ORDERAVAILABLE,
			"approvedAt":      now,
			"expiredAt":       ticketResp["expiredAt"],
		}
	}

	fmt.Println("---------5----------")

	id, err = s.mongoRepository[config.Instance().MongoDBName].UpdateTravolutionOrderById(ctx, updates)
	if err != nil {
		return travolution.Response{}, err
	}
	fmt.Println("---------6----------")
	fmt.Println("******************* orderNumber ********************** ", ticketResp["orderNumber"].(string))

	order, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderByOrderNumber(ctx, ticketResp["orderNumber"].(string))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", err.Error())
			err = fmt.Errorf("no document exist with orderNumber: %s", ticketResp["orderNumber"].(string))
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderByOrderNumber")
		}
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderNumber, %v", err), "", http.StatusInternalServerError, "GetOrderByOrderNumber")
	}
	fmt.Println("---------7----------")

	resp.Body = order
	resp.Code = "200"

	return resp, nil
}

// SearchTravolutionOrder get order from travolution
func (s *service) SearchTravolutionOrder(ctx context.Context, orderNumber string) (resp travolution.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "SearchTravolutionOrder",
		"Request ID", requestID,
	)

	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	level.Info(logger).Log(" info ", "travolution service call")

	var tsvc, _ = travolutionSvc.New(ctx)
	resp, err = tsvc.TravolutionGetOrder(ctx, orderNumber)
	if err != nil {
		level.Error(logger).Log(" travolution error ", err)
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}

	return resp, nil
}

// CancelTravolutionOrder cancel order with orderNumber
func (s *service) CancelTravolutionOrder(ctx context.Context, orderNumber string) (resp travolution.Response, err error) {
	var requestID string

	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "CancelTravolutionOrder",
		"Request ID", requestID,
	)

	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	order, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderByOrderNumber(ctx, orderNumber)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", err.Error())
			err = fmt.Errorf("no document exist with orderNumber: %s", orderNumber)
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderByOrderNumber")
		}
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderNumber, %v", err), "", http.StatusInternalServerError, "GetOrderByOrderNumber")
	}

	if order.Status != constant.ORDERAVAILABLE {
		resp.Code = "400"
		resp.Body = fmt.Errorf("orderNumber %s is not fit for cancellation ", orderNumber)
		return resp, fmt.Errorf("orderNumber %s is not fit for cancellation ", orderNumber)
	}

	updates := map[string]interface{}{
		"status":            constant.ORDERCANCELREQUEST,
		"cancelRequestedAt": time.Now().UTC().Format(time.RFC3339),
	}

	err = s.mongoRepository[config.Instance().MongoDBName].UpdateOrderByOrderNumber(ctx, orderNumber, updates)
	if err != nil {
		level.Error(logger).Log("repository error ", err)
		resp.Code = "500"
		resp.Body = fmt.Errorf("repository error on updating Cancel Request %w", err)
		return resp, err
	}

	level.Info(logger).Log(" info ", "travolution service call")

	// update cancel_request
	var tsvc, _ = travolutionSvc.New(ctx)
	resp, err = tsvc.TravolutionCancelOrder(ctx, orderNumber)
	if err != nil {
		level.Error(logger).Log(" travolution error ", err)
		resp.Code = "500"
		resp.Body = err.Error()
		return resp, err
	}

	order, err = s.mongoRepository[config.Instance().MongoDBName].GetOrderByOrderNumber(ctx, orderNumber)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", err.Error())
			err = fmt.Errorf("no document exist with orderNumber: %s", orderNumber)
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderByOrderNumber")
		}
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderNumber, %v", err), "", http.StatusInternalServerError, "GetOrderByOrderNumber")
	}

	resp.Body = order
	resp.Code = "200"

	return resp, nil
}

// Webhook for travolution order update for REDEEMED
func (s *service) OrderWebhookUpdate(ctx context.Context, payload domain.Webhook) (resp travolution.Response, err error) {
	var requestID string
	requestID = utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "OrderWebhookUpdate",
		"Request ID", requestID,
	)

	// Defer panic recovery
	defer func(ctx context.Context) {
		if r := recover(); r != nil {
			level.Info(logger).Log("info", "processing request went into panic mode", "panic", r)
			resp.Code = "500"
		}
	}(ctx)

	fmt.Println("=====================1========================")
	// insert event call to mongo
	_, err = s.mongoRepository[config.Instance().MongoDBName].UpsertTravolutionWebhook(ctx, payload)
	if err != nil {
		level.Error(logger).Log(" InsertTravolutionWebhook reposotory error ", fmt.Sprintf("%w", err))
		resp.Code = "500"
		resp.Body = fmt.Errorf("InsertTravolutionWebhook error in repository %w ", err)
		return resp, err
	}

	fmt.Println("=====================2========================")

	level.Info(logger).Log(" info ", "travolution service call")

	switch payload.EventType {

	case "BOOKING_ACCEPTED":
		// "BOOKING_ACCEPTED"
		//status = "AV" // available
		fmt.Println("=====================3========================")

		update := map[string]any{
			"bookingStatus": constant.BOOKINGAPPROVED,
			"status":        constant.ORDERAVAILABLE,
			"eventType":     payload.EventType,
			"dateAt":        payload.Data.DateAt,
			"updatedAt":     time.Now().UTC().Format(time.RFC3339),
		}
		_, err = s.mongoRepository[config.Instance().MongoDBName].UpsertWebhookToOrder(ctx, payload, update)
		if err != nil {
			level.Error(logger).Log(" UpsertWebhookToOrder reposotory error  ", fmt.Sprintf("%w", err))
			resp.Code = "500"
			resp.Body = fmt.Errorf("UpsertWebhookToOrder error in repository %w ", err)
			return resp, err
		}
	case "REDEEMED":
		// "REDEEMED"
		fmt.Println("=====================4========================")

		update := map[string]any{
			"status":     constant.ORDERAPPROVED,
			"approvedAt": payload.Data.DateAt,
			"eventType":  payload.EventType,
			"dateAt":     payload.Data.DateAt,
			"updatedAt":  time.Now().UTC().Format(time.RFC3339),
		}
		//status = "AP" // used
		_, err = s.mongoRepository[config.Instance().MongoDBName].UpsertWebhookToOrder(ctx, payload, update)
		if err != nil {
			level.Error(logger).Log(" UpsertWebhookToOrder reposotory error  ", fmt.Sprintf("%w", err))
			resp.Code = "500"
			resp.Body = fmt.Errorf("UpsertWebhookToOrder error in repository %w ", err)
			return resp, err
		}
	case "RESTORED":
		// "RESTORED"
		fmt.Println("=====================5========================")

		update := map[string]any{
			"status":    constant.ORDERAVAILABLE,
			"eventType": payload.EventType,
			"dateAt":    payload.Data.DateAt,
			"updatedAt": time.Now().UTC().Format(time.RFC3339),
		}
		//status = "AV" // available
		_, err = s.mongoRepository[config.Instance().MongoDBName].UpsertWebhookToOrder(ctx, payload, update)
		if err != nil {
			level.Error(logger).Log(" UpsertWebhookToOrder reposotory error  ", fmt.Sprintf("%w", err))
			resp.Code = "500"
			resp.Body = fmt.Errorf("UpsertWebhookToOrder error in repository %w ", err)
			return resp, err
		}

	case "CANCELED":
		// "CANCELED"
		fmt.Println("=====================6========================")

		update := map[string]any{
			"status":    constant.ORDERCANCELED,
			"eventType": payload.EventType,
			"dateAt":    payload.Data.DateAt,
			"updatedAt": time.Now().UTC().Format(time.RFC3339),
		}
		//status = "CL"
		_, err = s.mongoRepository[config.Instance().MongoDBName].UpsertWebhookToOrder(ctx, payload, update)
		if err != nil {
			level.Error(logger).Log(" UpsertWebhookToOrder reposotory error  ", fmt.Sprintf("%w", err))
			resp.Code = "500"
			resp.Body = fmt.Errorf("UpsertWebhookToOrder error in repository %w ", err)
			return resp, err
		}

	case "BOOKING_REJECTED":
		// "BOOKING_REJECTED"
		fmt.Println("=====================7========================")

		update := map[string]any{
			"bookingStatus": constant.BOOKINGREJECTED,
			"eventType":     payload.EventType,
			"dateAt":        payload.Data.DateAt,
			"updatedAt":     time.Now().UTC().Format(time.RFC3339),
		}
		//status = "RJ"
		_, err = s.mongoRepository[config.Instance().MongoDBName].UpsertWebhookToOrder(ctx, payload, update)
		if err != nil {
			level.Error(logger).Log(" UpsertWebhookToOrder reposotory error  ", fmt.Sprintf("%w", err))
			resp.Code = "500"
			resp.Body = fmt.Errorf("UpsertWebhookToOrder error in repository %w ", err)
			return resp, err
		}
		fmt.Println("=====================8========================")

		// BOOKING_REJECTED event, then immediately proceed with cancellation of order
		resp, err = s.CancelTravolutionOrder(ctx, payload.Data.OrderNumber)
		if err != nil {
			level.Error(logger).Log(" CancelTravolutionOrder service error  ", fmt.Sprintf("%w", err))
			resp.Code = "500"
			resp.Body = fmt.Errorf("CancelTravolutionOrder service error %w ", err)
			return resp, err
		}
		fmt.Println("=====================9========================")

	default:
		// Unknown / fallback
		resp.Body = fmt.Sprintf("Unsupported event type: %s", payload.EventType)
		resp.Code = "400"
		return resp, nil
	}
	fmt.Println("=====================10========================")

	order, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderByOrderNumber(ctx, payload.Data.OrderNumber)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", err.Error())
			err = fmt.Errorf("no document exist with orderNumber: %s", payload.Data.OrderNumber)
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderByOrderNumber")
		}
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderNumber, %v", err), "", http.StatusInternalServerError, "GetOrderByOrderNumber")
	}

	resp.Code = "200"
	resp.Body = order
	return resp, nil
}

// Detect and unmarshal into correct response struct
func ParseResponseBody(body interface{}, logger log.Logger, typeOfBooking string) (interface{}, error) {
	// Marshal interface{} back to JSON bytes
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	switch typeOfBooking {
	case "ticket":
		var resp travolution.TicketResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			level.Error(logger).Log("unmarshal error ", err)
			return nil, fmt.Errorf("failed to unmarshal TicketResponse: %w", err)
		}
		return resp, nil

	case "booking":
		var resp travolution.BookingResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			level.Error(logger).Log(" unmarshal error ", err)
			return nil, fmt.Errorf("failed to unmarshal BookingResponse: %w", err)
		}
		return resp, nil

	case "passpkg":
		var resp travolution.PassPkgResponse
		if err := json.Unmarshal(data, &resp); err != nil {
			level.Error(logger).Log(" unmarshal error ", err)
			return nil, fmt.Errorf("failed to unmarshal PassPkgResponse: %w", err)
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown typeOfBooking: %s", typeOfBooking)
	}
}
