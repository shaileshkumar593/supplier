package yanolja

import "swallow-supplier/mongo/domain/yanolja"

type ProductImages struct {
	ProductID int64           `json:"productId"`
	Images    []yanolja.Image `json:"images"`
}
