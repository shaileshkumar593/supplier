package implementation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"swallow-supplier/caches/cache"
	"swallow-supplier/config"
	customError "swallow-supplier/error"
	mapper "swallow-supplier/mapper/trip_to_yanolja"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/trip"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/services/pdfvoucher"
	"swallow-supplier/utils/constant"

	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostRequestFromGGT request all request from trip
func (s *service) PostRequestFromGGT(ctx context.Context, req trip.SwallowRequest) (resp yanolja.Response, err error) {

	logger := log.With(
		s.logger,
		"method", "PostRequestFromGGT",
		"request ", req,
	)
	var flag = true

	startTime := time.Now()
	serviceName := req.Header.ServiceName
	if serviceName == "" {
		return resp, customError.NewError(ctx, "leisure-api-1018", fmt.Sprintf("service name is empty, %v", err), "PostRequestFromGGT")
	}
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		resp.Code = "500"
		return resp, err
	}

	mrepo := s.mongoRepository[config.Instance().MongoDBName]
	level.Info(logger).Log("serviceName", serviceName)

	switch serviceName {
	case "CreatePreOrder":
		var preorderRequest trip.PreorderRequest
		err := json.Unmarshal([]byte(req.DecryptedBody), &preorderRequest)
		if err != nil {
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1019", fmt.Sprintf("unmarshaling error %v", err), "PostRequestFromGGT")
		}

		level.Info(logger).Log(" preorderRequest ", preorderRequest)
		//var reqYanolja = yanolja.WaitingForOrder{}
		reqYanolja, plus, err := mapper.PreOrderMapper(ctx, mrepo, logger, preorderRequest)
		if err != nil {
			resp.Code = "500"
			return resp, err
		}

		level.Info(logger).Log(" mapper ", reqYanolja)
		// For checking duplicate request
		sequenceIdExist, _ := s.mongoRepository[config.Instance().MongoDBName].GetSequenceIDByOtaOrderIDAndRequestCategory(ctx, preorderRequest.OtaOrderID, serviceName)
		if sequenceIdExist {
			record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyPartnerOrderId(ctx, preorderRequest.OtaOrderID)
			if err != nil {
				level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
				resp.Code = "500"
				return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyPartnerOrderId")
			}
			resp.Code = "200"
			resp.Body = trip.CreatePreOrder{
				PLU:   plus,
				Order: record,
			}
			logger.Log("*****************trip resp**************", resp)
			return resp, nil
		}

		level.Info(logger).Log("yanolja request : ", reqYanolja)
		respp, err := s.PostWaitForOrder(ctx, reqYanolja)
		if err != nil {
			resp.Body = err.Error()
			resp.Code = "500"
			return resp, err
		}

		// inserting request from trip to sync to redis
		err = s.mongoRepository[config.Instance().MongoDBName].InsertPreorderRequestFromTrip(ctx, preorderRequest)
		if err != nil {
			level.Error(logger).Log("repository error", "inserting trip priorder request ")
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on inserting trip priorder request, %v", err), "InsertPreorderRequestFromTrip")
		}
		order, _ := respp.Body.(domain.Model)
		createPreorderResp := trip.CreatePreOrder{
			PLU:   plus,
			Order: order,
		}
		resp.Body = createPreorderResp

	case "PayPreOrder":
		fmt.Println("PayPreOrder 33")
		var confirmTicket trip.PreOrderPaymentRequest
		err := json.Unmarshal([]byte(req.DecryptedBody), &confirmTicket)
		if err != nil {
			level.Error(logger).Log("error :", err)
			resp.Code = "500"
			return resp, err
		}
		fmt.Println("PayPreOrder 34")

		// For checking duplicate request
		sequenceIdExist, _ := s.mongoRepository[config.Instance().MongoDBName].GetSequenceIDByOtaOrderIDAndRequestCategory(ctx, confirmTicket.OtaOrderId, serviceName)
		if sequenceIdExist {
			OrderId, err := strconv.ParseInt(confirmTicket.SupplierOrderId, 10, 64)
			if err != nil {
				level.Error(logger).Log("parseInt error :", err)
				resp.Code = "500"
				return resp, err
			}
			record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, OrderId)
			if err != nil {
				level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
				resp.Code = "500"
				return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyOrderId")
			}

			tripresp, err := MakeTripResponse(ctx, logger, record)
			if err != nil {
				level.Error(logger).Log("error ", "error from makeTripResponse")
				resp.Code = "500"
				return resp, err
			}
			tripresp.Message = fmt.Sprintf("duplicate order confirmation, no order waiting list matching the requested order number %d", record.OrderId)
			resp.Code = "200"
			resp.Body = tripresp
			logger.Log(" trip duplicate order confirmation response ", resp)
			return resp, nil
		}

		var itemAry []domain.Item
		for _, item := range confirmTicket.Items {
			intemInfo := domain.Item{
				ItemId: item.ItemId,
				PLU:    item.PLU,
			}
			itemAry = append(itemAry, intemInfo)
		}
		fmt.Println("PayPreOrder 36")

		var confirmation yanolja.OrderConfirmation
		confirmation.OrderId, err = strconv.ParseInt(confirmTicket.SupplierOrderId, 10, 64)
		if err != nil {
			level.Error(logger).Log("parseInt error :", err)
			resp.Code = "500"
			return resp, err
		}
		fmt.Println("PayPreOrder 37")

		confirmation.PartnerOrderId = confirmTicket.OtaOrderId
		itemIdDetail := domain.ItemIdDetails{
			OrderId: confirmation.OrderId,
			Items:   itemAry,
		}
		fmt.Println("PayPreOrder 38")

		err = s.mongoRepository[config.Instance().MongoDBName].UpsertItemIdDetails(ctx, itemIdDetail)
		if err != nil {
			level.Error(logger).Log("repository error", "updating itemDetail in collection ", err)
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on updating itemdetail by orderId, %v", err), "UpsertItemIdDetails")
		}

		fmt.Println("PayPreOrder 39")

		level.Info(logger).Log("method call", "PostOrderCompletion")
		payDetail, err := s.PostOrderCompletion(ctx, confirmation)
		if err != nil {
			resp.Code = "500"
			return resp, err
		}
		fmt.Println("PayPreOrder 40")

		respData, err := MakeTripResponse(ctx, logger, payDetail.Body)
		if err != nil {
			resp.Code = "500"
			return resp, err
		}
		err = s.mongoRepository[config.Instance().MongoDBName].InsertPaymentRequestFromTrip(ctx, confirmTicket)
		if err != nil {
			level.Error(logger).Log("repository error", "inserting trip payment confirm request ")
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on inserting trip payment confirm request, %v", err), "InsertPaymentRequestFromTrip")
		}

		// for pdf voucher generation
		pdfVoucherFlag := IsPdfVoucher(ctx, logger, s, payDetail.Body)
		logger.Log("pdfVoucherCallFlag   :      ", pdfVoucherFlag)

		if pdfVoucherFlag {
			s.logger.Log("VoucherDataCreator called")
			VoucherDataCreator(ctx, logger, s, payDetail.Body)
		}

		resp.Code = "200"
		resp.Body = respData
		fmt.Println("PayPreOrder 41")

	case "CancelPreOrder":
		// timeout
		fmt.Println("CancelPreOrder 60")
		var timeoutCancel trip.PreOrderTimeoutCancellation
		err := json.Unmarshal([]byte(req.DecryptedBody), &timeoutCancel)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			resp.Code = "500"
			return resp, err
		}

		level.Info(logger).Log("method call", "PostCancelOrderByReqTimeOut")

		resp, _ = s.PostCancelOrderByReqTimeOut(ctx, timeoutCancel.OtaOrderId)

		record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyPartnerOrderId(ctx, timeoutCancel.OtaOrderId)
		if err != nil {
			level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyPartnerOrderId")
		}

		resp.Code = "200"
		resp.Body = record

	case "ForcedCancelOrder":
		// forced cancel
		// need to confirm and delete

		var forcedCancel trip.CancellationRequest
		err := json.Unmarshal([]byte(req.DecryptedBody), &forcedCancel)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			resp.Code = "500"
			return resp, err
		}

		level.Info(logger).Log("method call", "PostForcelyCancelOrder")

		resp, _ = s.PostForcelyCancelOrder(ctx, forcedCancel.OTAOrderID)

		record, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyPartnerOrderId(ctx, forcedCancel.OTAOrderID)
		if err != nil {
			level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "GetOrderbyPartnerOrderId")
		}
		resp.Code = "200"
		resp.Body = record

	case "CancelOrder":
		// full order cancel is async cancellation
		fmt.Println("CancelOrder 70")
		var fullCancelOrder trip.CancellationRequest
		err := json.Unmarshal([]byte(req.DecryptedBody), &fullCancelOrder)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			resp.Code = "500" // 500
			return resp, err
		}

		fmt.Println("CancelOrder 71")

		orderId, err := strconv.ParseInt(fullCancelOrder.SupplierOrderID, 10, 64)
		if err != nil {
			level.Error(logger).Log("error to trip  :", err)
			resp.Code = "500" // 500
			return resp, err
		}
		fmt.Println("CancelOrder 73")

		order, err := s.mongoRepository[config.Instance().MongoDBName].FindOrderByOrderIdAndPartnerOrderId(ctx, orderId, fullCancelOrder.OTAOrderID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				level.Error(logger).Log("repository error", "no record exist based on orderid and partnerorderid  ", err)
				resp.Code = "404"
				return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "FindOrderByOrderIdAndPartnerOrderId")
			}
			level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
			resp.Code = "500" // 500
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "FindOrderByOrderIdAndPartnerOrderId")
		}
		fmt.Println("CancelOrder 74")

		for i, item := range fullCancelOrder.Items {
			variantId, productversion, productid, err := GetVariantIdFromPlu(ctx, logger, item.PLU)
			if err != nil {
				resp.Code = "500" // 500
				resp.Body = err.Error()
				return resp, err
			}
			fmt.Println("CancelOrder 75")

			if order.SelectVariants[i].VariantID == variantId &&
				order.SelectVariants[i].ProductVersion == productversion &&
				order.SelectVariants[i].ProductID == productid &&
				order.SelectVariants[i].Quantity != int32(item.Quantity) {
				resp.Code = "400" //400
				resp.Body = fmt.Errorf("order quantity and cancel quantity donot match ")
				return resp, customError.NewError(ctx, "leisure-api-00023", fmt.Sprintf("requested cancel quantity %d not match with order quantity %d", int32(item.Quantity), order.SelectVariants[i].Quantity), "PostRequestFromGGT")
			}

		}
		level.Info(logger).Log("method call", "PostCancelOrderEntirly")
		detail, err := s.PostCancelOrderEntirly(ctx, orderId)
		fmt.Println("CancelOrder 76")

		if err != nil {
			level.Error(logger).Log("error to trip :", err)
			resp.Code = "500" //500
			resp.Body = err.Error()
			return resp, err
		}
		fmt.Println("CancelOrder 77")

		err = s.mongoRepository[config.Instance().MongoDBName].InsertFullCancelOrderRequestFromTrip(ctx, fullCancelOrder)
		if err != nil {
			level.Error(logger).Log("repository error", "inserting trip full cancel order request ")
			resp.Code = "500"
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on inserting trip full cancel order request, %v", err), "InsertPaymentRequestFromTrip")
		}

		resp.Body = detail.Body
		resp.Code = detail.Code
	case "QueryOrder":

		// order lookup
		var orderInquiry trip.OrderInquiry
		err := json.Unmarshal([]byte(req.DecryptedBody), &orderInquiry)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			resp.Code = "500"
			return resp, err
		}

		var inquiry yanolja.OrderConfirmation
		inquiry.OrderId, err = strconv.ParseInt(orderInquiry.SupplierOrderID, 10, 64)
		if err != nil {
			level.Error(logger).Log("error :", err)
			resp.Code = "500"
			return resp, err
		}

		resp.Body, err = s.mongoRepository[config.Instance().MongoDBName].FindOrderByOrderIdAndPartnerOrderId(ctx, inquiry.OrderId, orderInquiry.OTAOrderID)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				level.Error(logger).Log("repository error", "no record exist based on orderid and partnerorderid  ", err)
				resp.Code = "404"
				return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "FindOrderByOrderIdAndPartnerOrderId")
			}
			level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
			resp.Code = "500"
			resp.Body = err.Error()
			return resp, customError.NewError(ctx, "leisure-api-1015", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "FindOrderByOrderIdAndPartnerOrderId")
		}

		tripresp, err := MakeTripResponse(ctx, logger, resp.Body)
		if err != nil {
			return resp, err
		}

		for i, item := range tripresp.Items.Items {
			itemstr, err := cacheLayer.Get(ctx, item.PLU)
			if err != nil {
				level.Error(logger).Log("cache error ", err)
				return resp, err
			}
			tripresp.Items.Items[i].PLU = itemstr
		}

		resp.Body = tripresp

	default:
		level.Error(logger).Log("Unknown serviceName", serviceName)
		flag = false

	}
	if !flag {
		return resp, err
	}
	level.Info(logger).Log("response to trip ", resp)
	resp.Code = "200"
	endtime := time.Now()
	fmt.Println("**************servicename : ", serviceName, "   ***************", " Total time  : ", serviceName, endtime.Sub(startTime))
	return resp, nil
}

