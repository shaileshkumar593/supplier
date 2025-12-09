package trip

import "swallow-supplier/mongo/domain/yanolja"

//CreatePreOrder  is used  to return response for CreatePreOrder
type CreatePreOrder struct {
	PLU   map[string]string `json:"PLU"`
	Order yanolja.Model     `json:"order"`
}
