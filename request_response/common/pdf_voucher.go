package common

import (
	"swallow-supplier/mongo/domain/yanolja"
)

type PdfVoucherRequest struct {
	OrderId              int64               `bson:"orderId" json:"orderId" validate:"required"`
	OrderVariantID       int64               `bson:"orderVariantId" json:"orderVariantId" validate:"required"`
	OrderVariantItemID   int64               `bson:"orderVariantItemId" json:"orderVariantItemId" validate:"required"`
	OrderVariantItemName string              `bson:"orderVariantItemName" json:"orderVariantItemName" validate:"required"`
	VariantName          string              `bson:"variantName" json:"variantName" validate:"required"`
	VariantID            int64               `bson:"variantId" json:"variantId" validate:"required"`
	ProductID            int64               `bson:"productId" json:"productId" validate:"required"`
	ProductName          string              `bson:"productName" json:"productName" validate:"required"`
	ActualCustomerName   string              `bson:"actualCustomerName" json:"actualCustomerName" validate:"required"`
	PurchaseDate         string              `bson:"purchaseDate" json:"purchaseDate"`
	VisitingDate         string              `bson:"VisitingDate" json:"VisitingDate" validate:"required"`
	VisitingTime         string              `bson:"VisitingTime" json:"VisitingTime" validate:"required"`
	ConfirmationNumber   string              `bson:"partnerOrderId" json:"partnerOrderId" validate:"required"`
	ValidityPeriod       string              `bson:"ValidityPeriod" json:"ValidityPeriod" validate:"required"`
	ProductInfo          yanolja.ProductInfo `bson:"ProductInfo" json:"ProductInfo"`
}

type VoucherResponse struct {
	OrderId int64  `json:"orderId"`
	PdfUrl  string `json:"pdfUrl"`
}
