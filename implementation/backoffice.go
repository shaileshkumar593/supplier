package implementation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"swallow-supplier/config"
	customError "swallow-supplier/error"
	"swallow-supplier/mongo/domain/odoo"
	domain "swallow-supplier/mongo/domain/yanolja"
	"swallow-supplier/request_response/common"
	"swallow-supplier/request_response/yanolja"
	"swallow-supplier/utils"
	"swallow-supplier/utils/constant"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetProductViewSync  to sync product_view to odoo
func (s *service) GetProductSync(ctx context.Context) (resp common.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetProductView",
		"Request ID", requestID,
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

	products, err := s.mongoRepository[config.Instance().MongoDBName].FetchProductsBasedOnOdooSyncStatus(ctx)
	if err != nil {
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "", http.StatusInternalServerError, "GetOrderbyOrderId")
	}
	level.Info(logger).Log("info", "document inserted with id ", products)

	numberOfProduct := len(products)
	if numberOfProduct <= 0 {
		level.Error(logger).Log("error ", fmt.Sprintf("number of product is %d ", numberOfProduct))
		resp.Code = "404"
		resp.Body = fmt.Sprintf("number of product is %d ", numberOfProduct)
		return resp, customError.NewErrorCustom(ctx, "404", resp.Body.(string), "", http.StatusNotFound, "GetProductSync")
	}

	var odooProducts = make([]odoo.Product, 0)

	for _, product := range products {
		var addressDetail odoo.FacilityLocationDetail
		var productVariant odoo.ProductVariant
		var odooproduct odoo.Product

		odooproduct.ProductID = product.ProductID
		odooproduct.ProductName = product.ProductName
		odooproduct.ProductVersion = product.ProductVersion
		odooproduct.ProductStatusCode = product.ProductStatusCode
		odooproduct.ProductTypeCode = product.ProductTypeCode
		if product.ProductStatusCode == "IN_SALE" && product.IsUsed {
			odooproduct.ProductValidity = "InValid"
		} else {
			odooproduct.ProductValidity = "Not_Valid"
		}

		facilityInfo := product.ProductInfo.FacilityInfos
		addressDetails := make([]odoo.FacilityLocationDetail, 0)

		if len(facilityInfo) > 0 {
			for _, facility := range facilityInfo {
				addressDetail.Latitude = facility.Location.Latitude
				addressDetail.Longitude = facility.Location.Longitude
				addressDetail.Address = facility.Location.Address
				addressDetails = append(addressDetails, addressDetail)
			}
		}
		odooproduct.FacilityAddress = addressDetails

		var regionInfos []odoo.Region
		if len(product.Regions) > 0 {
			regionInfos = FlattenRegions(product.Regions)
		}

		if len(regionInfos) > 0 {
			odooproduct.Regions = regionInfos
		}

		odooproduct.IsIntegratedVoucher = product.IsIntegratedVoucher

		productVariants := make([]odoo.ProductVariant, 0)
		for _, optiongrp := range product.ProductOptionGroups {

			for _, variant := range optiongrp.Variants {
				productVariant.VariantID = variant.VariantID
				productVariant.ProductID = variant.ProductID
				if !strings.EqualFold(variant.VariantName, "Adult") || !strings.EqualFold(variant.VariantName, "Youth") ||
					!strings.EqualFold(variant.VariantName, "Child") || !strings.EqualFold(variant.VariantName, "children") {
					productVariant.VariantName = "N/A"
				}
				productVariant.VariantName = variant.VariantName
				productVariant.VariantDescription = variant.VariantDescription
				productVariant.RefundApprovalTypeCode = variant.RefundApprovalTypeCode
				productVariant.IsRefundableAfterExpiration = variant.IsRefundableAfterExpiration
				productVariant.RefundInfo = variant.RefundInfo
				productVariant.VariantStatusCode = variant.VariantStatusCode
				productVariant.OrderExpirationUsageProcessTypeCode = variant.OrderExpirationUsageProcessTypeCode
				productVariant.OrderExpirationDateTypeCode = variant.OrderExpirationDateTypeCode
				productVariant.ValidityStartDate = variant.SalePeriod.StartDateTime
				productVariant.ValidityEndDate = variant.SalePeriod.EndDateTime
				productVariant.Currency = variant.Price.Currency
				productVariant.SupplierCostPrice = variant.Price.CostPrice
				productVariant.SalePrice = variant.Price.SalePrice
				productVariant.RetailPrice = variant.Price.RetailPrice
				productVariant.DiscountSalePrice = variant.Price.DiscountSalePrice
				productVariant.IsRound = optiongrp.IsRound
				productVariant.IsSchedule = optiongrp.IsSchedule
				voucherdisplaycode := make([]string, 0)
				for _, item := range variant.VariantItems {
					voucherdisplaycode = append(voucherdisplaycode, item.VoucherDisplayTypeCode)
				}
				productVariant.VoucherDisplayCode = voucherdisplaycode
				productVariants = append(productVariants, productVariant)
			}
		}
		odooproduct.Variants = productVariants
		odooProducts = append(odooProducts, odooproduct)

	}

	upsertedProducts, err := s.mongoRepository[config.Instance().MongoDBName].UpsertOdooProduct(ctx, odooProducts)
	if err != nil {
		level.Error(logger).Log("repository error", "upsert record of odoo product  ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on upsert odoo order, %v", err), "", http.StatusInternalServerError, "UpsertOdooProduct")
	}

	resp.Body = upsertedProducts
	resp.Code = "200"
	resp.Status = http.StatusOK

	return resp, nil
}

