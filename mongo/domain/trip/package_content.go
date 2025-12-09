package trip

type PackageContent struct {
	SupplierProductId string              `bson:"supplierProductId" json:"supplierProductId" validate:"required"`
	SupplierName      string              `bson:"supplierName" json:"supplierName" validate:"Required"`
	OptionList        []PackageOptionList `bson:"optionList" json:"optionList" validate:"required"`
	SyncStatus        string              `bson:"syncStatus" json:"syncStatus" validate:"Required"`
	CreatedAt         string              `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt         string              `bson:"updatedAt" json:"updatedAt" validate:"required"`
}

type PackageOptionList struct {
	Option            []PackageOption    `bson:"option" json:"option"`
	OptionStatus      string             `bson:"optionStatus" json:"optionStatus" validate:"required,oneof='active' 'inactive'"`
	OptionDescription string             `bson:"optionDescription" json:"optionDescription"`
	BookingCutOffTime BookingCutOffTime  `bson:"bookingCutOffTime" json:"bookingCutOffTime" validate:"required"`
	BookingQuestions  []BookingQuestions `bson:"bookingQuestions" json:"bookingQuestions"`
	Unit              []UnitObj          `bson:"unit" json:"unit" validate:"required"`
}

type PackageOption struct {
	OptionCode string `bson:"optionCode" json:"optionCode" validate:"required"`
	ValueCode  string `bson:"valueCode" json:"valueCode" validate:"required"`
	ValueName  string `bson:"valueName" json:"valueName" validate:"required"`
}

type BookingCutOffTime struct {
	DayBeforeVisitDate string `bson:"dayBeforeVisitDate" json:"dayBeforeVisitDate" validate:"required"`
	Time               string `bson:"time" json:"time" validate:"required"`
}

type BookingQuestions struct {
	Code           string         `bson:"code" json:"code" validate:"required"`
	Name           string         `bson:"name" json:"name" validate:"required"`
	AnswerType     string         `bson:"answerType" json:"answerType" validate:"required,oneof='Single_Selection' 'Free_Text'"`
	AllowedAnswers AllowedAnswers `bson:"allowedAnswers" json:"allowedAnswers"`
	Description    string         `bson:"description" json:"description"`
}

type AllowedAnswers struct {
	Code string `bson:"Code" json:"Code" validate:"required"`
	Name string `bson:"Name" json:"Name" validate:"required"`
}

type UnitObj struct {
	PLU            string       `bson:"plu" json:"plu" validate:"required"`
	Reference      string       `bson:"reference" json:"reference"`
	TicketTypeCode string       `bson:"ticketTypeCode" json:"ticketTypeCode" validate:"required, oneof='Adult' 'Child' 'Senior' 'Youth' 'Infant' 'Student' 'Traveler' 'Customized'"`
	CustomCode     string       `bson:"customCode" json:"customCode"`
	Restrictions   Restrictions `bson:"restrictions" json:"restrictions" validate:"required"`
	Currency       CurrencyObj  `bson:"currency" json:"currency" validate:"required"`
}

type Restrictions struct {
	MinUnits          string `bson:"MinUnits" json:"MinUnits"`
	MaxUnits          string `bson:"maxUnits" json:"maxUnits" validate:"required"`
	UnitPax           string `bson:"unitPax" json:"unitPax" validate:"required"`
	CompanionRequired string `bson:"companionRequired" json:"companionRequired" validate:"required"`
}

type CurrencyObj struct {
	NetPriceCurrency    string `bson:"netPriceCurrency" json:"netPriceCurrency" validate:"required"`
	RetailPriceCurrency string `bson:"retailPriceCurrency" json:"retailPriceCurrency" validate:"required"`
}
