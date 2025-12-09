package trip

type PreOrderPaymentRequest struct {
	SequenceId           string                `json:"sequenceId" binding:"required"`
	OtaOrderId           string                `json:"otaOrderId" binding:"required"`
	SupplierOrderId      string                `json:"supplierOrderId" binding:"required"`
	ConfirmType          int                   `json:"confirmType" binding:"required"`
	OrderLastConfirmTime string                `json:"orderLastConfirmTime" binding:"required"`
	Items                []PreOrderPaymentItem `json:"items" binding:"required"`
	Coupons              []PreOrderCoupon      `json:"coupons"`
}

type PreOrderPaymentItem struct {
	ItemId string `json:"itemId" binding:"required"`
	PLU    string `json:"PLU" binding:"required"`
}

type PreOrderCoupon struct {
	Type           int     `json:"type" binding:"required"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Amount         float64 `json:"amount"`
	AmountCurrency string  `json:"amountCurrency"`
}
