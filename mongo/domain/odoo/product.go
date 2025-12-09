package odoo

type Product struct {
	Id                  string                   `bson:"_id" json:"id,omitempty"`
	ProductID           int64                    `bson:"product_id" json:"productId"`
	ProductVersion      int32                    `bson:"productVersion" json:"productVersion"`
	ProductName         string                   `bson:"productName" json:"productName"`
	ProductStatusCode   string                   `bson:"productStatusCode" json:"productStatusCode"` // Product Status Code
	ProductTypeCode     string                   `bson:"productTypeCode" json:"productTypeCode"`     // Product Type Code
	ProductValidity     string                   `bson:"productValidity" json:"productValidity"`
	Option1Variant      string                   `bson:"option1Variant" json:"option1Variant"`
	FacilityAddress     []FacilityLocationDetail `bson:"facilityAddress" json:"facilityAddress"`
	Regions             []Region                 `bson:"regions" json:"regions"`
	IsIntegratedVoucher bool                     `bson:"isIntegratedVoucher" json:"isIntegratedVoucher"`
	Variants            []ProductVariant         `json:"variants" json:"variants"`
	CreatedAt           string                   `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt           string                   `bson:"updatedAt" json:"updatedAt"`
}

// variant
type ProductVariant struct {
	VariantID                           int64    `bson:"variantId,omitempty" json:"variantId,omitempty"`
	ProductID                           int64    `bson:"productId,omitempty" json:"productId,omitempty"`
	VariantName                         string   `bson:"variantName,omitempty" json:"variantName,omitempty"`
	VariantDescription                  string   `bson:"variantDescription,omitempty" json:"variantDescription,omitempty"`
	RefundApprovalTypeCode              string   `bson:"refundApprovalTypeCode,omitempty" json:"refundApprovalTypeCode"` // validate:"oneof=DIRECT ADMIN"`
	IsRefundableAfterExpiration         bool     `bson:"isRefundableAfterExpiration,omitempty" json:"isRefundableAfterExpiration,omitempty"`
	RefundInfo                          string   `bson:"refundInfo,omitempty" json:"refundInfo,omitempty"`
	VariantStatusCode                   string   `bson:"variantStatusCode" json:"variantStatusCode"`
	OrderExpirationUsageProcessTypeCode string   `bson:"orderExpirationUsageProcessTypeCode,omitempty" json:"orderExpirationUsageProcessTypeCode,omitempty"` // validate:"oneof=NONE FORCED_USE"`
	OrderExpirationDateTypeCode         string   `bson:"orderExpirationDateTypeCode,omitempty" json:"orderExpirationDateTypeCode,omitempty"`
	ValidityStartDate                   string   `bson:"validityStartDate" json:"validityStartDate"`
	ValidityEndDate                     string   `bson:"validityEndDate" json:"validityEndDate"`
	ChannelNetPrice                     float64  `bson:"channelNetPrice" json:"channelNetPrice"`
	RetailPrice                         float64  `bson:"retailPrice" json:"retailPrice"`
	DiscountSalePrice                   float64  `bson:"discountSalePrice" json:"discountSalePrice"`
	Currency                            string   `bson:"currency" json:"currency"`
	SalePrice                           float64  `bson:"salePrice" json:"salePrice"`
	SupplierCostPrice                   float64  `bson:"supplierCostPrice" json:"supplierCostPrice"`
	IsRound                             bool     `bson:"isRound" json:"isRound"`
	IsSchedule                          bool     `bson:"isSchedule" json:"isSchedule"`
	VoucherDisplayCode                  []string `bson:"voucherDisplayCode" json:"voucherDisplayCode"`
}

// facility Address
type FacilityLocationDetail struct {
	Latitude  float64 `json:"latitude" json:"latitude"`
	Longitude float64 `json:"longitude" json:"longitude"`
	Address   string  `json:"address" json:"address"`
}

type Region struct {
	AreaCode  string `json:"areaCode" json:"areaCode"`
	Area      string `json:"area" json:"area"`
	City      string `json:"city" json:"city"`
	Country   string `json:"country" json:"country"`
	Continent string `json:"continent" json:"continent"`
}
