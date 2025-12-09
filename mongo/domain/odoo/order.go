package odoo

// Order represents the MongoDB model for the given dataset
type Order struct {
	Id                     string         `bson:"_id" json:"id,omitempty"`
	Channel                string         `bson:"channel" json:"channel"`
	Supplier               string         `bson:"supplier" json:"supplier"`
	OrderID                int64          `bson:"orderId" json:"orderId"`
	PartnerOrderID         string         `bson:"partnerOrderId" json:"partnerOrderId"`
	BookingStatus          string         `bson:"bookingStatus" json:"bookingStatus"`
	CustomerName           string         `bson:"customerName" json:"customerName"`
	BookedAt               string         `bson:"bookedAt" json:"bookedAt" default:""` ///createdAt
	ExpirationStatus       string         `bson:"expirationStatus" json:"expirationStatus"`
	TotalQuantity          int            `bson:"totalQuantity" json:"totalQuantity"`
	SupplierCostPrice      float32        `bson:"supplierCostPrice" json:"supplierCostPrice"` // sum of all cost price of ord
	SupplierSalePrice      float32        `bson:"supplierSalePrice" json:"supplierSalePrice"`
	Margin                 float32        `bson:"margin" json:"margin"`
	MarkupType             string         `bson:"markupType" json:"markupType"`
	MarkupValue            float32        `bson:"markupValue" json:"markupValue"`
	ChannelNetPriceFormula string         `bson:"channelNetPriceFormula" json:"channelNetPriceFormula"`
	ChannelNetPrice        float32        `bson:"channelNetPrice" json:"channelNetPrice"`
	SupplierCurrency       string         `bson:"supplierCurrency" json:"supplierCurrency"`
	AppliedMarkup          float32        `bson:"appliedMarkup" json:"appliedMarkup"`
	Variants               []OrderVariant `bson:"variants" json:"variants"`
	MarkupInfo             []MarkupDetail `bson:"markupInfo" json:"markupInfo"`
	NetPriceInfo           NetPriceDetail `bson:"netPriceInfo" json:"netPriceInfo"`
	CreatedAt              string         `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt              string         `bson:"updatedAt" json:"updatedAt"`
}

type OrderVariant struct {
	OrderVariantID             int64   `bson:"orderVariantId" json:"orderVariantId"`
	ProductId                  int64   `bson:"productId" json:"productId"`
	ProductVersion             int     `bson:"productVersion" json:"productVersion"`
	ProductName                string  `bson:"productName" json:"productName"`
	VariantId                  int64   `bson:"variantId" json:"variantId"`
	VariantName                string  `bson:"variantname" json:"variantName"`
	VariantSalePrice           float32 `bson:"variantSalePrice" json:"variantSalePrice"`
	Quantity                   int     `bson:"quantity" json:"quantity"`
	MarkupType                 string  `bson:"markupType" json:"markupType"`
	MarkupValue                float32 `bson:"markupValue" json:"markupValue"`
	VariantCostPrice           float32 `bson:"variantCostPrice" json:"variantCostPrice"`
	VisitDate                  string  `bson:"visitDate" json:"visitDate"`
	VisitTime                  string  `bson:"visitTime" json:"visitTime"`
	Currency                   string  `bson:"currency" json:"currency"`
	VariantStatusType          string  `bson:"variantStatusType" json:"variantStatusType"`
	RefundCode                 string  `bson:"refundCode" json:"refundCode"`
	RefundStatus               string  `bson:"refundStatus" json:"refundStatus"`
	OrderUsedOrRestoreDateTime string  `bson:"orderUsedDateTime,omitempty" json:"orderUsedDateTime,omitempty"`
	CanceledDateTime           string  `bson:"canceledDateTime,omitempty" json:"canceledDateTime" default:""`
	CancelStatusCode           string  `bson:"cancelStatusCode" json:"cancelStatusCode" validate:"required" default:""` // New field with default value
	RefundInfo                 string  `bson:"refundInfo" json:"refundInfo" validate:"required" default:""`
	OrderCancelTypeCode        string  `bson:"orderCancelTypeCode" json:"orderCancelTypeCode" validate:"required"`
	CancelRejectTypeCode       string  `bson:"cancelRejectTypeCode" json:"cancelRejectTypeCode" validate:"required" default:""`
	Message                    string  `bson:"message,omitempty" json:"message,omitempty"`                                      // check for use
	CancelFailReasonCode       string  `bson:"cancelFailReasonCode" json:"cancelFailReasonCode" validate:"required" default:""` // New field with default value
	ForceCancelTypeCode        string  `bson:"forceCancelTypeCode" json:"forceCancelTypeCode" validate:"required"`
	OrderVariantItemID         int64   `bson:"orderVariantItemId" json:"orderVariantItemId" validate:"required"`
	OrderVariantItemName       string  `bson:"orderVariantItemName" json:"orderVariantItemName" validate:"required"`
	VoucherCode                string  `bson:"voucherCode" json:"voucherCode" validate:"required"`
	VoucherDisplayTypeCode     string  `bson:"voucherDisplayTypeCode" json:"voucherDisplayTypeCode" validate:"required"`
	VoucherPdfUrl              string  `bson:"voucherPdfUrl" json:"voucherPdfUrl" default:""`
}

type MarkupDetail struct {
	ProductId      int64   `bson:"productId" json:"productId"`
	VariantId      int64   `bson:"variantId" json:"VariantId"`
	Quantity       int32   `bson:"quantity" json:"quantity"`
	CostPrice      float64 `bson:"costPrice" json:"costPrice"`
	TotalCostPrice float64 `bson:"totalCostPrice" json:"totalCostPrice"` // costprice*quantity
	TotalSalePrice float64 `bson:"totalSalePrice" json:"totalSalePrice"` // saleprice*quantity
	SalePrice      float64 `bson:"salePrice" json:"salePrice"`
	MarkupValue    float32 `bson:"markupValue" json:"markupValue"`
	MarkupType     string  `bson:"markupType" json:"markupType"`
}

type NetPriceDetail struct {
	NetCostPrice float32 `bson:"netCostPrice" json:"netCostPrice"`
	NetSalePrice float32 `bson:"netSalePrice" json:"netSalePrice"`
	MarkupValue  float32 `bson:"markupValue" json:"markupValue"`
	NetQuantity  int     `bson:"netQuantity" json:"netQuantity"`
	MarkupType   string  `bson:"markupType" json:"markupType"`
	Currency     string  `bson:"currency" json:"currency"`
}
