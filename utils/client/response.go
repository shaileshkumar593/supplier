package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Response response callback
type Response struct {
	Headers http.Header
	Data    interface{}
	Body    string
	Status  int
	Message string
	Error   error
}

// NewResponse create new response object
func NewResponse(response *http.Response) *Response {
	var body []byte
	var err error

	r := &Response{}

	// always close the response-body, even if content is not required
	defer response.Body.Close()

	r.Status = response.StatusCode
	r.Headers = response.Header

	if response == nil {
		r.Error = errors.New("http response is nil")
		return r
	}

	if body, err = io.ReadAll(response.Body); err != nil {
		r.Error = err
		return r
	}

	r.Data = body
	r.Body = string(body)

	return r
}

// HasError check if has error
func (r *Response) HasError() bool {
	return r.Error != nil
}

// GetError get error object
func (r *Response) GetError() error {
	return r.Error
}

// GetData get response data
func (r *Response) GetData() interface{} {
	return r.Data
}

// GetAsString get response data
func (r *Response) GetAsString() string {
	return r.Body
}

// GetAsJSON get response as json
func (r *Response) GetAsJSON(v interface{}) error {
	return json.Unmarshal([]byte(r.Body), v)
}
