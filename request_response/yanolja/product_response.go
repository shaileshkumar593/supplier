package yanolja

// Response Yanolja  response
type SubCategory Category
type SubRegional Regional

type Product struct {
	ProductID                   int                   `json:"productId"`
	ProductName                 string                `json:"productName"`
	ProductVersion              int                   `json:"productVersion"`
	Price                       ProductPrice          `json:"price"`
	ProductStatusCode           string                `json:"productStatusCode"`
	ProductTypeCode             string                `json:"productTypeCode"`
	SalePeriod                  SalePeriodS           `json:"salePeriod"`
	ProductBriefIntroduction    string                `json:"productBriefIntroduction"`
	ProductInfo                 ProductInfo           `json:"productInfo"`
	ProductOptionGroups         []ProductOptionGroups `json:"productOptionGroups"`
	SearchKeywords              []string              `json:"searchKeywords"`
	Categories                  []Category            `json:"categories"`
	Regions                     []Regional            `json:"regions"`
	Images                      []Image               `json:"images"`
	Videos                      []Video               `json:"videos"`
	Pictograms                  []Pictogram           `json:"pictograms"`
	IsCancelPenalty             bool                  `json:"isCancelPenalty"`
	IsReservationAfterPurchase  bool                  `json:"isReservationAfterPurchase"`
	PurchaseDateUsableTypeCode  string                `json:"purchaseDateUsableTypeCode"`
	IsAvailableOnPurchaseDate   bool                  `json:"isAvailableOnPurchaseDate"`
	IsIntegratedVoucher         bool                  `json:"isIntegratedVoucher"`
	IsRefundableAfterExpiration bool                  `json:"isRefundableAfterExpiration"`
	IsUsed                      bool                  `json:"isUsed"`
	SellerInfos                 []SellerInfos         `json:"sellerInfos"`
	ConvenienceTypeCodes        []string              `json:"convenienceTypeCodes"`
}

type ProductPrice struct {
	Currency      string `json:"currency"`
	RetailPrice   int    `json:"retailPrice"`
	SalePrice     int    `json:"salePrice"`
	SalePriceName string `json:"salePriceName"`
}

type SalePeriodS struct {
	StartDateTime string `json:"startDateTime"`
	EndDateTime   string `json:"endDateTime"`
	Timezone      string `json:"timezone"`
	Offset        string `json:"offset"`
}

type ProductInfo struct {
	ProductBasicInfo  string        `json:"productBasicInfo"`
	ProductUsageInfo  string        `json:"productUsageInfo"`
	FacilityInfos     FacilityInfos `json:"facilityInfos"`
	NoticeInfo        string        `json:"noticeInfo"`
	ServiceCenterInfo string        `json:"serviceCenterInfo"`
	RefundInfo        string        `json:"refundInfo"`
	VoucherUsageInfo  string        `json:"voucherUsageInfo"`
}

