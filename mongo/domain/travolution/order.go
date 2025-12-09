package travolution

// Order represents an order stored in MongoDB.
type Order struct {
	ID              string        `bson:"_id,omitempty" json:"id"`
	OrderNumber     string        `bson:"orderNumber" json:"orderNumber" validate:"required"`
	ReferenceNumber string        `bson:"referenceNumber" json:"referenceNumber" validate:"required,max=64"`
	VoucherType     int           `bson:"voucherType" json:"voucherType" validate:"required,oneof=1 2 3"` // VoucherSendType
	Status          string        `bson:"status" json:"status" validate:"required"`                       // AV, AP, CR, CL, PC, EP
	EventHistory    []WebHookData `bson:"eventHistory" json:"eventHistory" validate:"required"`
	Type            string        `bson:"type" json:"type" validate:"required"` // TK, BK, PAS, PKG
	Product         int           `bson:"product" json:"product" validate:"required"`
	Option          interface{}   `bson:"option" json:"option" validate:"required"`
	UnitAmounts     []UnitAmount  `bson:"unitAmounts" json:"unitAmounts" validate:"required,dive"`
	VoucherInfo     []VoucherInfo `bson:"voucherInfo" json:"voucherInfo"`

	// Booking-related
	BookingDate           string      `bson:"bookingDate,omitempty" json:"bookingDate,omitempty" validate:"omitempty,datetime=2006-01-02"` // Reservation Date
	BookingTime           string      `bson:"bookingTime,omitempty" json:"bookingTime,omitempty" validate:"omitempty"`                     // Reservation Time
	BookingAdditionalInfo interface{} `bson:"bookingAdditionalInfo,omitempty" json:"bookingAdditionalInfo,omitempty" validate:"omitempty"` // BAC = object, TRV = object[][]

	// Traveler details
	TravelerName          string `bson:"travelerName" json:"travelerName" validate:"required,max=64"`
	TravelerContactEmail  string `bson:"travelerContactEmail" json:"travelerContactEmail" validate:"required,email,max=64"`
	TravelerContactNumber string `bson:"travelerContactNumber" json:"travelerContactNumber" validate:"required,max=32"`
	TravelerNationality   string `bson:"travelerNationality" json:"travelerNationality" validate:"required,len=2,uppercase"`

	// System-generated booking status
	BookingStatus     string `bson:"bookingStatus,omitempty" json:"bookingStatus,omitempty"`
	BookingAt         string `bson:"bookingAt,omitempty" json:"bookingAt,omitempty"`
	ExpiredAt         string `bson:"expiredAt,omitempty" json:"expiredAt,omitempty"`
	ApprovedAt        string `bson:"approvedAt,omitempty" json:"approvedAt,omitempty"`
	CancelRequestedAt string `bson:"cancelRequestedAt,omitempty" json:"cancelRequestedAt,omitempty"`
	CanceledAt        string `bson:"canceledAt,omitempty" json:"canceledAt,omitempty"`

	// Audit fields
	CreatedAt string `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt string `bson:"updatedAt" json:"updatedAt" validate:"required"`
}

// UnitAmount represents the pricing unit and quantity.
type UnitAmount struct {
	Unit   interface{} `bson:"unit" json:"unit" validate:"required"`
	Amount int         `bson:"amount" json:"amount" validate:"gte=1"`
}

// VoucherInfo represents voucher details associated with an order.
type VoucherInfo struct {
	Product       int         `bson:"product" json:"product"`
	Option        interface{} `bson:"option" json:"option"`
	BookingStatus string      `bson:"bookingStatus,omitempty" json:"bookingStatus,omitempty"`
	Unit          interface{} `bson:"unit" json:"unit"`
	Amount        int         `bson:"amount" json:"amount"`
	CodeType      string      `bson:"codeType" json:"codeType"`
	VoucherCode   string      `bson:"voucherCode" json:"voucherCode"`
	VoucherFile   string      `bson:"voucherFile" json:"voucherFile"`
	VoucherLink   string      `bson:"voucherLink" json:"voucherLink"`
}

// BookingAdditionalItem represents one element when bookingAdditionalInfo is of type BAC (flat object).
type BookingAdditionalItem struct {
	AdditionalInfo string `json:"additionalInfo" validate:"required,uuid4"` // UUID of the info type
	Value          string `json:"value" validate:"required"`
}

type WebHookData struct {
	EventType string `bson:"eventType" json:"eventType" validate:"required"`
	DateAt    string `bson:"dateAt" json:"dateAt" validate:"required"`
}
