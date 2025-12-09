package trip

type ResponseForTripRequest struct {
	SequenceID      string `bson:"sequenceId" json:"sequenceId"`
	OtaOrderID      string `bson:"otaOrderId" json:"otaOrderId"`
	RequestCategory string `bson:"requestCategory" json:"requestCategory"`
}
