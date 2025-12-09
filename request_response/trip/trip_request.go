package trip

type SwallowRequest struct {
	Header        SwallowHeader `json:"header" binding:"required"`
	DecryptedBody string        `json:"body" binding:"required"`
}
type SwallowHeader struct {
	ServiceName string `json:"serviceName" binding:"required"`
	RequestTime string `json:"requestTime" binding:"required"`
}