// GetOrderSync to sync order to odoo
func (s *service) GetOrderSync(ctx context.Context) (resp common.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetOrderSync",
		"Request ID", requestID,
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

	orders, err := s.mongoRepository[config.Instance().MongoDBName].GetOrdersByOdooSyncStatus(ctx)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = fmt.Errorf("no record exist based on odoosyncstatus false: ")
			resp.Code = "200"
			resp.Status = http.StatusNoContent
			resp.Body = []domain.Model{}
			return resp, nil
		}

		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "", http.StatusInternalServerError, "GetOrderbyOrderId")
	}
	level.Info(logger).Log("info", "document inserted with id ", orders)

	var odooorder = make([]odoo.Order, 0)
	for _, order := range orders {

		odooRecord := oderDataForOdooSync(ctx, s, logger, order)
		odooorder = append(odooorder, odooRecord)
	}

	records, err := s.mongoRepository[config.Instance().MongoDBName].UpsertOdooOrder(ctx, odooorder)
	if err != nil {
		level.Error(logger).Log("repository error", "upsert record of odoo order  ", err)
		resp.Code = "500"
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on upsert odoo order, %v", err), "", http.StatusInternalServerError, "UpsertOdooProduct")
	}

	resp.Body = records
	resp.Code = "200"
	resp.Status = http.StatusOK
	return resp, nil
}

// GetOrderSyncByOrderIdInOdoo   get order sync in odoo_order
func (s *service) GetOrderSyncByOrderIdInOdoo(ctx context.Context, req yanolja.Order) (resp common.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "GetOrderSyncByOrderIdInOdoo",
		"Request ID", requestID,
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

	orderRec, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, int64(req.OrderId))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no record exist based on orderId ", err)
			err = fmt.Errorf("no document exist with orderId: %d", int64(req.OrderId))
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderbyOrderId")
		}
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		resp.Status = http.StatusInternalServerError
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "", http.StatusInternalServerError, "GetOrderbyOrderId")
	}
	level.Info(logger).Log("info", "document inserted with id ", orderRec.OrderId)

	if orderRec.OodoSyncStatus == false {
		var odooorder = make([]odoo.Order, 0)
		odooRecord := oderDataForOdooSync(ctx, s, logger, orderRec)
		odooorder = append(odooorder, odooRecord)
		odooorders, err := s.mongoRepository[config.Instance().MongoDBName].UpsertOdooOrder(ctx, odooorder)

		if err != nil {
			level.Error(logger).Log("repository error", "upsert record of odoo order  ", err)
			resp.Code = "500"
			return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on upsert odoo order, %v", err), "", http.StatusInternalServerError, "UpsertOdooProduct")
		}
		resp.Body = odooorders
		resp.Code = "204"
		resp.Status = http.StatusNoContent
		return resp, nil

	}

	odoorec, err := s.mongoRepository[config.Instance().MongoDBName].GetOdooOrderbyOrderId(ctx, int64(req.OrderId))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(logger).Log("repository error", "no record exist based on orderId ", err)
			err = fmt.Errorf("no document exist with orderId: %d", int64(req.OrderId))
			resp.Code = "404"
			return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOdooOrderbyOrderId")
		}
		level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
		resp.Code = "500"
		resp.Status = http.StatusInternalServerError
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "", http.StatusInternalServerError, "GetOdooOrderbyOrderId")
	}
	resp.Body = odoorec
	resp.Code = "204"
	resp.Status = http.StatusNoContent
	return resp, nil
}

