package yanolja

type PluRequest struct {
	ProductId      int64 `json:"productId"`
	VariantId      int64 `json:"variantId"`
	ProductVersion int32 `json:"productVersion"`
}
