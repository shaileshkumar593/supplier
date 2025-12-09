package trip

type PreOrderTimeoutCancellation struct {
	SequenceId string `json:"sequenceId" binding:"required"`
	OtaOrderId string `json:"otaOrderId" binding:"required"`
}
