package yanolja

type ProductView struct {
	Id                          string            `bson:"_id,omitempty" json:"id,omitempty"`              // MongoDB document ID
	ProductID                   int64             `bson:"productId" json:"productId" validate:"required"` // Product ID
	ProductName                 string            `bson:"productName" json:"productName" validate:"required"`
	SupplierName                string            `bson:"supplierName" json:"supplierName" validate:"required"`                                                               // Product Name
	ProductVersion              int32             `bson:"productVersion" json:"productVersion" validate:"required"`                                                           // Product Version
	ProductInfo                 ProductInfo       `bson:"productInfo" json:"productInfo" validate:"required"`                                                                 // Product Information
	Price                       Price             `bson:"price" json:"price" validate:"required"`                                                                             // Product Price
	ProductStatusCode           string            `bson:"productStatusCode" json:"productStatusCode" validate:"required,oneof=WAITING_FOR_SALE IN_SALE SOLD_OUT END_OF_SALE"` // Product Status Code
	ProductTypeCode             string            `bson:"productTypeCode" json:"productTypeCode" validate:"required,oneof=LEISURE GIFTICON RENTCAR"`                          // Product Type Code
	SalePeriod                  SalePeriod        `bson:"salePeriod" json:"salePeriod" validate:"required"`                                                                   // Product Sales Period
	ProductBriefIntroduction    string            `bson:"productBriefIntroduction" json:"productBriefIntroduction" validate:"required"`                                       // Brief Introduction
	Variants                    []Variant         `bson:"variants" json:"variants" validate:"required"`                                                                       // Variants
	Categories                  []ProductCategory `bson:"categories" json:"categories" validate:"required"`                                                                   // Categories
	Regions                     []Regional        `bson:"regions" json:"regions" validate:"required"`                                                                         // Regions
	Images                      []Image           `bson:"images" json:"images"`                                                                                               // Images
	IsIntegratedVoucher         bool              `bson:"isIntegratedVoucher" json:"isIntegratedVoucher" validate:"required"`                                                 // Integrated Voucher
	IsRefundableAfterExpiration bool              `bson:"isRefundableAfterExpiration" json:"isRefundableAfterExpiration" validate:"required"`                                 // Refundable After Expiration
	IsUsed                      bool              `bson:"isUsed" json:"isUsed" validate:"required"`
	PluDetails                  []PluDetail       `bson:"pluDetails" json:"PluDetails" validate:"required"`
	CreatedAt                   string            `bson:"createdAt" json:"createdAt" validate:"required"` // Created Timestamp
	UpdatedAt                   string            `bson:"updatedAt" json:"updatedAt"`                     // Updated Timestamp
}

type PluDetail struct {
	ProductOptionTypeCode string            `bson:"productOptionTypeCode" json:"productOptionTypeCode" validate:"oneof=SCHEDULE ROUND LIST"` // Product Option Type Code
	PLU                   []string          `bson:"plu" json:"plu" validate:"required"`
	PluHash               map[string]string `bson:"pluHash" json:"pluHash" validate:"required"`
}

type VariantPLU struct {
	VariantId   int64  `bson:"VariantId" json:"VariantId" validate:"required"`
	VariantName string `bson:"VariantName" json:"VariantName" validate:"required"`
}
