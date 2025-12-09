package yanolja

type ImageUrlForProcessing struct {
	Id            string `bson:"_id,omitempty" json:"id"`
	ProductId     int64  `bson:"productId" json:"productId"  validate:"required"`
	SupplierName  string `bson:"supplierName" json:"supplierName" validate:"required"`
	Url           string `bson:"url,omitempty" json:"url,omitempty"`
	ImageTypeCode string `bson:"imageTypeCode,omitempty" json:"imageTypeCode,omitempty" validate:"oneof=THUMBNAIL ROLLING DETAIL"`
	Status        string `bson:"status,omitempty" json:"status,omitempty"` // sync or not-sync
	CreatedAt     string `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt     string `bson:"updatedAt" json:"updatedAt"`
}
