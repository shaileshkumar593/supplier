package yanolja

// note ItemId and VariantId   is similar

type ItemIdDetails struct {
	Id        string `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB document ID
	OrderId   int64  `bson:"orderId" json:"orderId" validate:"required"`
	Items     []Item `bson:"items" json:"items" validate:"required"`
	CreatedAt string `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt string `bson:"updatedAt" json:"updatedAt"`
}

type Item struct {
	ItemId string `bson:"itemId" json:"itemId" validate:"required"`
	PLU    string `bson:"plu" json:"plu" validate:"required"`
}