func (s *service) PostOdooSyncStatus(ctx context.Context, req yanolja.OrderList) (resp common.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostOdooSyncStatus",
		"Request ID", requestID,
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

	odooSyncStatusMap := map[string]any{
		"oodoSyncStatus": true,
		"updatedAt":      time.Now().UTC().Format(time.RFC3339),
	}
	for _, orderid := range req.OrderId {
		_, err := s.mongoRepository[config.Instance().MongoDBName].GetOrderbyOrderId(ctx, orderid)
		if err != nil {

			if err == mongo.ErrNoDocuments {
				level.Error(logger).Log("repository error", "no record exist based on orderId ", err)
				resp.Code = "404"
				resp.Status = http.StatusBadRequest
				return resp, customError.NewErrorCustom(ctx, resp.Code, err.Error(), "", http.StatusNotFound, "GetOrderbyOrderId")
			}
			level.Error(logger).Log("repository error", "fetching record based on orderId ", err)
			resp.Code = "500"
			resp.Body = err
			resp.Status = http.StatusInternalServerError
			return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on fetching order by orderId, %v", err), "", http.StatusInternalServerError, "GetOrderbyOrderId")

		}

		level.Info(logger).Log("info ", fmt.Sprintf("record  exist with orderid %d", orderid))

		_, err = s.mongoRepository[config.Instance().MongoDBName].UpdateOrderByOrderId(ctx, orderid, odooSyncStatusMap)
		if err != nil {
			level.Error(logger).Log("repository error", "error in updating order for oodoSyncStatus", err)
			resp.Code = "500"
			resp.Body = err
			resp.Status = http.StatusInternalServerError
			return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error in updating order for oodoSyncStatus, %v", err), "", http.StatusInternalServerError, "UpdateOrderByOrderId")
		}

	}

	resp.Code = "204"
	resp.Body = "success"
	resp.Status = http.StatusNoContent
	return resp, nil
}