// GetProductsFromGGT  return  All product to channel based on listed productId
func (s *service) GetProductsFromGGT(ctx context.Context) (resp yanolja.Response, err error) {
	logger := log.With(
		s.logger,
		"method", "GetProductsFromGGT",
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

	products, err := s.mongoRepository[config.Instance().MongoDBName].GetAllProductViews(ctx)
	if err != nil {
		level.Error(logger).Log("repository error", "no product exist based on productId ", "error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, resp.Code, fmt.Sprintf("repository error on fetching product by productId, %v", err), "GetProductByProductId")

	}
	resp.Body = products
	resp.Code = "200"

	return resp, nil
}

// GetProductByIdFromGGT  return product to channel based on productId
func (s *service) GetProductByIdFromGGT(ctx context.Context, productId int64) (resp yanolja.Response, err error) {
	logger := log.With(
		s.logger,
		"method", "GetProductByIdFromGGT",
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

	product, err := s.mongoRepository[config.Instance().MongoDBName].GetProductViewByProductId(ctx, productId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no product exist based on productId ", "error ", err)
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetProductByProductId")
		}
		level.Error(logger).Log("repository error", "no product exist based on productId ", "error ", err)
		resp.Code = "500"
		return resp, customError.NewError(ctx, resp.Code, fmt.Sprintf("repository error on fetching product by productId, %v", err), "GetProductByProductId")

	}

	var productList = make([]domain.ProductView, 0)
	productList = append(productList, product)

	resp.Body = productList
	resp.Code = "200"
	return resp, nil
}

func MakeTripResponse(ctx context.Context, logger log.Logger, OrderModel interface{}) (tripresp trip.TripResponse, err error) {
	// Initialize the cache layer
	level.Info(logger).Log("function name  : ", MakeTripResponse)

	order, _ := OrderModel.(domain.Model)
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		return tripresp, err
	}
	// Convert int64 to string
	orderIdstr := strconv.FormatInt(order.OrderId, 10)

	itemstr, err := cacheLayer.Get(ctx, orderIdstr)
	if err != nil {
		level.Error(logger).Log("cache error ", err)
		return tripresp, err
	}

	if itemstr != "" {
		// Create a variable of the struct type
		var itemDetails domain.ItemIdDetails

		// Unmarshal the JSON string into the struct
		err = json.Unmarshal([]byte(itemstr), &itemDetails)
		if err != nil {
			level.Error(logger).Log("error unmarshalling JSON: %v", err)
			return tripresp, err
		}

		tripresp = trip.TripResponse{
			Items: itemDetails,
			Order: order,
		}
	} else {
		tripresp = trip.TripResponse{
			Order: order,
		}
	}

	return tripresp, nil
}

func GetVariantIdFromPlu(ctx context.Context, logger log.Logger, plu string) (variantId int64, productversion int32, productId int64, err error) {
	cacheLayer, err := cache.New(config.Instance().CacheName)
	if err != nil {
		level.Error(logger).Log("error", fmt.Sprintf("Error initializing cache layer: %s", err))
		return -1, -1, -1, err
	}

	pluVal, err := cacheLayer.Get(ctx, plu)
	if err != nil {
		level.Error(logger).Log("cache error ", err)
		return -1, -1, -1, err
	}

	detail := strings.Split(pluVal, "|") // ProductID-ProductVersion-VariantID
	variantId, err = strconv.ParseInt(detail[2], 10, 64)
	if err != nil {
		level.Error(logger).Log("error :", err)
		return -1, -1, -1, err
	}
	productId, err = strconv.ParseInt(detail[0], 10, 64)
	if err != nil {
		level.Error(logger).Log("error :", err)
		return 1, -1, -1, err
	}

	intConvrt, err := strconv.ParseInt(detail[1], 10, 64)
	if err != nil {
		level.Error(logger).Log("error :", err)
		return 1, -1, -1, err
	}

	productversion = int32(intConvrt)

	return
}

func VoucherDataCreator(ctx context.Context, logger log.Logger, s *service, data interface{}) {
	order := data.(domain.Model)
	var pdfVoucherReqList = make([]common.PdfVoucherRequest, 0)

	for _, orderVariant := range order.OrderVariants {

		product, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, orderVariant.ProductID)
		if err != nil {
			level.Error(logger).Log("database error", "retrieving product by productId", err)
			_, err = s.PostForcelyCancelOrder(ctx, order.PartnerOrderID)
			if err != nil {
				level.Error(logger).Log("error", "yanolja forced cancellation error")
			}
			return
		}

		var date, timeStr string

		for _, productOption := range product.ProductOptionGroups {
			if productOption.IsSchedule == true && productOption.IsRound == true {
				date = orderVariant.Date
				timeStr = orderVariant.Time

			} else if productOption.IsSchedule == true && productOption.IsRound == false {
				date = orderVariant.Date
			} else {
				date = ""
				timeStr = ""
			}

		}
		for _, variantItem := range orderVariant.OrderVariantItems {
			pdfVoucherReq := common.PdfVoucherRequest{
				OrderId:              order.OrderId,
				OrderVariantID:       orderVariant.OrderVariantID,
				OrderVariantItemID:   variantItem.OrderVariantItemID,
				OrderVariantItemName: variantItem.OrderVariantItemName,
				VariantName:          orderVariant.VariantName,
				VariantID:            orderVariant.VariantID,
				ProductID:            product.ProductID,
				ProductName:          product.ProductName,
				ActualCustomerName:   order.ActualCustomer.Name,
				PurchaseDate:         order.UpdatedAt,
				VisitingDate:         date,
				VisitingTime:         timeStr,
				ConfirmationNumber:   order.PartnerOrderID,
				ValidityPeriod:       orderVariant.ValidityPeriod.StartDateTime.Format("2006-01-02") + " to " + orderVariant.ValidityPeriod.EndDateTime.Format("2006-01-02"),
				ProductInfo:          product.ProductInfo,
			}
			pdfVoucherReqList = append(pdfVoucherReqList, pdfVoucherReq)
		}
	}

	UpdateVoucherPdf(ctx, logger, s, pdfVoucherReqList, order.PartnerOrderID)

}

