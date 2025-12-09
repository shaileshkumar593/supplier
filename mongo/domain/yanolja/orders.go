package yanolja

import (
	"time"
)

// Customer represents the customer information.
type Customer struct {
	Name  string `bson:"name" json:"name" validate:"required, min=2,max=30,regexp=^[a-zA-Z]{2,30}$"`
	Tel   string `bson:"tel" json:"tel" validate:"required, regexp=^(010|011|016|017|018|019|\\+82)[-\\s]?\\d{3,4}[-\\s]?\\d{4}$"`
	Email string `bson:"email" json:"email" validate:"required, regexp=^[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*@[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*\\.[a-zA-Z]{2,3}$"`
}

// SelectVariant represents a selected variant in the order.
type SelectVariant struct {
	ProductID        int64   `bson:"productId" json:"productId" validate:"required"`
	ProductName      string  `bson:"productName, omitempty" json:"productName,omitempty"`
	ProductVersion   int32   `bson:"productVersion" json:"productVersion" validate:"required"`
	VariantID        int64   `bson:"variantId" json:"variantId" validate:"required"`
	Date             string  `bson:"date" json:"date" validate:"required,datetime=2006-01-02,regexp=^(19|20)\\d{2}-(0[1-9]|1[012])-(0[1-9]|[ 12][0-9]|3[0-1])$"`
	Time             string  `bson:"time" json:"time" validate:"required,datetime=15:04,regexp=^([1-9]|[01][0-9]|2[0-3]):([0-5][0-9])$"` // Corrected tags
	Quantity         int32   `bson:"quantity" json:"quantity" validate:"required,gt=0"`
	Currency         string  `bson:"currency" json:"currency" validate:"required"` // 3-character currency code
	PartnerSalePrice float32 `bson:"partnerSalePrice" json:"partnerSalePrice" validate:"required,gt=0"`
	CostPrice        float32 `bson:"costPrice" json:"costPrice" validate:"required,gt=0"`
}

// ValidityPeriod represents the validity period of an order variant.
type ValidityPeriod struct {
	StartDateTime time.Time `bson:"startDateTime" json:"startDateTime" validate:"required"`
	EndDateTime   time.Time `bson:"endDateTime" json:"endDateTime" validate:"required,gtfield=StartDateTime"`
	Timezone      string    `bson:"timezone" json:"timezone" validate:"required"`
	Offset        string    `bson:"offset" json:"offset" validate:"required"`
}

// Voucher represents the voucher information.
type Voucher struct {
	VoucherProvideStatusCode string `bson:"voucherProvideStatusCode" json:"voucherProvideStatusCode" validate:"required"`
	VoucherCode              string `bson:"voucherCode" json:"voucherCode" validate:"required"`
	VoucherDisplayTypeCode   string `bson:"voucherDisplayTypeCode" json:"voucherDisplayTypeCode" validate:"required"`
	VoucherSync              bool   `bson:"voucherSync" json:"voucherSync" validate:"required" default:"false"`
	VoucherPdfUrl            string `bson:"voucherPdfUrl" json:"voucherPdfUrl" default:""`
	PdfUrlCreatedAt          string `bson:"pdfUrlCreatedAt" json:"pdfUrlCreatedAt"`
}

// OrderVariantItem represents an item within an order variant.
type OrderVariantItem struct {
	OrderVariantItemID   int64          `bson:"orderVariantItemId" json:"orderVariantItemId" validate:"required"`
	OrderVariantItemName string         `bson:"orderVariantItemName" json:"orderVariantItemName" validate:"required"`
	ValidityPeriod       ValidityPeriod `bson:"validityPeriod" json:"validityPeriod" validate:"required,dive"`
	Voucher              Voucher        `bson:"voucher" json:"voucher" validate:"required,dive"`
}

