package trip

type GooglePlaceIdOfProduct struct {
	Id        string  `bson:"_id,omitempty" json:"id"`
	Latitude  float64 `bson:"latitude" json:"latitude" validate:"required"`
	Longitude float64 `bson:"longitude" json:"longitude" validate:"required"`
	ProductID int64   `bson:"productId" json:"productId" validate:"required"` // Product ID
	PlaceId   string  `bson:"placeId" json:"placeId"`
	CreatedAt string  `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt string  `bson:"updatedAt" json:"updatedAt"`
}