// UpdateVoucherPdf call voucher generator to generate voucher and store in order
func UpdateVoucherPdf(ctx context.Context, logger log.Logger, s *service, voucherReq []common.PdfVoucherRequest, partnerOrderId string) {
	level.Info(logger).Log("request data UpdateVoucherPdf", voucherReq)

	pdfVoucherSvc, _ := pdfvoucher.New(ctx)
	res, err := pdfVoucherSvc.NotifyToVoucherPdfUpdate(ctx, voucherReq)
	level.Info(logger).Log("error in pdf voucher ", err)
	if err != nil {
		level.Error(logger).Log("error", "pdf voucher call error ")
		_, err := s.PostForcelyCancelOrder(ctx, partnerOrderId)
		if err != nil {
			level.Error(logger).Log("error", "yanolja forced cancellation error")
			return
		}
		return
	}
	level.Info(logger).Log("response from pdfVoucher generator", res.Body, reflect.TypeOf(res.Body))

	var respBody map[string]interface{}
	if err := json.Unmarshal([]byte(res.Body.(string)), &respBody); err != nil {
		level.Error(logger).Log("Error decoding JSON:", err)

		_, err := s.PostForcelyCancelOrder(ctx, partnerOrderId)
		if err != nil {
			level.Error(logger).Log("error", "yanolja forced cancellation error")
			return
		}
		return
	}

	level.Info(logger).Log("response from pdfVoucher generator", res)

	if res.Code == "200" {
		orderIDRaw, _ := respBody["orderId"]
		orderIDFloat, _ := orderIDRaw.(float64)
		orderId := int64(orderIDFloat)

		order, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, orderId)
		if err != nil {
			level.Error(logger).Log("error", "repository error on GetOrderbyOrderId")
			_, err := s.PostForcelyCancelOrder(ctx, partnerOrderId)
			if err != nil {
				level.Error(logger).Log("error", "yanolja forced cancellation error")
				return
			}
			return
		}

		level.Info(logger).Log("order data ", order)

		level.Info(logger).Log("info ", "call to trip")

		// call to trip
		level.Info(logger).Log("info ", "$$$$$$$$$$$$$$$$$$$$$ call to trip to process voucher #########################")
		processTripNotification(ctx, logger, s, "VoucherUpdateNotify", order, constant.TRIPPAYMENTREQUEST)
	}

}