// OrderVariant represents a variant of an order.
type OrderVariant struct {
	OrderVariantID              int64                 `bson:"orderVariantId" json:"orderVariantId" validate:"required"`
	ProductID                   int64                 `bson:"productId" json:"productId" validate:"required"`
	ProductVersion              int32                 `bson:"productVersion" json:"productVersion" validate:"required"`
	VariantID                   int64                 `bson:"variantId" json:"variantId" validate:"required"`
	VariantName                 string                `bson:"variantName" json:"variantName" validate:"required"`
	Date                        string                `bson:"date" json:"date" validate:"required,datetime=2006-01-02"` // Validate as date in YYYY-MM-DD format
	Time                        string                `bson:"time" json:"time" validate:"required,datetime=15:04"`      // Validate as time in HH:MM format
	ValidityPeriod              ValidityPeriod        `bson:"validityPeriod" json:"validityPeriod" validate:"required,dive"`
	OrderVariantStatusTypeCode  string                `bson:"orderVariantStatusTypeCode" json:"orderVariantStatusTypeCode" validate:"required"`
	RefundApprovalTypeCode      string                `bson:"refundApprovalTypeCode" json:"refundApprovalTypeCode" validate:"required"`
	UsedRestoreDateTime         string                `bson:"usedRestoreDateTime,omitempty" json:"usedRestoreDateTime" default:""`
	UsedRestoreDateTimeTimezone string                `bson:"usedRestoreDateTimeTimezone,omitempty" json:"usedRestoreDateTimeTimezone" default:""`
	UsedRestoreDateTimeOffset   string                `bson:"usedRestoreDateTimeOffset,omitempty" json:"usedRestoreDateTimeOffset" default:""`
	CanceledDateTime            string                `bson:"canceledDateTime,omitempty" json:"canceledDateTime" default:""`
	CanceledDateTimeTimezone    string                `bson:"canceledDateTimeTimezone,omitempty" json:"canceledDateTimeTimezone" default:""`
	CanceledDateTimeOffset      string                `bson:"canceledDateTimeOffset,omitempty" json:"canceledDateTimeOffset" default:""`
	CancelStatusCode            string                `bson:"cancelStatusCode" json:"cancelStatusCode" validate:"required" default:""` // New field with default value
	RefundInfo                  string                `bson:"refundInfo" json:"refundInfo" validate:"required" default:""`
	OrderCancelTypeCode         string                `bson:"orderCancelTypeCode" json:"orderCancelTypeCode" validate:"required"`
	CancelRejectTypeCode        string                `bson:"cancelRejectTypeCode" json:"cancelRejectTypeCode" validate:"required" default:""`
	Message                     string                `bson:"message,omitempty" json:"message,omitempty"`                                      // check for use
	CancelFailReasonCode        string                `bson:"cancelFailReasonCode" json:"cancelFailReasonCode" validate:"required" default:""` // New field with default value
	OrderVariantItems           []OrderVariantItem    `bson:"orderVariantItems" json:"orderVariantItems" validate:"required,dive,min=1"`
	ReconciliationByDate        []ReconcilationDetail `bson:"reconciliationByDate" json:"reconciliationByDate"validate:"required"`
	ForceCancelTypeCode         string                `bson:"forceCancelTypeCode" json:"forceCancelTypeCode" validate:"required"`
	CancelVariantSync           bool                  `bson:"cancelVariantSync" json:"cancelVariantSync" validate:"required" default:"false""`
}

type ReconcilationDetail struct {
	ReconciliationDate       string `bson:"reconciliationDate" json:"reconciliationDate" validate:"required"`
	ReconcileOrderStatusCode string `bson:"reconcileOrderStatusCode"json:"reconcileOrderStatusCode" validate:"required"`
}

// OrderResponse represents the full order response.
type Model struct {
	Id                            string          `bson:"_id,omitempty" json:"id"`
	Suppliers                     string          `bson:"suppliers,omitempty" json:"suppliers,omitempty"`
	PartnerOrderID                string          `bson:"partnerOrderId" json:"partnerOrderId" validate:"required, min=0, max <=50"`
	PartnerOrderGroupID           string          `bson:"partnerOrderGroupId" json:"partnerOrderGroupId"`
	OrderId                       int64           `bson:"orderId" json:"orderId" validate:"required"`
	PartnerOrderChannelCode       string          `bson:"partnerOrderChannelCode" json:"partnerOrderChannelCode"`
	PartnerOrderChannelName       string          `bson:"partnerOrderChannelName" json:"partnerOrderChannelName" validate:"min=0, max=255"`
	TotalSelectedVariantsQuantity int32           `bson:"totalSelectedVariantsQuantity" json:"totalSelectedVariantsQuantity" validate:"gt=0"`
	OrderStatusCode               string          `bson:"orderStatusCode" json:"orderStatusCode" validate:"required"`
	Customer                      Customer        `bson:"customer" json:"customer" validate:"required,dive"`
	ActualCustomer                Customer        `bson:"actualCustomer" json:"actualCustomer" validate:"required,dive"`
	SelectVariants                []SelectVariant `bson:"selectVariants" json:"selectVariants" validate:"required,dive"`
	OrderVariants                 []OrderVariant  `bson:"orderVariants" json:"orderVariants"`
	OrderExpired                  bool            `bson:"orderExpired" json:"orderExpired" default:"false"`
	OodoSyncStatus                bool            `bson:"oodoSyncStatus" json:"oodoSyncStatus" default:"false"`
	CreatedAt                     string          `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt                     string          `bson:"updatedAt" json:"updatedAt"`
}
