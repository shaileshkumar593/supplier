package yanolja

import (
	"swallow-supplier/mongo/domain/yanolja"
)

type ProductsImageUrl struct {
	ProductId int64           `json:"productId"  validate:"required"`
	Images    []yanolja.Image `bson:"images" json:"images"`
}

type ProductsImageUrlSync struct {
	Message string             `json:"message"`
	Data    []ProductsImageUrl `json:"data"`
}
