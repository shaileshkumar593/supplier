package travolution

type VoucherInfoResponse struct {
	Product       int    `json:"product"`
	Option        string `json:"option"`
	BookingStatus string `json:"bookingStatus,omitempty"`
	Unit          int    `json:"unit"`
	Amount        int    `json:"amount"`
	CodeType      string `json:"codeType"`
	VoucherCode   string `json:"voucherCode"`
	VoucherFile   string `json:"voucherFile"`
	VoucherLink   string `json:"voucherLink"`
}

// TicketResponse represents the response for a ticket order
type TicketResponse struct {
	OrderNumber     string                `json:"orderNumber"`
	ReferenceNumber int                   `json:"referenceNumber"`
	VoucherType     int                   `json:"voucherType"`
	VoucherInfo     []VoucherInfoResponse `json:"voucherInfo"`
	ExpiredAt       string                `json:"expiredAt"`
}

type PassPkgResponse struct {
	OrderNumber     string                `json:"orderNumber"`
	ReferenceNumber string                `json:"referenceNumber"`
	VoucherType     int                   `json:"voucherType"`
	VoucherInfo     []VoucherInfoResponse `json:"voucherInfo"`
	ExpiredAt       string                `json:"expiredAt"`
}

type BookingResponse struct {
	OrderNumber     string                `json:"orderNumber"`
	ReferenceNumber string                `json:"referenceNumber"` // nullable
	VoucherType     int                   `json:"voucherType"`
	VoucherInfo     []VoucherInfoResponse `json:"voucherInfo"`
	ExpiredAt       string                `json:"expiredAt"`
}
