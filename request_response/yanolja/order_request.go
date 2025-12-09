package yanolja

// Blocking of order in advance
type WaitingForOrder struct {
	PartnerOrderID                string         `json:"partnerOrderId" binding:"required" validate:"required"`
	PartnerOrderGroupID           string         `json:"partnerOrderGroupId"`
	PartnerOrderChannelCode       string         `json:"partnerOrderChannelCode"`
	PartnerOrderChannelName       string         `json:"partnerOrderChannelName"`
	TotalSelectedVariantsQuantity int32          `json:"totalSelectedVariantsQuantity"`
	Customer                      CustomerDetail `json:"customer" binding:"required" validate:"required"`
	ActualCustomer                CustomerDetail `json:"actualCustomer" binding:"required" validate:"required"`
	SelectVariants                []VariantInfo  `json:"selectVariants" binding:"required" validate:"required"`
}

type CustomerDetail struct {
	Name  string `json:"name" binding:"required" validate:"required"`
	Tel   string `json:"tel" binding:"required" validate:"required"`
	Email string `json:"email"`
}

type VariantInfo struct {
	ProductID        int64   `json:"productId" binding:"required" validate:"required"`
	ProductVersion   int32   `json:"productVersion" binding:"required" validate:"required"`
	VariantID        int64   `json:"variantId" binding:"required" validate:"required"`
	Date             *string `json:"date"`
	Time             *string `json:"time"`
	Quantity         int32   `json:"quantity" binding:"required" validate:"required"`
	Currency         string  `json:"currency" binding:"required" validate:"required"`
	PartnerSalePrice float32 `json:"partnerSalePrice" binding:"required" validate:"required"`
	CostPrice        float32 `json:"costPrice" binding:"required" validate:"required"`
}

// pre order confirmation request
type OrderConfirmation struct {
	OrderId        int64  `json:"orderId"`
	PartnerOrderId string `json:"partnerOrderId"`
}

type PartnerOrder struct {
	PartnerOrderId string `json:"partnerOrderId"`
}

type Order struct {
	OrderId int64 `json:"orderId" binding:"required" validate:"required"`
}

type OrderList struct {
	OrderId []int64 `json:"orderId" binding:"required" validate:"required"`
}

//
type OrderReconcileReq struct {
	ReconciliationDate       string `json:"reconciliationDate" binding:"required" validate:"required,datetime=2006-01-02"`
	ReconcileOrderStatusCode string `json:"reconcileOrderStatusCode" binding:"required" validate:"required,oneof=UNKNOWN CREATED CANCELING CANCELED"`
	PageNumber               int    `json:"pageNumber" binding:"required" validate:"required,min=1"`
	PageSize                 int    `json:"pageSize" binding:"required" validate:"required,min=1,max=1000"`
}

type OrderRequestFromGGT struct {
	TypeOfRequest string `json:"typeOfRequest" validate:"required"`
}
