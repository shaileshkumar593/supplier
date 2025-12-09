package yanolja

type EverlandGetRequest struct {
	ChannelCode   string `json:"ChannelCode" binding:"required" validate:"required"`
	CustomerEmail string `json:"customerEmail" validate:"required"`
}