type FacilityInfos struct {
	FacilityID                 int      `json:"facilityId"`
	FacilityName               string   `json:"facilityName"`
	Location                   Location `json:"location"`
	PhoneNumber                string   `json:"phoneNumber"`
	AdministrativeBuildingCode int      `json:"administrativeBuildingCode"`
	IsRepFacility              bool     `json:"isRepFacility"`
	FacilityDetailInfo         string   `json:"facilityDetailInfo"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

type ProductOptionGroups struct {
	ProductOptionGroupID          int              `json:"productOptionGroupId"`
	ProductID                     int              `json:"productId"`
	ProductOptionGroupName        string           `json:"productOptionGroupName"`
	ProductOptionGroupDescription string           `json:"productOptionGroupDescription"`
	IsSchedule                    bool             `json:"isSchedule"`
	IsRound                       bool             `json:"isRound"`
	ProductOptions                []ProductOptions `json:"productOptions"`
	Variants                      []Variants       `json:"variants"`
}

type ProductOptions struct {
	ProductOptionID       int                  `json:"productOptionId"`
	ProductOptionName     string               `json:"productOptionName"`
	HierarchicalOrder     int                  `json:"hierarchicalOrder"`
	ProductOptionTypeCode string               `json:"productOptionTypeCode"`
	ProductOptionItems    []ProductOptionItems `json:"productOptionItems"`
}
type ProductOptionItems struct {
	ProductOptionItemID   int      `json:"productOptionItemId"`
	ProductOptionItemName string   `json:"productOptionItemName"`
	SortOrder             int      `json:"sortOrder,omitempty"`
	Schedules             []string `json:"schedules,omitempty"`
	Rounds                []string `json:"rounds,omitempty"`
}
type Variants struct {
	VariantID                           int                  `json:"variantId"`
	ProductID                           int                  `json:"productId"`
	VariantName                         string               `json:"variantName"`
	VariantDescription                  string               `json:"variantDescription"`
	RefundApprovalTypeCode              string               `json:"refundApprovalTypeCode"`
	QuantityPerPerson                   int                  `json:"quantityPerPerson"`
	QuantityPerPurchase                 int                  `json:"quantityPerPurchase"`
	QuantityPerPersonValidityDays       any                  `json:"quantityPerPersonValidityDays"`
	SortOrder                           int                  `json:"sortOrder"`
	IsRefundableAfterExpiration         bool                 `json:"isRefundableAfterExpiration"`
	IsAvailableOnPurchaseDate           bool                 `json:"isAvailableOnPurchaseDate"`
	RefundInfo                          string               `json:"refundInfo"`
	IsDisplay                           bool                 `json:"isDisplay"`
	OrderExpirationUsageProcessTypeCode string               `json:"orderExpirationUsageProcessTypeCode"`
	OrderExpirationDateTypeCode         string               `json:"orderExpirationDateTypeCode"`
	Price                               VariantPrice         `json:"price"`
	Fee                                 FeePolicy            `json:"fee_policy"`
	SalePeriod                          ProductSalesPeriod   `json:"salePeriod"`
	ProductOptionItems                  []ProductOptionItems `json:"productOptionItems"`
	VariantItems                        []VariantItems       `json:"variantItems"`
	VariantStatusCode                   string               `json:"variantStatusCode"`
}

type VariantPrice struct {
	Currency          string `json:"currency"`
	RetailPrice       int    `json:"retailPrice"`
	DiscountSalePrice int    `json:"discountSalePrice"`
	SalePrice         int    `json:"salePrice"`
	CostPrice         int    `json:"costPrice"`
}

type FeePolicy struct {
	FeeTypeCode      string  `json:"feeTypeCode"`
	FeeRate          any     `json:"feeRate"`
	RevenueShareRate float64 `json:"revenueShareRate"`
}

type ProductSalesPeriod struct {
	StartDateTime string `json:"startDateTime"`
	EndDateTime   string `json:"endDateTime"`
	Timezone      string `json:"timezone"`
	Offset        string `json:"offset"`
}

type VariantItems struct {
	SupplyItemID           int              `json:"supplyItemId"`
	SupplyItemName         string           `json:"supplyItemName"`
	ValidityPeriodTypeCode string           `json:"validityPeriodTypeCode"`
	ValidityPeriod         PeriodOfValidity `json:"validityPeriod"`
	ValidityDays           int              `json:"validityDays"`
	SellerInfoID           string           `json:"sellerInfoId"`
	IsVoucherUsed          bool             `json:"isVoucherUsed"`
	VoucherDisplayTypeCode string           `json:"voucherDisplayTypeCode"`
}

type PeriodOfValidity struct {
	StartDateTime string `json:"startDateTime"`
	EndDateTime   string `json:"endDateTime"`
	Timezone      string `json:"timezone"`
	Offset        string `json:"offset"`
}

type Category struct {
	CategoryCode  string        `json:"categoryCode"`
	CategoryLevel int           `json:"categoryLevel"`
	SubCategories []SubCategory `json:"subCategories"`
}

type Regional struct {
	RegionID       int           `json:"regionId"`
	RegionCode     string        `json:"regionCode"`
	ParentRegionID int           `json:"parentRegionId"`
	RegionName     string        `json:"regionName"`
	RegionLevel    int           `json:"regionLevel"`
	SubRegions     []SubRegional `json:"subRegions"`
}

type Image struct {
	ImageTypeCode string   `json:"imageTypeCode"`
	ImageUrls     []string `json:"imageUrls"`
}

type Video struct {
	VideoTypeCode          string `json:"videoTypeCode"`
	VideoURL               string `json:"videoUrl"`
	VideoThumbnailImageURL string `json:"videoThumbnailImageUrl"`
}

type Pictogram struct {
	PictogramName    string `json:"pictogramName"`
	IconName         string `json:"iconName"`
	PictogramContent string `json:"pictogramContent"`
}

type SellerInfos struct {
	ID                   string `json:"id"`
	Owner                string `json:"owner"`
	CompanyName          string `json:"companyName"`
	CompanyAddress       string `json:"companyAddress"`
	CompanyEmail         string `json:"companyEmail"`
	CompanyPhoneNumber   string `json:"companyPhoneNumber"`
	BusinessNumber       string `json:"businessNumber"`
	MailOrderSalesNumber string `json:"mailOrderSalesNumber"`
}

type VariantStock struct {
	ProductId                  int64                       `json:"productId"`
	VariantId                  int64                       `json:"variantId"`
	InventoryTypeCode          string                      `json:"inventoryTypeCode"`
	QuantityPerPerson          int32                       `json:"quantityPerPerson"`
	QuantityPerPurchase        int32                       `json:"quantityPerPurchase"`
	Quantity                   int32                       `json:"quantity"`
	VariantScheduleInventories []VariantScheduleInventorie `json:"variantScheduleInventories"`
}

type VariantScheduleInventorie struct {
	Date                    string                `json:"date"`
	Quantity                int32                 `json:"quantity"`
	VariantRoundInventories []VariantEpisodeStock `json:"variantRoundInventories"`
}

type VariantEpisodeStock struct {
	Time     string `json:"time"`
	Quantity int32  `json:"quantity"`
}

type VariantInventoryResp struct {
	VariantId                     int64 `json:"variantId"`
	QuantityPerPersonValidityDays int32 `json:"quantityPerPersonValidityDays"`
	QuantityPerPerson             int32 `json:"quantityPerPerson"`
	QuantityPerPurchase           int32 `json:"quantityPerPurchase"`
	Quantity                      int32 `json:"quantity"`
}

// ProductNameDoc is a lightweight struct just to decode productName
type ProductNameStruct struct {
	ProductName string `bson:"productName"`
}
