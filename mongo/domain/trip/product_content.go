package trip

type ProuctContent struct {
	Id                 string                `bson:"_id,omitempty" json:"id"`
	SupplierProductId  string                `bson:"supplierProductId" json:"supplierProductId" validate:"required,max=40"`
	SupplierName       string                `bson:"supplierName" json:"supplierName" validate:"Required"`
	Reference          string                `bson:"reference,omitempty" json:"reference,omitempty"`
	ContractId         int                   `bson:"contractId" json:"contractId" validate:"required"`
	PrimaryLanguage    string                `bson:"primaryLanguage" json:"primaryLanguage" validate:"required,oneof='zh-CN' 'zh-HK' 'en-US' 'ko-KR' 'ja-JP' 'th-TH'"` // Language code`
	Status             string                `bson:"status" json:"status" validate:"required,oneof='active' 'inactive'"`                                               // Product status`
	Category           []ProductCategory     `bson:"category" json:"category" validate:"required,min=1,dive,required"`
	Tags               []ProductTags         `bson:"tags" json:"tags" validate:"required"`
	Title              string                `bson:"title" json:"title" validate:"required,max=100"`
	Poi                []ProductPoi          `bson:"poi,omitempty" json:"poi,omitempty" validate:"omitempty,max=100,dive"`
	Destination        []DestinationObj      `bson:"destination,omitempty" json:"destination,omitempty" validate:"omitempty,max=10,dive"`
	Departure          []Departure           `bson:"departure" json:"departure" validate:"required,max=100"`
	TicketInfo         TicketInfoObj         `bson:"ticketInfo" json:"ticketInfo" validate:"required"`
	RedemptionInfo     RedemptionInfoObj     `bson:"redemptionInfo" json:"redemptionInfo" validate:"required"`
	ServiceLanguage    []ServiceLanguageObj  `bson:"serviceLanguage" json:"serviceLanguage" validate:"required"`
	Gallery            []ProductImageObj     `bson:"gallery" json:"gallery" validate:"required,max=20"`
	Highlight          []string              `bson:"highlight" json:"highlight" validate:"required,max=3"`
	Description        string                `bson:"description" json:"description" validate:"required,max=8000"`
	HowToUse           []string              `bson:"howToUse,omitempty" json:"howToUse,omitempty" validate:"required,max=20"`
	AdditionalInfo     string                `bson:"additionalInfo" json:"additionalInfo" validate:"required,max=4000"`
	GuestInformation   GuestInformationObj   `bson:"guestInformation" json:"guestInformation"`
	BookingSettings    BookingSettingsObj    `bson:"bookingSettings" json:"bookingSettings" validate:"required"`
	CancellationPolicy CancellationPolicyObj `bson:"cancellationPolicy" json:"cancellationPolicy" validate:"required"`
	TicketType         []TicketTypeObj       `bson:"ticketType" json:"ticketType" validate:"required, oneof 8 16 24"`
	Option             []ProductOptionObj    `bson:"option" json:"option"`
	MetaData           string                `bson:"metaData" json:"metaData" validate:"max=8000"`
	SyncStatus         string                `bson:"syncStatus" json:"syncStatus" validate:"required"`
	CreatedAt          string                `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt          string                `bson:"updatedAt" json:"updatedAt" validate:"required"`
}

type ProductCategory struct {
	Code string `bson:"code" json:"code" validate:"required,oneof=ATTRACTION_TICKET SHOW_EVENT_TICKET MULTI_ATTRACTION_TICKETS CRUISE_TICKET DAY_TOUR HELICOPTER_TOURS HOT_AIR_BALLOON_RIDES PARACHUTING JUNGLE_ZIPLINING DIVING ROCK_CLIMBING WI_FI SPEEDBOATING SEAPLANE_RIDES SAILING TOUR_GUIDE BOAT_SPEEDBOAT_TICKETS AIRPORT_TRANSFERS POINT_TO_POINT_TRANSFER MEALS GOLF CYCLING OFF_ROADING CLASSES_COURSES SKIING PARAGLIDING CANOEING_KAYAKING SPA_TREATMENTS TRAVEL_PHOTO_SHOOTS SIGHTSEEING_BUS_RIDES SOUVENIRS"`
}

type ProductTags struct {
	Id      string `bson:"id" json:"id" validate:"required,,max=16"`
	TagName string `bson:"tagName" json:"tagName" validate:"required,max=128"`
}

type ProductPoi struct {
	SupplierPOI   SupplierPOIObj `bson:"supplierPOI,omitempty" json:"supplierPOI,omitempty"`
	GooglePlaceId string         `bson:"googlePlaceId,omitempty" json:"googlePlaceId,omitempty"`
}

type SupplierPOIObj struct {
	SupplierId      string                  `bson:"supplierId" json:"supplierId" validate:"required"`
	MappingElements LocationMappingElements `bson:"mappingElements" json:"mappingElements" validate:"required"`
}

type LocationMappingElements struct {
	Name          string  `bson:"name" json:"name" validate:"required,max=100"`
	AddressDetail string  `bson:"addressDetail" json:"addressDetail" validate:"required,max=100"`
	Latitude      float64 `bson:"latitude" json:"latitude" validate:"required"`
	Longitude     float64 `bson:"longitude" json:"longitude" validate:"required"`
}

type DestinationObj struct {
	SupplierDestination SupplierDestinationObj `bson:"supplierDestination" json:"supplierDestination" validate:"required"`
	GooglePlaceId       string                 `bson:"googlePlaceId" json:"googlePlaceId"`
}

type SupplierDestinationObj struct {
	SupplierId      string                  `bson:"supplierId" json:"supplierId" validate:"required,max=200"`
	MappingElements LocationMappingElements `bson:"mappingElements" json:"mappingElements" validate:"required"`
}

type Departure struct {
	SupplierDeparture SupplierDepartureObj `bson:"supplierDeparture" json:"supplierDeparture" validate:"required"`
	GooglePlaceId     string               `bson:"googlePlaceId" json:"googlePlaceId" validate:"required"`
}

type SupplierDepartureObj struct {
	SupplierId      string                  `bson:"supplierId" json:"supplierId" validate:"required,max=200"`
	MappingElements LocationMappingElements `bson:"mappingElements" json:"mappingElements" validate:"required"`
}

type TicketInfoObj struct {
	DeliveryMethods string `bson:"deliveryMethods" json:"deliveryMethods" validate:"required,oneof=DIGITAL PRINT VALID_ID"`
}

type RedemptionInfoObj struct {
	RedemptionType     string               `bson:"redemptionType" json:"redemptionType" validate:"required,oneof Direct_Entry Need_Ticket_Exchange Meet_at_Start_Point Pick_Up_Everyone"`
	RedemptionLocation []RedemptionLocation `bson:"redemptionLocation" json:"redemptionLocation" validate:"max=20"`
	Description        string               `bson:"description" json:"description" validate:"required,max=400"`
}

type RedemptionLocation struct {
	SupplierLocation SupplierLocationObj `bson:"supplierLocation" json:"supplierLocation"`
	GooglePlaceId    string              `bson:"" json:"" validate:"required"`
}

type SupplierLocationObj struct {
	SupplierId      string                  `bson:"supplierId" json:"supplierId" validate:"required,max=200"`
	MappingElements LocationMappingElements `bson:"mappingElements" json:"mappingElements" validate:"required"`
}

type ServiceLanguageObj struct {
	LanguageCode string `bson:"languageCode" json:"languageCode" validate:"required,oneof=zh-CN zh-CAN en th id ms fil vi ko ja de fr es it ru pt"`
}

type ProductImageObj struct {
	TripImageId string `bson:"tripImageId" json:"tripImageId" validate:"required"`
}

type GuestInformationObj struct {
	Type string   `bson:"type" json:"type" validate:"required,oneof=PER_PERSON PER_ORDER"`
	Code []string `bson:"code" json:"code" validate:"required,oneof=GUEST_NAME COUNTRY BIRTH_DATE PassportNo"`
}

type BookingSettingsObj struct {
	BookingType             BookingTypeObj `bson:"bookingType" json:"bookingType" validate:"required"`
	PaymentConfirmationTime int            `bson:"paymentConfirmationTime" json:"paymentConfirmationTime" validate:"required"`
}

type BookingTypeObj struct {
	DateType  string       `bson:"dateType" json:"dateType" validate:"required,oneof=DATE_REQUIRED DATE_NOT_REQUIRED"`
	DateLimit DateLimitObj `bson:"dateLimit" json:"dateLimit" validate:"required"`
}

type DateLimitObj struct {
	DateLimitType       string       `bson:"dateLimitType" json:"dateLimitType" validate:"required,oneof=Single_date Multi_date Customized Unlimited"`
	MultiDateDuration   int          `bson:"multiDateDuration, omitempty" json:"multiDateDuration,omitempty"`
	CustomizedDateRange DateRangeObj `bson:"customizedDateRange" json:"customizedDateRange"`
}

type DateRangeObj struct {
	FromDate string `bson:"fromDate" json:"fromDate" validate:"required,datetime=2006-01-02"`
	ToDate   string `bson:"toDate" json:"toDate" validate:"required,datetime=2006-01-02"`
}

type CancellationPolicyObj struct {
	Type             string            `bson:"type" json:"type" validate:"required,oneof=Non_Cancellable Free_Cancel By_Visit_Date"`
	RateList         []CancellationFee `bson:"rateList" json:"rateList" validate:"max=5"`
	ConfirmationTime int               `bson:"confirmationTime" json:"confirmationTime" validate:"required"`
}

type CancellationFee struct {
	DayBeforeVisitDate int    `bson:"dayBeforeVisitDate" json:"dayBeforeVisitDate" validate:"required"`
	Time               string `bson:"time" json:"time" validate:"required"`
	Unit               string `bson:"unit" json:"unit" validate:"required"`
	Value              int    `bson:"value" json:"value" validate:"required"`
}

type TicketTypeObj struct {
	Code         string        `bson:"code,omitempty" json:"code,omitempty" validate:"required,oneof=Adult Child Senior Youth Infant Student Traveler Customized"`
	CustomCode   string        `bson:"customCode,omitempty" json:"customCode,omitempty" validate:"max=20"`
	CustomName   string        `bson:"customName,omitempty" json:"customName,omitempty" validate:"max=30"`
	Restrictions CrowdLimitObj `bson:"restrictions,omitempty" json:"restrictions,omitempty"`
	Description  string        `bson:"description,omitempty" json:"description,omitempty"`
}

type CrowdLimitObj struct {
	MinAge int `bson:"minAge,omitempty" json:"minAge,omitempty" validate:"min=1"`
	MaxAge int `bson:"maxAge,omitempty" json:"maxAge,omitempty" validate:"min=1"`
}

type ProductOptionObj struct {
	OptionCode string `bson:"optionCode" json:"optionCode" validate:"required,max=4,oneof=Option Time_Slot"`
}
