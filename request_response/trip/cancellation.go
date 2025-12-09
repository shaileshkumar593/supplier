package trip

type CancellationRequest struct {
	SequenceID      string             `json:"sequenceId" binding:"required"`
	OTAOrderID      string             `json:"otaOrderId" binding:"required"`
	SupplierOrderID string             `json:"supplierOrderId" binding:"required"`
	ConfirmType     int                `json:"confirmType" binding:"required"`
	Items           []CancellationItem `json:"items" binding:"required"`
}

// CancellationItem defines the structure for individual items in the preOrderCancel request.
type CancellationItem struct {
	ItemID          string            `json:"itemId" binding:"required"`
	PLU             string            `json:"PLU" binding:"required"`
	LastConfirmTime string            `json:"lastConfirmTime,omitempty"`
	CancelType      int               `json:"cancelType" binding:"required"`
	Quantity        int               `json:"quantity" binding:"required"`
	Passengers      []PassengerDetail `json:"passengers,omitempty"`
	Amount          float64           `json:"amount,omitempty"`
	AmountCurrency  string            `json:"amountCurrency ,omitempty"`
}

// PassengerDetail defines the structure for passenger-related information in the preOrderCancel request.
type PassengerDetail struct {
	PassengerID string `json:"passengerId" binding:"required"`
}
