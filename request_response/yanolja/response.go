package yanolja

type Response struct {
	Code        string       `json:"code"`
	TraceID     string       `json:"traceId"`
	Body        interface{}  `json:"body,omitempty"`
	Page        NumberOfPage `json:"page"`
	Collection  bool         `json:"collection"`
	ContentType *string      `json:"contentType"`
}

type SupplierResponse struct {
	TraceID     string       `json:"traceId"`
	Body        interface{}  `json:"body,omitempty"`
	Page        NumberOfPage `json:"page"`
	Collection  bool         `json:"collection"`
	ContentType *string      `json:"contentType"`
}

/*
	 type ReconcilationOrderResp struct {
		Body        interface{}  `json:"body, omitempty"`
		Page        NumberOfPage `json:"page"`
		Collection  bool         `json:"collection"`
		ContentType string       `json:"contentType"`
	}
*/
type NumberOfPage struct {
	Number            int `json:"number"`
	Size              int `json:"size"`
	TotalPageCount    int `json:"totalPageCount"`
	TotalElementCount int `json:"totalElementCount"`
}

type ResponseBody struct {
	Code       string      `json:"code"`
	Detail     string      `json:"detail"`
	Message    string      `json:"message"`
	Properties interface{} `json:"properties"`
}

type ErrorMsg struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
