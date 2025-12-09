package trip

type ProductsImageUrl struct {
	ProductId int64   `bson:"productId" json:"productId"  validate:"required"`
	Images    []Image `bson:"images" json:"images"`
}

type Image struct {
	ImageTypeCode string              `json:"imageTypeCode,omitempty" validate:"oneof=THUMBNAIL ROLLING DETAIL"`
	ImageURLs     []map[string]string `json:"imageUrls,omitempty"`
}
