package yanolja

/*type OrderDetailResp struct {
	OrderID                       int            `json:"orderId"`
	PartnerOrderID                string         `json:"partnerOrderId"`
	PartnerOrderGroupID           string         `json:"partnerOrderGroupId"`
	TotalSelectedVariantsQuantity int            `json:"totalSelectedVariantsQuantity"`
	OrderStatusCode               string         `json:"orderStatusCode"`
	Customer                      CustomerDetail `json:"customer"`
	ActualCustomer                CustomerDetail `json:"actualCustomer"`
	OrderVariants                 []OrderVariant `json:"orderVariants"`
}

type OrderVariant struct {
	OrderVariantID             int                `json:"orderVariantId"`
	ProductID                  int                `json:"productId"`
	ProductVersion             int                `json:"productVersion"`
	VariantID                  int                `json:"variantId"`
	VariantName                string             `json:"variantName"`
	Date                       string             `json:"date"`
	Time                       string             `json:"time"`
	ValidityPeriod             ValidityPeriod     `json:"validityPeriod"`
	OrderVariantStatusTypeCode string             `json:"orderVariantStatusTypeCode"`
	RefundApprovalTypeCode     string             `json:"refundApprovalTypeCode"`
	UsedDateTime               string             `json:"usedDateTime,omitempty"`
	UsedDateTimeTimezone       string             `json:"usedDateTimeTimezone,omitempty"`
	UsedDateTimeOffset         string             `json:"usedDateTimeOffset,omitempty"`
	CanceledDateTime           string             `json:"canceledDateTime,omitempty"`
	CanceledDateTimeTimezone   string             `json:"CanceledDateTimeTimezone,omitempty"`
	CanceledDateTimeOffset     string             `json:"CanceledDateTimeOffset,omitempty"`
	OrderVariantItems          []OrderVariantItem `json:"orderVariantItems"`
}

type ValidityPeriod struct {
	StartDateTime string `json:"startDateTime"`
	EndDateTime   string `json:"endDateTime"`
	Timezone      string `json:"timezone"`
	Offset        string `json:"offset"`
}
type OrderVariantItem struct {
	OrderVariantItemID   int            `json:"orderVariantItemId"`
	OrderVariantItemName string         `json:"orderVariantItemName"`
	ValidityPeriod       ValidityPeriod `json:"validityPeriod"`
	Voucher              Voucher        `json:"voucher"`
}

type Voucher struct {
	VoucherProvideStatusCode string `json:"voucherProvideStatusCode"`
	VoucherCode              string `json:"voucherCode"`
	VoucherDisplayTypeCode   string `json:"voucherDisplayTypeCode"`
}*/

type CancelOrder struct {
	OrderId              int64  `json:"orderId"`
	PartnerOrderId       string `json:"partnerOrderId"`
	CancelStatusCode     string `json:"cancelStatusCode"`
	CancelFailReasonCode string `json:"cancelFailReasonCode"`
}

type OrderStatusLookupResp struct {
	OrderId        int64           `json:"orderId" validate:"required"`
	PartnerOrderId string          `json:"partnerOrderId" validate:"required"`
	OrderVariants  []OrderVariants `json:"orderVariants" validate:"required"`
}

type OrderVariants struct {
	OrderVariantID             int64  `json:"orderVariantId" validate:"required"`
	OrderVariantStatusTypeCode string `json:"orderVariantStatusTypeCode" validate:"required"`
}
type OrderReconcilationResp struct {
	Orders []OrderReconcilation `json:"orders" validate:"required"`
}

type OrderReconcilation struct {
	ReconciliationDate         string  `json:"reconciliationDate" validate:"required"`
	ProductId                  int64   `json:"productId" validate:"required"`
	VariantId                  int64   `json:"variantId" validate:"required"`
	OrderVariantId             string  `json:"orderVariantId" validate:"required"`
	PartnerSalePrice           float32 `json:"partnerSalePrice" validate:"required"`
	ReconcileOrderStatusCode   string  `json:"reconcileOrderStatusCode" validate:"required"`
	PartnerOrderId             string  `json:"partnerOrderId" validate:"required"`
	Pin                        string  `json:"pin" validate:"required"`
	PartnerOrderChannelPin     string  `json:"partnerOrderChannelPin" validate:"required"`
	PartnerOrderChannelName    string  `json:"partnerOrderChannelName"`
	PartnerOrderChannelCode    string  `json:"partnerOrderChannelCode"`
	PartnerOrderChannelOrderId string  `json:"partnerOrderChannelOrderId"`
}

/*type PreOrderPrecaution struct {
	ProductId                   int64       `json:"productId" validate:"required"`
	VariantId                   int64       `json:"variantId" validate:"required"`
	VariantStatusCode           string      `json:"variantStatusCode"`
	ProductStatusCode           string      `json:"productStatusCode"`
	ProductSalePeriod           SalePeriodS `json:"productSalePeriod"`
	OrderSalePeriod             SalePeriodS `json:"orderSalePeriod"`
	IsCancelPenalty             bool        `json:"isCancelPenalty"`
	IsRefundableAfterExpiration bool        `json:"isRefundableAfterExpiration"`
	IsAvailableOnPurchaseDate   bool        `json:"isAvailableOnPurchaseDate"`
	RefundInfo                  string      `json:"refundInfo"`
}*/
