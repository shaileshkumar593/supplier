package yanolja

// not in use now
type ImageDetailForSync struct {
	ID             string `bson:"_id,omitempty" json:"id"`
	ProductId      int64  `bson:"productId" json:"productId"`
	TripImageId    string `bson:"tripImageId" json:"tripImageId" validate:"required"`
	ProcessedUrl   string `bson:"ProcessedUrl,omitempty" json:"ProcessedUrl,omitempty"`
	ImageTypeCode  string `bson:"imageTypeCode,omitempty" json:"imageTypeCode,omitempty"`
	Size           string `bson:"size" json:"size"`
	AspectRatio    string `bson:"aspectRatio" json:"aspectRatio"`
	PixelDimension string `bson:"pixelDimension" json:"pixelDimension"`
	Status         string `bson:"status,omitempty" json:"status,omitempty"` //Sync, Not-Sync how to accommodate multiple channel (pub-sub/kafka?)
	CreatedAt      string `bson:"createdAt" json:"createdAt"`
	UpdatedAt      string `bson:"updatedAt" json:"updatedAt"`
}
