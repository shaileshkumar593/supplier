package travolution

type ProductReq struct {
	ProductUid int    `json:"productUid,omitempty"`
	Take       int    `json:"take,omitempty"`
	Skip       int    `json:"skip,omitempty"`
	Lang       string `json:"lang,omitempty" validate:"oneof=ko en ja zh-CN zh-TW"`
}

type OptionRequest struct {
	ProductUid int         `json:"productUid" validate:"required"`
	OptionUid  interface{} `json:"optionUid,omitempty"`
	Lang       string      `json:"lang,omitempty" validate:"oneof=ko en ja zh-CN zh-TW"`
}

type UnitPriceRequest struct {
	ProductUid int         `json:"productUid" validate:"required"`
	OptionUid  interface{} `json:"optionUid" validate:"required"` // int or string (PKG, PAS)
	UnitUid    interface{} `json:"unitUid,omitempty"`             // optional: int or string (PKG, PAS)
}

// BookingScheduleReq
type BookingScheduleReq struct {
	ProductUid int         `json:"productUid" validate:"required"`
	OptionUid  interface{} `json:"optionUid" validate:"required"`
	Date       string      `json:"date,omitempty"` // format: YYYYMMDD
	Time       string      `json:"time,omitempty"` // format: HHMM
}

type BookingAdditionalInfoRequest struct {
	ProductUID        int         `json:"productUid" validate:"required"`
	OptionUID         interface{} `json:"optionUid" validate:"required"`
	AdditionalInfoUID interface{} `json:"additionalInfoUid,omitempty"`
}