// PostChangeOdooSyncStatusToFalse
func (s *service) PostChangeOdooSyncStatusToFalse(ctx context.Context) (resp common.Response, err error) {
	requestID := utils.GenerateUUID("GGT", true)

	logger := log.With(
		s.logger,
		"method", "PostOdooSyncStatus",
		"Request ID", requestID,
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

	err = s.mongoRepository[config.Instance().MongoDBName].UpdateOdooSyncStatusToFalseIfTrue(ctx)

	if err != nil {
		level.Error(logger).Log("repository error", "error in updating odooSyncStatus flag of order ", "error ", err)
		resp.Code = "500"
		resp.Body = err
		resp.Status = http.StatusInternalServerError
		return resp, customError.NewErrorCustom(ctx, "500", fmt.Sprintf("repository error on updating odooSyncStatus flag of order, %v", err), "", http.StatusInternalServerError, "UpdateOdooSyncStatusToFalseIfTrue")

	}
	resp.Code = "204"
	resp.Body = "success"
	resp.Status = http.StatusNoContent

	return resp, nil
}

func FlattenRegions(regions []domain.Regional) []odoo.Region {
	var result = make([]odoo.Region, 0)
	for _, region := range regions {
		result = append(result, flattenRegion(region, odoo.Region{})...)
	}
	fmt.Println("done")
	return result
}

func flattenRegion(region domain.Regional, parent odoo.Region) []odoo.Region {
	// Set region name according to level
	switch region.RegionLevel {
	case 1:
		parent.Continent = region.RegionName
	case 2:
		parent.Country = region.RegionName
	case 3:
		parent.City = region.RegionName
	case 4:
		parent.Area = region.RegionName
	case 5:
		parent.AreaCode = region.RegionName
	}

	// If this is a leaf node, return the final RegionInfo
	if len(region.SubRegions) == 0 {
		return []odoo.Region{parent}
	}

	// Recurse into subregions
	var result []odoo.Region
	for _, subRegion := range region.SubRegions {
		result = append(result, flattenRegion(subRegion, parent)...)
	}
	return result
}

func oderDataForOdooSync(ctx context.Context, s *service, logger log.Logger, order domain.Model) odoo.Order {
	var odooRec odoo.Order

	odooRec.OrderID = order.OrderId
	odooRec.Supplier = order.Suppliers
	odooRec.Channel = order.PartnerOrderChannelCode
	if order.OrderExpired {
		odooRec.ExpirationStatus = constant.EXPIREDSTAUS
	} else {
		odooRec.ExpirationStatus = constant.VALIDSTAUS
	}
	if order.OrderStatusCode == constant.ORDERDONE {
		odooRec.BookingStatus = constant.ORDERCOMPLETE
	} else if order.OrderStatusCode == constant.ORDERPREPARE {
		odooRec.BookingStatus = constant.ORDERNOTCOMPLETE
	}
	odooRec.BookedAt = order.CreatedAt
	odooRec.CreatedAt = order.CreatedAt
	odooRec.UpdatedAt = order.UpdatedAt
	odooRec.CustomerName = order.Customer.Name
	odooRec.PartnerOrderID = order.PartnerOrderID

	/* odooRec.MarkupValue = constant.MARKUPPERCENTAGE
	odooRec.MarkupType = constant.PERCENTAGE */
	var variant odoo.OrderVariant

	variantlist := make([]odoo.OrderVariant, 0)

	netInfo, markuplist := FindCurrencyCodeAndMarkup(order.SelectVariants, odooRec.Channel)
	odooRec.SupplierCurrency = netInfo.Currency

	for _, ov := range order.OrderVariants {

		variant.Currency = odooRec.SupplierCurrency
		variant.VariantId = ov.VariantID
		variant.ProductId = ov.ProductID

		productName, err := s.mongoRepository[config.Instance().MongoDBName].GetProductNameByProductID(ctx, variant.ProductId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				level.Error(logger).Log("repository error", "no record exist based on productId ", err)

				return odoo.Order{}
			}
			level.Error(logger).Log("repository error", "fetching record based on productId ", err)

			return odoo.Order{}

		}
		variant.ProductName = productName
		variant.ProductVersion = int(ov.ProductVersion)
		variant.OrderVariantID = ov.OrderVariantID
		if ov.Date != "" {
			variant.VisitDate = ov.Date
		} else {
			variant.VisitDate = ""
		}
		if ov.Time != "" {
			variant.VisitTime = ov.Time
		} else {

			variant.VisitTime = ""
		}

		variant.VariantName = ov.VariantName

		detail := GetVariantPrice(variant.VariantId, markuplist)
		variant.MarkupType = detail.MarkupType
		variant.MarkupValue = detail.MarkupValue

		variant.VariantCostPrice = float32(detail.CostPrice)
		variant.VariantSalePrice = float32(detail.SalePrice)

		if ov.OrderVariantStatusTypeCode == constant.ORDERVARIANTUSEDSTATUS {
			variant.OrderUsedOrRestoreDateTime = ov.UsedRestoreDateTime
		}
		if ov.OrderVariantStatusTypeCode == constant.ORDERVARIANTCANCELEDSTATUS {
			variant.CancelStatusCode = ov.CancelStatusCode
			variant.CanceledDateTime = ov.CanceledDateTime
		} else if ov.OrderVariantStatusTypeCode == constant.ORDERVARIANTCANCELINGSTATUS {
			variant.CancelStatusCode = ov.CancelStatusCode
			variant.CancelFailReasonCode = ov.CancelFailReasonCode
			if ov.CancelRejectTypeCode != "" {
				variant.CancelRejectTypeCode = ov.CancelRejectTypeCode
				variant.Message = ov.Message

			}

		}

		variant.VariantName = ov.VariantName
		variant.Quantity = 1
		variant.RefundCode = ov.RefundApprovalTypeCode
		if variant.RefundCode == "ADMIN" {
			variant.RefundStatus = strings.ToUpper("Administrator approval refund")
		} else if variant.RefundCode == "DIRECT" {
			variant.RefundStatus = "INSTANT REFUND"
		}

		variant.VariantStatusType = ov.OrderVariantStatusTypeCode
		// //var voucher = make([]odoo.VoucherInfo, 0)
		// var voucherDetail odoo.VoucherInfo

		for _, item := range ov.OrderVariantItems {
			variant.OrderVariantItemID = item.OrderVariantItemID
			variant.OrderVariantItemName = item.OrderVariantItemName
			if item.Voucher.VoucherDisplayTypeCode == "" {
				variant.VoucherDisplayTypeCode = "NONE"
				variant.VoucherCode = ""
				variant.VoucherPdfUrl = item.Voucher.VoucherPdfUrl
			} else {
				variant.VoucherDisplayTypeCode = item.Voucher.VoucherDisplayTypeCode
				variant.VoucherCode = item.Voucher.VoucherCode
				variant.VoucherPdfUrl = ""

			}

		}
		//variant.Voucher = voucherDetail
		variantlist = append(variantlist, variant)

	}

	odooRec.Variants = variantlist

	odooRec.SupplierCostPrice = float32(netInfo.NetCostPrice) //float32(costprice)
	odooRec.SupplierSalePrice = float32(netInfo.NetSalePrice) //float32(saleprice)
	odooRec.MarkupType = netInfo.MarkupType
	odooRec.MarkupValue = netInfo.MarkupValue
	odooRec.AppliedMarkup = odooRec.MarkupValue
	odooRec.Margin = TotaMargineCompute(markuplist, order.SelectVariants)
	odooRec.ChannelNetPriceFormula = "supplierCostPrice+margin"
	odooRec.ChannelNetPrice = odooRec.SupplierCostPrice + odooRec.Margin

	//odooRec.ChannelNetPrice = float32(costprice) + odooRec.Margin
	odooRec.TotalQuantity = netInfo.NetQuantity
	odooRec.NetPriceInfo = netInfo
	odooRec.MarkupInfo = markuplist
	return odooRec
}

