package trip

type OrderInquiry struct { // for sandbox only
	SequenceID      string `json:"sequenceId" binding:"required"` // Trip Batch Number in the format yyyyMMdd + 32-bit GUID
	OTAOrderID      string `json:"otaOrderId" binding:"required"` // Trip order number
	SupplierOrderID string `json:"supplierOrderId,omitempty"`     // Supplier order number, optional
}

type OrderInquiryResponse struct {
	OTAOrderID      string          `json:"otaOrderId" binding:"required"`      // Trip order number
	SupplierOrderID string          `json:"supplierOrderId" binding:"required"` // Supplier order number
	Items           []OrderLineItem `json:"items" binding:"required,dive"`      // Array of order line item nodes
}

// OrderLineItem represents an individual order line item.
type OrderLineItem struct {
	ItemID         string                  `json:"itemId" binding:"required"`                                      // Order line item number, must be 0 if pre-order created but not paid
	UseStartDate   string                  `json:"useStartDate,omitempty" binding:"omitempty,datetime=2006-01-02"` // Format: yyyy-MM-dd
	UseEndDate     string                  `json:"useEndDate,omitempty" binding:"omitempty,datetime=2006-01-02"`   // Format: yyyy-MM-dd
	OrderStatus    int                     `json:"orderStatus" binding:"required,oneof=0 1 2"`                     // Strictly defined values
	Quantity       int                     `json:"quantity" binding:"required,gt=0"`                               // Must be greater than 0
	UseQuantity    int                     `json:"useQuantity" binding:"required,min=0"`                           // Must be non-negative
	CancelQuantity int                     `json:"cancelQuantity" binding:"required,min=0"`                        // Must be non-negative
	Passengers     []CancellationPassenger `json:"passengers,omitempty" binding:"omitempty,dive"`                  // Required when cancelType=2
	Vouchers       []Voucher               `json:"vouchers,omitempty" binding:"omitempty,dive"`                    // Array of voucher nodes
}

// CancellationPassenger represents a cancellation passenger node.
type CancellationPassenger struct {
	PassengerID     string `json:"passengerId" binding:"required"`                 // Passenger number
	PassengerStatus int    `json:"passengerStatus" binding:"required,oneof=0 1 2"` // 0: Pending; 1: Used; 2: Canceled
}

// Voucher represents a voucher node.
type Voucher struct {
	VoucherID     string `json:"voucherId" binding:"required"`                 // Voucher number
	VoucherStatus int    `json:"voucherStatus" binding:"required,oneof=0 1 2"` // 0: Pending; 1: Used; 2: Canceled
}
