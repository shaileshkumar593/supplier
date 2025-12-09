package common

type Response struct {
	Code    string      `json:"code"`
	TraceID string      `json:"traceId"`
	Status  int         `json:"status"`
	Body    interface{} `json:"body,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