func FindCurrencyCodeAndMarkup(selectVariants []domain.SelectVariant, channelCode string) (netPriceInfo odoo.NetPriceDetail, markuplst []odoo.MarkupDetail) {
	var totalcostprice float64 = 0.0
	var totalsaleprice float64 = 0.0
	var orderQuantity int = 0
	currency := make([]string, 0)
	var currencyCode string
	var MarkupList = make([]odoo.MarkupDetail, 0)
	var markupstruct odoo.MarkupDetail

	for _, selectvariant := range selectVariants {
		totalcostprice = 0.0
		totalsaleprice = 0.0
		if len(currency) == 0 {
			currency = append(currency, selectvariant.Currency)
		} else {
			for _, currencyval := range currency {
				if currencyval == selectvariant.Currency {
					continue
				} else {
					currency = append(currency, selectvariant.Currency)
				}
			}
		}
		totalcostprice = totalcostprice + float64(selectvariant.CostPrice*float32(selectvariant.Quantity))
		totalsaleprice := totalsaleprice + float64(selectvariant.PartnerSalePrice*float32(selectvariant.Quantity))
		markupstruct.SalePrice = float64(selectvariant.PartnerSalePrice)
		markupstruct.CostPrice = float64(selectvariant.CostPrice)
		markupstruct.VariantId = selectvariant.VariantID
		if channelCode == "GGT_EVERLAND" {
			price, markup, markuptype := MargineAndNetPriceDetail(selectvariant.ProductID, totalcostprice, selectvariant.Quantity)
			fmt.Println("========GGT_EVERLAND ======== ", price, markup, markuptype)
			markupstruct.MarkupType = markuptype
			markupstruct.MarkupValue = float32(markup)
			markupstruct.TotalCostPrice = float64(price)
			markupstruct.Quantity = selectvariant.Quantity
			markupstruct.TotalSalePrice = totalsaleprice
			markupstruct.ProductId = selectvariant.ProductID

		} else {
			markupstruct.MarkupType = constant.PERCENTAGE
			markupstruct.Quantity = selectvariant.Quantity
			markupstruct.MarkupValue = constant.MARKUPPERCENTAGE
			for _, margine := range constant.MARGINEDETAIL {
				if margine.ProductId == selectvariant.ProductID && margine.MargineType == constant.PERCENTAGE {
					markupstruct.MarkupValue = margine.Value
					break
				}
			}
			markupstruct.TotalCostPrice = totalcostprice + (totalcostprice * (float64(markupstruct.MarkupValue / 100)))
			markupstruct.TotalSalePrice = totalsaleprice
			markupstruct.ProductId = selectvariant.ProductID
		}

		orderQuantity += int(selectvariant.Quantity)
		MarkupList = append(MarkupList, markupstruct)
	}

	fmt.Println("^^^^^^^^^^^^^^^^ MarkupList ^^^^^^^^^^^^^^^^^^^^ ", MarkupList)
	if len(currency) == 1 {
		currencyCode = currency[0]
	} else {
		currencyCode = constant.DEFAULTCURRENCY
	}
	netPriceInfo.NetQuantity = orderQuantity
	netPriceInfo.Currency = currencyCode

	markuptype, markupval := FindMarkupTypeAndValue(MarkupList)
	netPriceInfo.MarkupType = markuptype
	netPriceInfo.MarkupValue = markupval

	netPriceInfo.NetCostPrice, netPriceInfo.NetSalePrice = CalculateTotalPrice(selectVariants)

	return netPriceInfo, MarkupList
}

