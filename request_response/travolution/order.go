package travolution

import (
	domain "swallow-supplier/mongo/domain/travolution"
)

type VoucherSendType int

const (
	DoNotSend VoucherSendType = iota
	SendEmail
	SendKakao
	SendEmailAndKakao
)

// OrderRequest represents the incoming request payload for creating an order.
type OrderRequest struct {
	Product               int                 `json:"product" validate:"required"`                                    // Product Unique ID
	Option                interface{}         `json:"option" validate:"required"`                                     // Option Unique ID
	UnitAmounts           []domain.UnitAmount `json:"unitAmounts" validate:"required,dive"`                           // Pricing Unit Unique ID and Order Quantity
	BookingDate           string              `json:"bookingDate,omitempty" validate:"omitempty,datetime=2006-01-02"` // Reservation Date (YYYY-MM-DD)
	BookingTime           string              `json:"bookingTime,omitempty" validate:"omitempty"`                     // Reservation Time (HH:mm)
	BookingAdditionalInfo interface{}         `json:"bookingAdditionalInfo,omitempty" validate:"omitempty"`           // BAC = object, TRV = object[][]
	ReferenceNumber       string              `json:"referenceNumber" validate:"omitempty,max=64"`                    // Partner’s Own Order Number
	VoucherSendType       VoucherSendType     `json:"voucherSendType" validate:"omitempty,oneof=0 1 2 3"`             // Voucher send type (enum) value = 0 everytime
	TravelerName          string              `json:"travelerName" validate:"required,max=64"`                        // Traveler’s Name
	TravelerContactEmail  string              `json:"travelerContactEmail" validate:"required,email,max=64"`          // Traveler’s Email address
	TravelerContactNumber string              `json:"travelerContactNumber" validate:"required,max=32"`               // Traveler’s Contact number
	TravelerNationality   string              `json:"travelerNationality" validate:"required,len=2,uppercase"`        // ISO 3166-1 alpha-2
}

// UnitAmount holds pricing unit ID and quantity.
/* type UnitAmount struct {
	Unit   int `json:"unit" validate:"required"`   // Pricing Unit Unique ID
	Amount int `json:"amount" validate:"required"` // Order quantity
} */

// BookingAdditionalItem represents one element of bookingAdditionalInfo
// when the type is "BAC" (flat object).
type BookingAdditionalItem struct {
	AdditionalInfo string `json:"additionalInfo" validate:"required,uuid4"` // UUID of the info type
	Value          string `json:"value" validate:"required"`                // Provided value
}

// WebhookPayload  for taking REDEEMED, RESTORED, CANCELED, BOOKING_ACCEPTED, BOOKING_REJECTED

type TicketRequest struct {
	Product               int                 `json:"product" validate:"required"`                             // Product Unique ID
	Option                int                 `json:"option" validate:"required"`                              // Option Unique ID
	UnitAmounts           []domain.UnitAmount `json:"unitAmounts" validate:"required,dive"`                    // List of pricing units with quantity
	ReferenceNumber       string              `json:"referenceNumber" validate:"required,max=64"`              // Partner’s own order number
	VoucherSendType       int                 `json:"voucherSendType" validate:"oneof=0 1 2 3"`                // Voucher send type enum
	TravelerName          string              `json:"travelerName" validate:"required,max=64"`                 // Traveler’s name
	TravelerContactEmail  string              `json:"travelerContactEmail" validate:"required,email,max=64"`   // Traveler’s email
	TravelerContactNumber string              `json:"travelerContactNumber" validate:"required,max=32"`        // Traveler’s contact number
	TravelerNationality   string              `json:"travelerNationality" validate:"required,len=2,uppercase"` // ISO 3166-1 alpha-2
}

type PASSOrPKGRequest struct {
	Product               int                 `json:"product" bson:"product" validate:"required"`                                         // Product Unique ID
	Option                string              `json:"option" bson:"option" validate:"required,oneof=PAS PKG"`                             // Option Unique ID (enum)
	UnitAmounts           []domain.UnitAmount `json:"unitAmounts" bson:"unitAmounts" validate:"required,dive"`                            // List of pricing units with quantity
	ReferenceNumber       string              `json:"referenceNumber" bson:"referenceNumber" validate:"required,max=64"`                  // Partner’s own order number
	VoucherSendType       int                 `json:"voucherSendType" bson:"voucherSendType" validate:"oneof=0 1 2 3"`                    // Voucher send type enum
	TravelerName          string              `json:"travelerName" bson:"travelerName" validate:"required,max=64"`                        // Traveler’s name
	TravelerContactEmail  string              `json:"travelerContactEmail" bson:"travelerContactEmail" validate:"email,max=64"`           // Traveler’s email
	TravelerContactNumber string              `json:"travelerContactNumber" bson:"travelerContactNumber" validate:"required,max=32"`      // Traveler’s contact number
	TravelerNationality   string              `json:"travelerNationality" bson:"travelerNationality" validate:"required,len=2,uppercase"` // ISO 3166-1 alpha-2
}

type BookingRequest struct {
	Product               int                 `json:"product" validate:"required"`                                    // Product Unique ID
	Option                int                 `json:"option" validate:"required"`                                     // Option Unique ID
	UnitAmounts           []domain.UnitAmount `json:"unitAmounts" validate:"required,dive"`                           // Pricing Unit Unique ID and Order Quantity
	BookingDate           string              `json:"bookingDate,omitempty" validate:"omitempty,datetime=2006-01-02"` // Reservation Date (YYYY-MM-DD)
	BookingTime           string              `json:"bookingTime,omitempty" validate:"omitempty"`                     // Reservation Time (HH:mm)
	BookingAdditionalInfo interface{}         `json:"bookingAdditionalInfo,omitempty" validate:"omitempty"`           // BAC = object, TRV = object[][]
	ReferenceNumber       string              `json:"referenceNumber" validate:"omitempty,max=64"`                    // Partner’s Own Order Number
	VoucherSendType       VoucherSendType     `json:"voucherSendType" validate:"omitempty,oneof=0 1 2 3"`             // Voucher send type (enum) value = 0 everytime
	TravelerName          string              `json:"travelerName" validate:"required,max=64"`                        // Traveler’s Name
	TravelerContactEmail  string              `json:"travelerContactEmail" validate:"required,email,max=64"`          // Traveler’s Email address
	TravelerContactNumber string              `json:"travelerContactNumber" validate:"required,max=32"`               // Traveler’s Contact number
	TravelerNationality   string              `json:"travelerNationality" validate:"required,len=2,uppercase"`        // ISO 3166-1 alpha-2
}
