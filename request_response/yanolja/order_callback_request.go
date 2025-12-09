package yanolja

// All callback

//used to send acknowlogment of cancelled order(full)
type CancellationAck struct {
	OrderId             int64  `json:"orderId" binding:"required" validate:"required"`
	PartnerOrderId      string `json:"partnerOrderId" binding:"required" validate:"required"`
	OrderCancelTypeCode string `json:"orderCancelTypeCode" binding:"required" validate:"required"`
}

//Order is not cancelled
type RefusalToCancel struct {
	OrderId              int64  `json:"orderId" binding:"required" validate:"required"`
	PartnerOrderId       string `json:"partnerOrderId" binding:"required" validate:"required"`
	OrderVariantID       int64  `json:"orderVariantId" binding:"required" validate:"required"`
	CancelRejectTypeCode string `json:"cancelRejectTypeCode" validate:"required"`
	Message              string `json:"message" validate:"required"`
}

// Check the order Variant status
type OrderStatusLookup struct {
	OrderId        int64  `json:"orderId" binding:"required" validate:"required"`
	PartnerOrderId string `json:"partnerOrderId" binding:"required" validate:"required"`
	OrderVariantID int64  `json:"orderVariantId" binding:"required"`
}

// ForcedOrderCancellation is used to set the cancellation reason of ordervariant
type ForcedOrderCancellation struct {
	OrderId             int64               `json:"orderId" binding:"required" validate:"required"`
	PartnerOrderId      string              `json:"partnerOrderId" binding:"required" validate:"required"`
	ForceCancelVariants []CancelledVariants `json:"forceCancelVariants" binding:"required" validate:"required"`
}

type CancelledVariants struct {
	OrderVariantID      int64  `json:"orderVariantId" binding:"required" validate:"required"`
	ForceCancelTypeCode string `json:"forceCancelTypeCode" binding:"required" validate:"required" `
}

// individual order varient Item voucher update
type IndividualVoucherUpdate struct {
	OrderId                int64  `json:"orderId" binding:"required" validate:"required"`
	PartnerOrderId         string `json:"partnerOrderId" binding:"required" validate:"required"`
	OrderVariantID         int64  `json:"orderVariantId" binding:"required" validate:"required"`
	OrderVariantItemId     int64  `json:"orderVariantItemId" binding:"required" validate:"required"`
	VoucherDisplayTypeCode string `json:"voucherDisplayTypeCode" binding:"required" validate:"required"`
	VoucherCode            string `json:"voucherCode" binding:"required" validate:"required"`
}

//combined order variant item update
type CombinedVoucherUpdate struct {
	OrderId                int64                    `json:"orderId" binding:"required" validate:"required"`
	PartnerOrderId         string                   `json:"partnerOrderId" binding:"required" validate:"required"`
	ProductID              int64                    `json:"productId" binding:"required" validate:"required"`
	OrderVariantIds        []OrderVariantAndItemIds `json:"orderVariantIds" binding:"required" validate:"required"`
	VoucherDisplayTypeCode string                   `json:"voucherDisplayTypeCode" binding:"required" validate:"required"`
	VoucherCode            string                   `json:"voucherCode" binding:"required" validate:"required"`
}

type OrderVariantAndItemIds struct {
	OrderVariantID      int64   `json:"orderVariantId" binding:"required" validate:"required"`
	OrderVariantItemIds []int64 `json:"orderVariantItemIds" binding:"required" validate:"required"`
}

//
type ProcessingOrRestoringReq struct {
	OrderId            int64  `json:"orderId" binding:"required" validate:"required"`
	EventType          string `json:"eventType" binding:"required" validate:"required"`
	PartnerOrderId     string `json:"partnerOrderId" binding:"required" validate:"required"`
	OrderVariantId     int64  `json:"orderVariantId" binding:"required" validate:"required"`
	OrderVariantItemId int64  `json:"orderVariantItemId" binding:"required" validate:"required"`
	DateTime           string `json:"dateTime" validate:"required,datetime=2006-01-02T15:04:05Z07:00"` // The date of use, restoration, or cancellation in UTC format
	DateTimeTimezone   string `json:"dateTimeTimeZone"  validate:"required"`                           // UTC timezone
	DateTimeOffset     string `json:"dateTimeOffset"  validate:"required"`                             // UTC offset
}
