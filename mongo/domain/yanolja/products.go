package yanolja

// Product represents the MongoDB document structure for the product.
type Product struct {
	Id                          string               `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB document ID
	SupplierName                string               `bson:"supplierName" json:"supplierName" validate:"required"`
	ProductID                   int64                `bson:"productId" json:"productId" validate:"required"`                                                                     // Product ID
	ProductName                 string               `bson:"productName" json:"productName" validate:"required"`                                                                 // Product Name
	ProductVersion              int32                `bson:"productVersion" json:"productVersion" validate:"required"`                                                           // Product Version
	Price                       Price                `bson:"price" json:"price" validate:"required"`                                                                             // Product Price
	ProductStatusCode           string               `bson:"productStatusCode" json:"productStatusCode" validate:"required,oneof=WAITING_FOR_SALE IN_SALE SOLD_OUT END_OF_SALE"` // Product Status Code
	ProductTypeCode             string               `bson:"productTypeCode" json:"productTypeCode" validate:"required,oneof=LEISURE GIFTICON RENTCAR"`                          // Product Type Code
	SalePeriod                  SalePeriod           `bson:"salePeriod" json:"salePeriod" validate:"required"`                                                                   // Product Sales Period
	ProductBriefIntroduction    string               `bson:"productBriefIntroduction" json:"productBriefIntroduction" validate:"required"`                                       // Brief Introduction
	ProductInfo                 ProductInfo          `bson:"productInfo" json:"productInfo" validate:"required"`                                                                 // Product Information
	ProductOptionGroups         []ProductOptionGroup `bson:"productOptionGroups" json:"productOptionGroups" validate:"required"`                                                 // Product Option Groups
	SearchKeywords              []string             `bson:"searchKeywords,omitempty" json:"searchKeywords,omitempty"`
	Categories                  []ProductCategory    `bson:"categories" json:"categories" validate:"required"`
	Regions                     []Regional           `bson:"regions" json:"regions" validate:"required"`
	Images                      []Image              `bson:"images" json:"images"`
	TextFromImages              []TextFromImage      `bson:"textFromImages" json:"textFromImages"`
	Videos                      []Video              `bson:"videos" json:"videos"`
	Pictograms                  []Pictogram          `bson:"pictograms" json:"pictograms"`
	IsCancelPenalty             bool                 `bson:"isCancelPenalty" json:"isCancelPenalty" validate:"required"`
	IsReservationAfterPurchase  bool                 `bson:"isReservationAfterPurchase" json:"isReservationAfterPurchase" validate:"required"`
	PurchaseDateUsableTypeCode  string               `bson:"purchaseDateUsableTypeCode" json:"purchaseDateUsableTypeCode" validate:"required,oneof=ALL NONE CUSTOM"`
	IsAvailableOnPurchaseDate   bool                 `bson:"isAvailableOnPurchaseDate" json:"isAvailableOnPurchaseDate" validate:"required"`
	IsIntegratedVoucher         bool                 `bson:"isIntegratedVoucher" json:"isIntegratedVoucher" validate:"required"`
	IsRefundableAfterExpiration bool                 `bson:"isRefundableAfterExpiration" json:"isRefundableAfterExpiration" validate:"required"`
	IsUsed                      bool                 `bson:"isUsed" json:"isUsed" validate:"required"`
	SellerInfos                 []SellerInfo         `bson:"sellerInfos,omitempty" json:"sellerInfos,omitempty"`
	ConvenienceTypeCode         []string             `bson:"convenienceTypeCode,omitempty" json:"convenienceTypeCode,omitempty"`
	ImageScheduleStatus         bool                 `bson:"imageScheduleStatus" json:"imageScheduleStatus"`
	ViewScheduleStatus          bool                 `bson:"viewScheduleStatus" json:"viewScheduleStatus"`
	ContentScheduleStatus       bool                 `bson:"contentScheduleStatus" json:"contentScheduleStatus"`
	OodoSyncStatus              bool                 `bson:"oodoSyncStatus" json:"oodoSyncStatus"`
	CreatedAt                   string               `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt                   string               `bson:"updatedAt" json:"updatedAt"`
}

