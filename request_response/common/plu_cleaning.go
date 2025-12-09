package common

type ProductValidityAndVersion struct {
	ProductID      int64      `bson:"productId" json:"productId" `
	ProductName    string     `bson:"productName" json:"productName" validate:"required"`
	ProductVersion int        `bson:"productVersion" json:"productVersion"`
	SalePeriod     SalePeriod `bson:"salePeriod" json:"salePeriod"`
}

type SalePeriod struct {
	StartDateTime string `bson:"startDateTime,omitempty" json:"startDateTime,omitempty"` // UTC start date
	EndDateTime   string `bson:"endDateTime,omitempty" json:"endDateTime,omitempty"`     // UTC end date
	TimeZone      string `bson:"timezone,omitempty" json:"timezone,omitempty"`           // Time Zone
	Offset        string `bson:"offset,omitempty" json:"offset,omitempty"`               // Offset
}