// IsPdfVoucher for checking PDF voucher condition for all variant.
// on assumption that  all variant has pdf voucher
func IsPdfVoucher(ctx context.Context, logger log.Logger, s *service, data interface{}) bool {
	logger.Log("IsPdfVoucher called")
	pdfVoucherCallFlag := false

	order := data.(domain.Model)
	for _, orderVariant := range order.OrderVariants {
		product, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductByProductId(ctx, orderVariant.ProductID)
		if err != nil {
			level.Error(logger).Log("database error", "retrieving product by productId", err)
			_, err = s.PostForcelyCancelOrder(ctx, order.PartnerOrderID)
			if err != nil {
				level.Error(logger).Log("error", "yanolja forced cancellation error")
			}
		}
		logger.Log("Product data ", product)

		var isProductPdfTypeVoucherFlag bool = false

		for _, productOption := range product.ProductOptionGroups {
			for _, productVariant := range productOption.Variants {
				for _, item := range productVariant.VariantItems {
					logger.Log("VoucherDisplayTypeCode", item.VoucherDisplayTypeCode)
					if item.VoucherDisplayTypeCode == "NONE" {
						isProductPdfTypeVoucherFlag = true
						continue
					} else {
						isProductPdfTypeVoucherFlag = false
						break
					}

				}

			}
			logger.Log("isProductPdfTypeVoucherFlag", isProductPdfTypeVoucherFlag)
		}

		// single orderVariantItem is available
		for _, variantItem := range orderVariant.OrderVariantItems {
			logger.Log("VoucherProvideStatusCode", variantItem.Voucher.VoucherProvideStatusCode, "isProductPdfTypeVoucherFlag in order loop", isProductPdfTypeVoucherFlag)
			if variantItem.Voucher.VoucherProvideStatusCode == "NON_PROVIDED" && isProductPdfTypeVoucherFlag {
				pdfVoucherCallFlag = true
				continue
			} else {
				pdfVoucherCallFlag = false
				break
			}

		}

	}
	logger.Log("pdfVoucherCallFlag", pdfVoucherCallFlag)

	return pdfVoucherCallFlag
}
