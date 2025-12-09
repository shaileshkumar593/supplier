package travolution

// Response format for success case
type Response struct {
	Code            string      `json:"code"`
	Body            interface{} `json:"body,omitempty"`
	Contents        Content     `json:"contents"`
	HtmlTypeContent *string     `json:"contentType"`
}

type Content struct {
	Highlights  *string `json:"highlights"`
	Description *string `json:"description"`
	SubContent  *string `json:"subContent"`
}

type ErrorMsg struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
