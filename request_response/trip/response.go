package trip

import "swallow-supplier/mongo/domain/yanolja"

type TripResponse struct {
	Items   yanolja.ItemIdDetails `json:"items,omitempty"`
	Order   yanolja.Model         `json:"order" validate:"required"`
	Message string                `json:"message,omitempty"`
}
