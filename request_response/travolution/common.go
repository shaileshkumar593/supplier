package travolution

import (
	domain "swallow-supplier/mongo/domain/travolution"
)

type RawProduct struct {
	UID                      int                   `json:"uid"`
	Type                     string                `json:"type"`
	Status                   int                   `json:"status"`
	SaleTarget               string                `json:"saleTarget"`
	HasBookingAdditionalInfo string                `json:"hasBookingAdditionalInfo"`
	VoucherType              int                   `json:"voucherType"`
	Titles                   map[string]string     `json:"titles"`
	Images                   domain.ProductImages  `json:"images"`
	Contents                 domain.ProductContent `bson:"contents,omitempty" json:"contents,omitempty"`
	Options                  []RawProductOption    `bson:"options,omitempty" json:"options,omitempty"`
}

type RawProductOption struct {
	UID                   string                     `json:"uid"`
	Names                 map[string]string          `json:"names"`
	Notice                string                     `json:"notice,omitempty"`
	UnitsPrice            []RawOptionUnitPrice       `json:"unitPrice,omitempty"`
	BookingSchedules      []domain.BookingSchedule   `json:"bookingSchedules,omitempty"`
	AdditionalBookingInfo []RawBookingAdditionalInfo `json:"additionalBookingInfo,omitempty"`
}

type RawOptionUnitPrice struct {
	UID           string            `json:"uid"`
	Currency      string            `json:"currency"`
	OriginalPrice int               `json:"originalPrice"`
	B2BPrice      int               `json:"B2Bprice"`
	B2CPrice      int               `json:"B2Cprice"`
	MinAmount     int               `json:"minAmount,omitempty"`
	MaxAmount     int               `json:"maxAmount,omitempty"`
	Names         map[string]string `json:"names"`
}

type RawBookingAdditionalInfo struct {
	UID        string                           `json:"uid"`
	Type       string                           `json:"type"`
	AnswerType string                           `json:"answerType"`
	Titles     map[string]string                `json:"titles"`
	Options    []domain.BookingAdditionalOption `json:"options"`
}
