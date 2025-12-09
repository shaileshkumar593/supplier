package trip

import "swallow-supplier/mongo/domain/yanolja"

type CallBackDetail struct {
	Id                  string          `bson:"_id,omitempty" json:"id,omitempty"`
	ChannelCallBackInfo CallBackRequest `bson:"channelCallBackInfo" json:"channelCallBackInfo" validate:"required"`
	Supplier            string          `bson:"supplier" json:"supplier" validate:"required"`
	Distributor         string          `bson:"distributor" json:"distributor" validate:"required"`
	GGTToChannelStatus  string          `bson:"ggtToChannelStatus" json:"ggtToChannelStatus" oneof:"'PROCESSING' 'FAILED' 'SUCCESS'"`
	ChannelToGGTStatus  string          `bson:"channelToGGTStatus" json:"channelToGGTStatus" oneof:"'FAILED' 'SUCCESS'"`
	CreatedAt           string          `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt           string          `bson:"updatedAt" json:"updatedAt"`
}

type CallBackRequest struct {
	Message       string        `json:"message" binding:"required"`
	SequenceIdGen string        `json:"sequenceIdGen"`
	ChannelReq    ChannelReqest `json:"channelReq" validate:"required"`
}
type ChannelReqest struct {
	Items yanolja.ItemIdDetails `json:"items,omitempty"`
	Order yanolja.Model         `json:"order" validate:"required"`
}
