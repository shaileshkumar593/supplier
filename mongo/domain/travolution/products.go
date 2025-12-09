package travolution

type Product struct {
	Id                       string            `bson:"_id,omitempty" json:"id,omitempty"`
	SupplierName             string            `bson:"supplierName" json:"supplierName"`
	ProductUID               int               `bson:"productUid" json:"productUid" validate:"required"`
	Type                     string            `bson:"type" json:"type"`
	Status                   int               `bson:"status" json:"status"`
	SaleTarget               string            `bson:"saleTarget" json:"saleTarget"`
	HasBookingAdditionalInfo string            `bson:"hasBookingAdditionalInfo" json:"hasBookingAdditionalInfo"`
	VoucherType              int               `bson:"voucherType" json:"voucherType"`
	Titles                   map[string]string `bson:"titles" json:"titles"`
	Images                   ProductImages     `bson:"images" json:"images"`
	Contents                 ProductContent    `bson:"contents,omitempty" json:"contents,omitempty"`
	Options                  []Option          `bson:"options,omitempty" json:"options,omitempty"`
	CreatedAt                string            `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt                string            `bson:"updatedAt" json:"updatedAt"`
}

type ProductImages struct {
	Main string   `bson:"main" json:"main"`
	List []string `bson:"list" json:"list"`
}

type ProductContent struct {
	Highlights  string              `bson:"highlights,omitempty" json:"highlights,omitempty"`
	Description string              `bson:"description,omitempty" json:"description,omitempty"`
	SubContent  []ProductSubContent `bson:"subContent,omitempty" json:"subContent,omitempty"`
}

// ProductSubContent defines each subsection of product content.
// define operations hour, how to use, notice, refund policy
type ProductSubContent struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}

type Option struct {
	OptionUID               string                  `bson:"optionUid" json:"optionUid" validate:"required"`
	Names                   map[string]string       `bson:"names" json:"names"`
	Notice                  string                  `bson:"notice,omitempty" json:"notice,omitempty"`
	UnitAndPriceDetails     []UnitAndPriceDetail    `bson:"unitAndPriceDetails,omitempty" json:"unitAndPriceDetails,omitempty"`
	BookingSchedules        []BookingSchedule       `bson:"bookingSchedules,omitempty" json:"bookingSchedules,omitempty"`
	AdditionalBookingDetail []BookingAdditionalInfo `bson:"additionalBookingDetail,omitempty" json:"additionalBookingDetail,omitempty"`
}
type UnitAndPriceDetail struct {
	UnitUID       string            `bson:"UnitUid" json:"UnitUid"`
	Currency      string            `bson:"currency" json:"currency"`
	OriginalPrice float32           `bson:"originalPrice" json:"originalPrice"`
	B2BPrice      float32           `bson:"B2Bprice" json:"B2Bprice"`
	B2CPrice      float32           `bson:"B2Cprice" json:"B2Cprice"`
	MinAmount     float32           `bson:"minAmount,omitempty" json:"minAmount,omitempty"`
	MaxAmount     float32           `bson:"maxAmount,omitempty" json:"maxAmount,omitempty"`
	Names         map[string]string `bson:"names" json:"names"`
}

type BookingSchedule struct {
	Date  string `bson:"date" json:"date"`
	Time  string `bson:"time" json:"time"`
	Stock int    `bson:"stock" json:"stock"`
}

type BookingAdditionalInfo struct {
	AdditionalInfoUID string                    `bson:"additionalInfoUid" json:"additionalInfoUid"`
	Type              string                    `bson:"type" json:"type"`
	AnswerType        string                    `bson:"answerType" json:"answerType"`
	Titles            map[string]string         `bson:"titles" json:"titles"`
	Options           []BookingAdditionalOption `bson:"options" json:"options"`
}

type BookingAdditionalOption struct {
	Titles map[string]string `bson:"titles" json:"titles"`
	Value  string            `bson:"value" json:"value"`
}
