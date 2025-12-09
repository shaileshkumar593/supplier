package trip

import "swallow-supplier/mongo/domain/trip"

type ImageSyncToTripRequest struct {
	Message string            `json:"message"`
	Data    []ImageSyncToTrip `json:"data"`
}

type ImageSyncToTrip struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	ProductId int64  `bson:"productId" json:"productId"`
	Url       string `bson:"url,omitempty" json:"url,omitempty"`
	Status    string `bson:"status,omitempty" json:"status,omitempty"`
	ImageId   string `bson:"imageId" json:"imageId"`
}

// ProductContentApiSync  content sync for product
type ProductContentApiSync struct {
	Message string               `json:"message"`
	Data    []trip.ProuctContent `json:"data"`
}

//PackageContentApiSync  content sync to package
type PackageContentApiSync struct {
	Message string                `json:"message"`
	Data    []trip.PackageContent `json:"data"`
}