type TextFromImage struct {
	ImageTypeCode string   `bson:"imageTypeCode,omitempty" json:"imageTypeCode,omitempty" validate:"oneof=THUMBNAIL ROLLING DETAIL"`
	ImageUrls     []string `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Text          []string `bson:"text,omitempty" json:"text,omitempty"`
}

// Image represents an image associated with a product
type Image struct {
	ImageTypeCode string   `bson:"imageTypeCode,omitempty" json:"imageTypeCode,omitempty" validate:"oneof=THUMBNAIL ROLLING DETAIL"`
	ImageURLs     []string `bson:"imageUrls,omitempty" json:"imageUrls,omitempty"`
}

// Video represents a video associated with a product
type Video struct {
	VideoTypeCode          string `bson:"videoTypeCode" json:"videoTypeCode" validate:"required,oneof=ROLLING"`
	VideoURL               string `bson:"videoUrl" json:"videoUrl" validate:"required"`
	VideoThumbnailImageURL string `bson:"videoThumbnailImageUrl" json:"videoThumbnailImageUrl" validate:"required"`
}

// Pictogram represents pictograms related to a product
type Pictogram struct {
	PictogramName    string `bson:"pictogramName,omitempty" json:"pictogramName,omitempty"`
	IconName         string `bson:"iconName,omitempty" json:"iconName,omitempty"`
	PictogramContent string `bson:"pictogramContent,omitempty" json:"pictogramContent,omitempty"`
}

// SellerInfo represents information about a seller
type SellerInfo struct {
	ID                   string `bson:"id,omitempty" json:"id,omitempty"`
	Owner                string `bson:"owner,omitempty" json:"owner,omitempty"`
	CompanyName          string `bson:"companyName,omitempty" json:"companyName,omitempty"`
	CompanyAddress       string `bson:"companyAddress,omitempty" json:"companyAddress,omitempty"`
	CompanyEmail         string `bson:"companyEmail,omitempty" json:"companyEmail,omitempty"`
	CompanyPhoneNumber   string `bson:"companyPhoneNumber,omitempty" json:"companyPhoneNumber,omitempty"`
	BusinessNumber       string `bson:"businessNumber,omitempty" json:"businessNumber,omitempty"`
	MailOrderSalesNumber string `bson:"mailOrderSalesNumber,omitempty" json:"mailOrderSalesNumber,omitempty"`
}

// Region represents a regional classification
type Regional struct {
	RegionID       int64      `bson:"regionId,omitempty" json:"regionId,omitempty"`
	RegionCode     string     `bson:"regionCode,omitempty" json:"regionCode,omitempty"`
	ParentRegionID int64      `bson:"parentRegionId,omitempty" json:"parentRegionId,omitempty"`
	RegionName     string     `bson:"regionName,omitempty" json:"regionName,omitempty"`
	RegionLevel    int32      `bson:"regionLevel,omitempty" json:"regionLevel,omitempty"`
	SubRegions     []Regional `bson:"subRegions,omitempty" json:"subRegions,omitempty"`
}

// Category represents a category of a product
type ProductCategory struct {
	CategoryCode  string            `bson:"categoryCode,omitempty" json:"categoryCode,omitempty"`
	CategoryLevel int32             `bson:"categoryLevel,omitempty" json:"categoryLevel,omitempty"`
	SubCategories []ProductCategory `bson:"subCategories,omitempty" json:"subCategories,omitempty"`
}

// Price represents the price details for the product.
type Price struct {
	Currency      string  `bson:"currency" json:"currency" validate:"required"`                   // Currency Code
	RetailPrice   float64 `bson:"retailPrice" json:"retailPrice"`                                 // MSRP
	SalePrice     float64 `bson:"salePrice" json:"salePrice" validate:"required"`                 // Sale Price
	SalePriceName float64 `bson:"discountSalePrice" json:"discountSalePrice" validate:"required"` // Discount Sale Price
}

// SalePeriod represents the sales period for the product.
type SalePeriod struct {
	StartDateTime string `bson:"startDateTime,omitempty" json:"startDateTime,omitempty"` // UTC start date
	EndDateTime   string `bson:"endDateTime,omitempty" json:"endDateTime,omitempty"`     // UTC end date
	TimeZone      string `bson:"timezone,omitempty" json:"timezone,omitempty"`           // Time Zone
	Offset        string `bson:"offset,omitempty" json:"offset,omitempty"`               // Offset
}

// ProductInfo contains the information about the product.
type ProductInfo struct {
	ProductBasicInfo  string         `bson:"productBasicInfo" json:"productBasicInfo" validate:"required"` // Basic Product Information
	ProductUsageInfo  string         `bson:"productUsageInfo" json:"productUsageInfo" validate:"required"` // Product Usage Information
	FacilityInfos     []FacilityInfo `bson:"facilityInfos" json:"facilityInfos"`                           // List of Facility Infos (optional)
	NoticeInfo        string         `bson:"noticeInfo,omitempty" json:"noticeInfo,omitempty"`
	ServiceCenterInfo string         `bson:"serviceCenterInfo,omitempty" json:"serviceCenterInfo,omitempty"`
	RefundInfo        string         `bson:"refundInfo,omitempty" json:"refundInfo,omitempty"`
	VoucherUsageInfo  string         `bson:"voucherUsageInfo,omitempty" json:"voucherUsageInfo,omitempty"`
}

// FacilityInfo contains information about facilities related to the product.
type FacilityInfo struct {
	FacilityID                 int64    `bson:"facilityId,omitempty" json:"facilityId,omitempty"`                                 // Facility ID
	FacilityName               string   `bson:"facilityName,omitempty" json:"facilityName,omitempty"`                             // Facility Name
	Location                   Location `bson:"location" json:"location,omitempty"`                                               // Location Information
	PhoneNumber                string   `bson:"phoneNumber,omitempty" json:"phoneNumber,omitempty"`                               // Facility Phone Number
	AdministrativeBuildingCode string   `bson:"administrativeBuildingCode,omitempty" json:"administrativeBuildingCode,omitempty"` // Administrative Code (optional)
	IsRepFacility              bool     `bson:"isRepFacility,omitempty" json:"isRepFacility,omitempty"`                           // Whether the facility represents the product
	FacilityDetailInfo         string   `bson:"facilityDetailInfo,omitempty" json:"facilityDetailInfo,omitempty"`                 // Facility Details (optional)
	NoticeInfo                 string   `bson:"noticeInfo,omitempty" json:"noticeInfo,omitempty"`                                 // Additional Information (optional)
	ServiceCenterInfo          string   `bson:"serviceCenterInfo,omitempty" json:"serviceCenterInfo,omitempty"`                   // Customer Service Information (optional)
	RefundInfo                 string   `bson:"refundInfo,omitempty" json:"refundInfo,omitempty"`                                 // Refund Information (optional)
	VoucherUsageInfo           string   `bson:"voucherUsageInfo,omitempty" json:"voucherUsageInfo,omitempty"`                     // Voucher Usage Information (optional)
}

// Location represents location details for a facility.
type Location struct {
	Longitude float64 `bson:"longitude" json:"longitude,omitempty"` // Longitude
	Latitude  float64 `bson:"latitude" json:"latitude,omitempty"`   // Latitude
	Address   string  `bson:"address" json:"address,omitempty"`     // Address
}

// ProductOptionGroup represents a group of product options.
type ProductOptionGroup struct {
	ProductOptionGroupID          int64           `bson:"productOptionGroupId" json:"productOptionGroupId" validate:"required"`                   // Product Option Group ID
	ProductID                     int64           `bson:"productId" json:"productId" validate:"required"`                                         // Product ID
	ProductOptionGroupName        string          `bson:"productOptionGroupname" json:"productOptionGroupname" validate:"required"`               // Product Option Group Name
	ProductOptionGroupDescription string          `bson:"productOptionGroupDescription,omitempty" json:"productOptionGroupDescription,omitempty"` // Product Option Group Description (optional)
	IsSchedule                    bool            `bson:"isSchedule" json:"isSchedule" validate:"required"`                                       // Is Schedule
	IsRound                       bool            `bson:"isRound" json:"isRound" validate:"required"`                                             // Is Round
	ProductOptions                []ProductOption `bson:"productOptions" json:"productOptions" validate:"required"`
	Variants                      []Variant       `bson:"variants" json:"variants"` // Product Options
}

// ProductOption represents a product option.
type ProductOption struct {
	ProductOptionID       int64               `bson:"productOptionId,omitempty" json:"productOptionId,omitempty"`                              // Product Option ID
	ProductOptionName     string              `bson:"productOptionName,omitempty" json:"productOptionName,omitempty"`                          // Product Option Name
	HierarchicalOrder     int32               `bson:"hierarchicalOrder,omitempty" json:"hierarchicalOrder,omitempty"`                          // Hierarchical Order
	ProductOptionTypeCode string              `bson:"productOptionTypeCode" json:"productOptionTypeCode" validate:"oneof=SCHEDULE ROUND LIST"` // Product Option Type Code
	ProductOptionItems    []ProductOptionItem `bson:"productOptionItems,omitempty" json:"productOptionItems,omitempty"`                        // Product Option Items
}

// ProductOptionItem represents an item in a product option.
type ProductOptionItem struct {
	ProductOptionItemID   int64    `bson:"productOptionItemId" json:"productOptionItemId" validate:"required"`     // Product Option Item ID
	ProductOptionItemName string   `bson:"productOptionItemName" json:"productOptionItemName" validate:"required"` // Product Option Item Name
	SortOrder             int32    `bson:"sortOrder" json:"sortOrder" validate:"required"`                         // Sort Order
	Schedules             []string `bson:"schedules,omitempty" json:"schedules,omitempty"`                         // Schedules (optional)
	Rounds                []string `bson:"rounds,omitempty" json:"rounds,omitempty"`                               // Rounds (optional)
}

// Variant represents a variant of a product
type Variant struct {
	VariantID                           int64               `bson:"variantId,omitempty" json:"variantId,omitempty"`
	ProductID                           int64               `bson:"productId,omitempty" json:"productId,omitempty"`
	VariantName                         string              `bson:"variantName,omitempty" json:"variantName,omitempty"`
	VariantDescription                  string              `bson:"variantDescription,omitempty" json:"variantDescription,omitempty"`
	RefundApprovalTypeCode              string              `bson:"refundApprovalTypeCode,omitempty" json:"refundApprovalTypeCode" validate:"oneof=DIRECT ADMIN"`
	QuantityPerPerson                   int32               `bson:"quantityPerPerson,omitempty" json:"quantityPerPerson,omitempty"`
	QuantityPerPurchase                 int32               `bson:"quantityPerPurchase,omitempty" json:"quantityPerPurchase,omitempty"`
	QuantityPerPersonValidityDays       int32               `bson:"quantityPerPersonValidityDays,omitempty" json:"quantityPerPersonValidityDays,omitempty"`
	SortOrder                           int32               `bson:"sortOrder,omitempty" json:"sortOrder,omitempty"`
	IsRefundableAfterExpiration         bool                `bson:"isRefundableAfterExpiration,omitempty" json:"isRefundableAfterExpiration,omitempty"`
	IsAvailableOnPurchaseDate           bool                `bson:"isAvailableOnPurchaseDate,omitempty" json:"isAvailableOnPurchaseDate,omitempty"`
	RefundInfo                          string              `bson:"refundInfo,omitempty" json:"refundInfo,omitempty"`
	IsDisplay                           bool                `bson:"isDisplay,omitempty" json:"isDisplay,omitempty"`
	OrderExpirationUsageProcessTypeCode string              `bson:"orderExpirationUsageProcessTypeCode,omitempty" json:"orderExpirationUsageProcessTypeCode,omitempty" validate:"oneof=NONE FORCED_USE"`
	OrderExpirationDateTypeCode         string              `bson:"orderExpirationDateTypeCode,omitempty" json:"orderExpirationDateTypeCode,omitempty"`
	Price                               VariantPrice        `bson:"price" json:"price"`
	Fee                                 Fee                 `bson:"fee" json:"fee"`
	SalePeriod                          SalePeriod          `bson:"salePeriod" json:"salePeriod"`
	ProductOptionItems                  []ProductOptionItem `bson:"productOptionItems" json:"productOptionItems"`
	VariantItems                        []VariantItem       `bson:"variantItems" json:"variantItems"`
	VariantStatusCode                   string              `bson:"variantStatusCode" json:"variantStatusCode"`
}

// VariantPrice represents the pricing of a product variant
type VariantPrice struct {
	Currency          string  `bson:"currency" json:"currency" validate:"required"`
	RetailPrice       float64 `bson:"retailPrice,omitempty" json:"retailPrice,omitempty"`
	DiscountSalePrice float64 `bson:"discountSalePrice" json:"discountSalePrice" validate:"required"`
	SalePrice         float64 `bson:"salePrice" json:"salePrice" validate:"required"`
	CostPrice         float64 `bson:"costPrice" json:"costPrice" validate:"required"`
}

// Fee represents the fee policy for a variant
type Fee struct {
	FeeTypeCode      string  `bson:"feeTypeCode" json:"feeTypeCode" validate:"oneof=FIXED_RATE REVENUE_SHARE"`
	FeeRate          float64 `bson:"feeRate,omitempty" json:"feeRate,omitempty"`
	RevenueShareRate float64 `bson:"revenueShareRate,omitempty" json:"revenueShareRate,omitempty"`
}

// VariantItem represents an item within a variant
type VariantItem struct {
	SupplyItemID           int64            `bson:"supplyItemId,omitempty" json:"supplyItemId,omitempty"`
	SupplyItemName         string           `bson:"supplyItemName,omitempty" json:"supplyItemName,omitempty"`
	ValidityPeriodTypeCode string           `bson:"validityPeriodTypeCode,omitempty" json:"validityPeriodTypeCode,omitempty" validate:"oneof=FIX BUY"`
	ValidityPeriod         PeriodOfValidity `bson:"validityPeriod,omitempty" json:"validityPeriod,omitempty"`
	ValidityDays           int32            `bson:"validityDays,omitempty" json:"validityDays,omitempty"`
	SellerInfoID           string           `bson:"sellerInfoId,omitempty" json:"sellerInfoId,omitempty"`
	IsVoucherUsed          bool             `bson:"isVoucherUsed,omitempty" json:"isVoucherUsed,omitempty"`
	VoucherDisplayTypeCode string           `bson:"voucherDisplayTypeCode" json:"voucherDisplayTypeCode" validate:"oneof=NONE BARCODE QR"`
}

type PeriodOfValidity struct {
	StartDateTime string `bson:"startDateTime,omitempty" json:"startDateTime,omitempty"` // UTC start date
	EndDateTime   string `bson:"endDateTime,omitempty" json:"endDateTime,omitempty"`     // UTC end date
	TimeZone      string `bson:"timezone,omitempty" json:"timezone,omitempty"`           // Time Zone
	Offset        string `bson:"offset,omitempty" json:"offset,omitempty"`               // Offset
}
