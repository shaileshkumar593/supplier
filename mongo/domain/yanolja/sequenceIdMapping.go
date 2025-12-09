package yanolja

// need to be checked in future for relevant use cases need discussion
type SequenceIdDetail struct {
	Id                  string `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB document ID
	SequenceId          string `bson:"sequenceId,,omitempty" json:"sequenceId" default:""`
	SequenceIdGen       string `bson:"sequenceIdGen,,omitempty" json:"sequenceIdGen" default:""`
	OrderId             int64  `bson:"orderId" json:"orderId" validate:"required"`
	PartnerOrderGroupID string `bson:"partnerOrderGroupID" json:"partnerOrderGroupID" validate:"required"`
	PartnerOrderID      string `bson:"partnerOrderId" json:"partnerOrderId" validate:"required, min=0, max <=50"`
	CreatedAt           string `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt           string `bson:"updatedAt" json:"updatedAt"`
}