func MargineAndNetPriceDetail(productId int64, cost float64, quantity int32) (netPrice float64, margineVal int, margineType string) {
	for _, margine := range constant.MARGINEDETAIL {
		margineVal = int(margine.Value)
		margineType = margine.MargineType
		if margine.ProductId == productId {
			switch margine.MargineType {
			case "FLATVALUE":
				netPrice = cost + float64(int(margine.Value)*int(quantity))
			case "PERCENTAGE":
				mergineval := margine.Value
				netPrice = cost + (cost * (float64(mergineval / 100)))
			}
			break
		}
	}

	return netPrice, margineVal, margineType
}

func FindMarkupTypeAndValue(data []odoo.MarkupDetail) (markuptype string, markupval float32) {
	if len(data) == 0 {
		return "", 0.0
	}

	if len(data) == 1 {
		markuptype = data[0].MarkupType
		markupval = data[0].MarkupValue
	}

	return markuptype, markupval
}

func CalculateTotalPrice(data []domain.SelectVariant) (costpriceTotal float32, salepriceTotal float32) {

	if len(data) == 0 {
		return 0, 0
	}
	var totalCost, totalSale float32 = data[0].CostPrice, data[0].PartnerSalePrice
	quantity := float32(data[0].Quantity)
	if len(data) == 1 {
		return totalCost * quantity, totalSale * quantity
	}

	/* for _, item := range data {
		totalCost += item.CostPrice * float32(item.Quantity)
		totalSale += item.PartnerSalePrice * float32(item.Quantity)
	}

	costpriceTotal = totalCost
	salepriceTotal = totalSale */
	return
}

func TotaMargineCompute(data []odoo.MarkupDetail, selectVariant []domain.SelectVariant) (totalMargine float32) {
	for i, markuplst := range data {
		if markuplst.ProductId == selectVariant[i].ProductID {
			if strings.ToLower(markuplst.MarkupType) == "percentage" {
				totalMargine += (selectVariant[i].CostPrice * float32(selectVariant[i].Quantity)) * (float32(markuplst.MarkupValue / 100))
			} else if strings.ToLower(markuplst.MarkupType) == "flatvalue" {
				totalMargine += float32(markuplst.MarkupValue) * float32(selectVariant[i].Quantity)
			}
		}
	}

	return
}

func GetVariantPrice(variantId int64, variantDetail []odoo.MarkupDetail) (markupdetail odoo.MarkupDetail) {

	for _, value := range variantDetail {
		if variantId == value.VariantId {
			markupdetail = value
			break
		}
	}

	return
}
