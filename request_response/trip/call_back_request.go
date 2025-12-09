package trip

type CallBackRequest struct {
	Message       string       `json:"message" binding:"required"`
	SequenceIdGen string       `json:"sequenceIdGen" binding:"required"`
	TripResp      TripResponse `json:"TripResp" validate:"required"`
}
