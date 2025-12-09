package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/go-querystring/query"
	"github.com/sony/gobreaker"

	"swallow-supplier/config"
	customContext "swallow-supplier/context"
)

const (
	// TimestampFormat Timestamp Format
	TimestampFormat = "20060102150405"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Request generic callback interface
type Request struct {
	Ctx           context.Context
	Client        HTTPClient
	Cert          Certificate
	Auth          Auth
	Headers       map[string]string
	Logger        log.Logger
	RequestID     string
	Method        string
	URL           string
	ContentType   string
	Data          []byte
	CustomRequest CustomRequest
}

// Auth authentication object
type Auth struct {
	User string
	Pass string
}

// CustomRequest some custom request to handle
type CustomRequest struct {
	RequestTimeout string
	Retries        int
}

type CircuitBreakerManager struct {
	breakers map[string]*gobreaker.CircuitBreaker
	mutex    sync.RWMutex
}

var breakerManager *CircuitBreakerManager
var once sync.Once

const (
	// RequestTimeout default connection timeout
	RequestTimeout = "120s"
)

const (
	// ContentType content-type header
	ContentType = "Content-Type"

	// ContentTypeFormURLEncoding x-www-form-url-encoded type
	ContentTypeFormURLEncoding = "application/x-www-form-urlencoded; param=value"

	// ContentTypeJSON json type
	ContentTypeJSON = "application/json"

	// ContentTypeXML xml type
	ContentTypeXML = "text/xml"

	// ContentTypeFormData form-data type
	ContentTypeFormData = "multipart/form-data"
)

// NewRequest create new request instance
func NewRequest(cr CustomRequest) *Request {
	request := &Request{}

	request.Headers = make(map[string]string)

	var timeout, _ = time.ParseDuration(RequestTimeout)
	if cr.RequestTimeout != "" {
		// example format of customRequest.RequestTimeout
		// 500ms
		// 30s
		// 60m
		// 2h

		timeout, _ = time.ParseDuration(cr.RequestTimeout)
	}

	request.CustomRequest = cr

	request.Client = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        5000,
			MaxIdleConnsPerHost: 1000,
			MaxConnsPerHost:     1000,
			IdleConnTimeout:     90 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				},
				PreferServerCipherSuites: true,
				MinVersion:               tls.VersionTLS12,
				MaxVersion:               tls.VersionTLS13,
			},
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}

	request.Logger = log.NewJSONLogger(os.Stdout)
	request.Logger = log.With(request.Logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	return request
}

// NewHTTPSRequest create a new request instance with certificate
func NewHTTPSRequest(cr CustomRequest, certFile string, keyFile string) *Request {
	request := NewRequest(cr)
	request.SetCertificate(certFile, keyFile)
	return request
}

// Get send get request to client
func Get(ctx context.Context, cr CustomRequest, serviceName string, url string) (*Response, error) {
	r := NewRequest(cr)
	return r.Get(ctx, serviceName, url)
}

// Post send post request to client
func Post(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	r := NewRequest(cr)
	return r.Post(ctx, serviceName, url, contentType, payload)
}

// Put send put request to client
func Put(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	r := NewRequest(cr)
	return r.Put(ctx, serviceName, url, contentType, payload)
}

// Patch send patch request to client
func Patch(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	r := NewRequest(cr)
	return r.Patch(ctx, serviceName, url, contentType, payload)
}

// Delete send delete request to client
func Delete(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	r := NewRequest(cr)
	return r.Delete(ctx, serviceName, url, contentType, payload)
}

// HTTPSGet send get request to client with certificate
func HTTPSGet(ctx context.Context, cr CustomRequest, serviceName string, url string, certFile string, keyFile string) (*Response, error) {
	r := NewHTTPSRequest(cr, certFile, keyFile)
	return r.Get(ctx, serviceName, url)
}

// HTTPSPost send post request to client with certificate
func HTTPSPost(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}, certFile string, keyFile string) (*Response, error) {
	r := NewHTTPSRequest(cr, certFile, keyFile)
	return r.Post(ctx, serviceName, url, contentType, payload)
}

// HTTPSPut send put request to client with certificate
func HTTPSPut(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}, certFile string, keyFile string) (*Response, error) {
	r := NewHTTPSRequest(cr, certFile, keyFile)
	return r.Put(ctx, serviceName, url, contentType, payload)
}

// HTTPSPatch send patch request to client with certificate
func HTTPSPatch(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}, certFile string, keyFile string) (*Response, error) {
	r := NewHTTPSRequest(cr, certFile, keyFile)
	return r.Patch(ctx, serviceName, url, contentType, payload)
}

// HTTPSDelete send delete request to client with certificate
func HTTPSDelete(ctx context.Context, cr CustomRequest, serviceName string, url string, contentType string, payload interface{}, certFile string, keyFile string) (*Response, error) {
	r := NewHTTPSRequest(cr, certFile, keyFile)
	return r.Delete(ctx, serviceName, url, contentType, payload)
}

// DialFunc ...
func DialFunc(network string, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Duration(time.Second*3))
}

// SetCertificate set cert and key files
func (r *Request) SetCertificate(certFile string, keyFile string) {
	r.Cert = Certificate{
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}

// HasCertificate checking if cert is not empty
func (r *Request) HasCertificate() bool {
	return (Certificate{}) != r.Cert
}

// AddHeader add new header
func (r *Request) AddHeader(key string, value string) {
	r.Headers[key] = value
}

// GetHeaders get all additional header
func (r *Request) GetHeaders() map[string]string {
	return r.Headers
}

// EncodeBasedOnContentType Encode convert the payload based on specified encoding type
func (r *Request) EncodeBasedOnContentType(data interface{}, contentType string) (string, string) {
	ct := strings.ToLower(contentType)

	switch ct {
	case ContentTypeFormURLEncoding:
		switch v := data.(type) {
		case map[string]interface{}:
			values := url.Values{}
			for k, val := range v {
				if val.(string) != "" {
					values.Set(k, val.(string))
				}
			}
			return values.Encode(), ContentTypeFormURLEncoding
		default:
			newData, _ := query.Values(v)
			return newData.Encode(), ContentTypeFormURLEncoding
		}

	case ContentTypeJSON:
		if data == nil || data == "" {
			return "{}", ContentTypeJSON
		}
		newData, _ := json.Marshal(data)
		return string(newData), ContentTypeJSON

	case ContentTypeXML:
		var b bytes.Buffer
		xml.NewEncoder(&b).Encode(data)
		return b.String(), ContentTypeXML
	}

	if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		newData, _ := query.Values(data)
		return newData.Encode(), ContentTypeFormURLEncoding
	}

	return "", contentType
}

// Post send post request to client
func (r *Request) Post(ctx context.Context, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	return r.Send(ctx, serviceName, url, http.MethodPost, contentType, payload)
}

// Get send get request to client
func (r *Request) Get(ctx context.Context, serviceName string, url string) (*Response, error) {
	return r.Send(ctx, serviceName, url, http.MethodGet, "", "")
}

// Put send put request to client
func (r *Request) Put(ctx context.Context, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	return r.Send(ctx, serviceName, url, http.MethodPut, contentType, payload)
}

// Delete send delete request to clint
func (r *Request) Delete(ctx context.Context, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	return r.Send(ctx, serviceName, url, http.MethodDelete, contentType, payload)
}

// Patch send patch request to client
func (r *Request) Patch(ctx context.Context, serviceName string, url string, contentType string, payload interface{}) (*Response, error) {
	return r.Send(ctx, serviceName, url, http.MethodPatch, contentType, payload)
}

// TimeMs Returns Unix time in milliseconds for benchmarking Svc performance.
func TimeMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// Send http request to the given details
func (r *Request) Send(ctx context.Context, serviceName string, url string, method string, contentType string, payload interface{}) (resp *Response, err error) {
	var (
		startTime        int64
		data             []byte
		dataString       string
		validContentType string
	)

	logger := log.With(r.Logger, "method", "request.Send")
	logger = log.With(r.Logger,
		"method", "request.Send",
		customContext.CtxLabelTraceID, customContext.CtxTraceID(ctx),
		customContext.CtxLabelRequestID, customContext.GetCtxHeader(ctx, customContext.CtxLabelRequestID),
		customContext.CtxLabelChannelCode, customContext.GetCtxHeader(ctx, customContext.CtxLabelChannelCode),
	)

	defer func() {
		if recovery := recover(); recovery != nil {
			logger.Log(
				"Request ID", r.RequestID,
				"URL", url,
				"Method", method,
				"Recover from panic", recovery,
				"Stack Trace", string(debug.Stack()),
				"Elapse Time", fmt.Sprintf("%d ms", TimeMs()-startTime),
				"End", time.Now().Format(TimestampFormat),
			)
		} else {
			if err != nil {
				logger.Log(
					"Request ID", r.RequestID,
					"URL", url,
					"Method", method,
					"Send Request Error", err.Error(),
					"Elapse Time", fmt.Sprintf("%d ms", TimeMs()-startTime),
					"End", time.Now().Format(TimestampFormat),
				)
			} else {
				if resp != nil {
					if resp.Error != nil {
						logger.Log(
							"Request ID", r.RequestID,
							"URL", url,
							"Method", method,
							"Send Response Error", err.Error(),
							"Elapse Time", fmt.Sprintf("%d ms", TimeMs()-startTime),
							"End", time.Now().Format(TimestampFormat),
						)
					} else {
						if !strings.Contains(url, "download") && !strings.Contains(url, "enumerations") {
							logger.Log(
								"Request ID", r.RequestID,
								"URL", url,
								"Method", method,
								//"Send Response", strings.TrimSpace(resp.GetAsString()),
								"Elapse Time", fmt.Sprintf("%d ms", TimeMs()-startTime),
								"End", time.Now().Format(TimestampFormat),
							)
						}
					}
				}
			}
		}
	}()

	startTime = TimeMs()

	//var req string
	if strings.Contains(contentType, ContentTypeFormData) {
		validContentType = contentType
		bufferData, _ := payload.(bytes.Buffer)
		data = bufferData.Bytes()
	} else {
		dataString, validContentType = r.EncodeBasedOnContentType(payload, contentType)
		data = []byte(dataString)
		//req = string(data)
	}

	logger.Log(
		"Start", time.Now().Format(TimestampFormat),
		"Request ID", r.RequestID,
		"URL", url,
		"Method", method,
		"ContentType", contentType,
		"Headers", r.GetHeaders(),
		"req body length", len(data),
		//"Request ", req,
	)

	if r.Client == nil {
		logger.Log(
			"Error", "http.Client has not been initialized",
		)
		return nil, errors.New("http.Client has not been initialized")
	}

	if r.HasCertificate() {
		if err = r.InitTLS(); err != nil {
			logger.Log(
				"Initializing TLS Error", err.Error(),
			)
		}
		logger.Log(
			"Certificate", r.Cert.CertFile,
			"Key Certificate", r.Cert.KeyFile,
		)
	}

	r.Ctx = ctx
	r.Method = method
	r.URL = url
	r.ContentType = validContentType
	r.Data = data

	var responseValue interface{}
	if config.Instance().CircuitBreakerEnable == "1" {
		cb := GetBreakerManager().GetBreaker(serviceName, r)
		responseValue, err = cb.Execute(r.SendRequest)
	} else {
		responseValue, err = r.SendRequest()
	}

	if err != nil {
		logger.Log(
			"HTTP Request Error", err.Error(),
		)
		return nil, err
	}

	resp = responseValue.(*Response)

	return resp, nil
}

// SendRequest do the actual http request call
func (r *Request) SendRequest() (resp interface{}, err error) {
	logger := log.With(r.Logger, "method", "request.SendRequest")
	logger = log.With(r.Logger,
		"method", "request.SendRequest",
		customContext.CtxLabelTraceID, customContext.CtxTraceID(r.Ctx),
		customContext.CtxLabelRequestID, customContext.GetCtxHeader(r.Ctx, customContext.CtxLabelRequestID),
		customContext.CtxLabelChannelCode, customContext.GetCtxHeader(r.Ctx, customContext.CtxLabelChannelCode),
	)

	var response *http.Response
	//payload := strings.NewReader(string(r.Data))
	//.Println("777777777777777777777777777777777777 : ", *payload)
	buf := bytes.NewBuffer(nil) // Create a new, empty buffer
	buf.Write(r.Data)           // Write the current data
	/* req, err := http.NewRequest(r.Method, r.URL, buf)
	if err != nil {
		panic(err)
	} */
	request, err := http.NewRequestWithContext(r.Ctx, r.Method, r.URL, buf) //  strings.NewReader(string(r.Data)), bytes.NewBuffer(r.Data)

	for key, val := range r.Headers {
		request.Header.Add(key, val)
	}

	if user, pass := r.GetBasicAuth(); "" != user && "" != pass {
		request.SetBasicAuth(user, pass)
	}

	if r.Method != http.MethodGet {
		contentLength := strconv.Itoa(len(r.Data))

		request.Header.Set("Content-Type", r.ContentType)
		request.Header.Set("Content-Length", contentLength)
	}

	var originalBody []byte
	if request != nil && request.Body != nil {
		originalBody, err = copyBody(request.Body)
		if err != nil {
			return resp, err
		}
		resetBody(request, originalBody)
	}
	var retries = r.CustomRequest.Retries
	for retries > 0 {
		logger.Log(
			"Sending Request Attempt", retries,
		)

		//logger.Log("***** request body  sent *********** : ", request.Body)
		response, err = r.Client.Do(request)
		if err != nil {
			logger.Log(
				"Sending HTTP Request Error", err.Error(),
			)

			if request.Body != nil && len(originalBody) > 0 {
				resetBody(request, originalBody)
			}

			retries--
			continue
		}

		break
	}
	if response != nil {
		resp = NewResponse(response)
	}

	return resp, err
}

func GetBreakerManager() *CircuitBreakerManager {
	once.Do(func() {
		breakerManager = &CircuitBreakerManager{
			breakers: make(map[string]*gobreaker.CircuitBreaker),
		}
	})
	return breakerManager
}

func (cbm *CircuitBreakerManager) GetBreaker(serviceName string, r *Request) *gobreaker.CircuitBreaker {
	cbm.mutex.RLock()
	breaker, exists := cbm.breakers[serviceName]
	cbm.mutex.RUnlock()

	if exists {
		return breaker
	}

	cbm.mutex.Lock()
	defer cbm.mutex.Unlock()

	// Double-check to avoid race condition
	if breaker, exists = cbm.breakers[serviceName]; exists {
		return breaker
	}

	breaker = r.FormatCircuitBreakerSettings(serviceName)
	cbm.breakers[serviceName] = breaker
	return breaker
}

// SetBasicAuth sets the basic authentication
func (r *Request) SetBasicAuth(user string, pass string) {
	r.Auth.User = user
	r.Auth.Pass = pass
}

// GetBasicAuth gets basic authentication
func (r *Request) GetBasicAuth() (user string, pass string) {
	return r.Auth.User, r.Auth.Pass
}

// InitTLS load the certificate and set TLS config
func (r *Request) InitTLS() error {
	var cert tls.Certificate
	var err error

	logger := log.With(r.Logger, "method", "request.InitTLS")

	if r.Cert.CertFile == "" || r.Cert.KeyFile == "" {
		return errors.New("No certificate defined")
	}

	if _, err = os.Stat(r.Cert.KeyFile); err != nil {
		return errors.New("Missing key file")
	}

	if _, err = os.Stat(r.Cert.CertFile); err != nil {
		return errors.New("Missing certificate file")
	}

	if cert, err = tls.LoadX509KeyPair(r.Cert.CertFile, r.Cert.KeyFile); err != nil {
		logger.Log(
			"LoadX509KeyPair Error", err.Error(),
		)
		return err
	}

	ssl := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		Rand:               rand.Reader,
	}

	r.Client = &http.Client{
		Transport: &http.Transport{
			Dial:            DialFunc,
			TLSClientConfig: ssl,
		},
	}

	return nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// PrepareFormData formats the required request form-data
func PrepareFormData(createFormDataHeader bool, stringInput interface{}, fileInput interface{}) (b bytes.Buffer, contentType string, err error) {
	stringMap, _ := stringInput.(map[string]string)
	fileMap, _ := fileInput.(map[string]*os.File)

	w := multipart.NewWriter(&b)
	defer w.Close()

	for index, file := range fileMap {
		var fw io.Writer
		fileInfo, _ := file.Stat()

		if createFormDataHeader {
			buffer := make([]byte, 512)
			if _, err := file.Read(buffer); err != nil {
				return b, "", err
			}
			if _, err := file.Seek(0, 0); err != nil {
				return b, "", err
			}
			var filename = fileInfo.Name()

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition",
				fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
					escapeQuotes(index), escapeQuotes(filename)))
			h.Set("Content-Type", http.DetectContentType(buffer))
			fw, err = w.CreatePart(h)
		} else {
			fw, err = w.CreateFormFile(index, fileInfo.Name())
		}

		_, err = io.Copy(fw, file)
		file.Close()
	}

	for key, str := range stringMap {
		var fw io.Writer

		fw, err = w.CreateFormField(key)
		_, err = io.Copy(fw, strings.NewReader(str))
	}

	contentType = w.FormDataContentType()

	return
}

// FormatCircuitBreakerSettings formats settings required for circuit breaker
func (r *Request) FormatCircuitBreakerSettings(name string) *gobreaker.CircuitBreaker {
	logger := log.With(r.Logger, "method", "request.FormatCircuitBreakerSettings")
	logger = log.With(r.Logger,
		"method", "request.FormatCircuitBreakerSettings",
		customContext.CtxLabelTraceID, customContext.CtxTraceID(r.Ctx),
		customContext.CtxLabelRequestID, customContext.GetCtxHeader(r.Ctx, customContext.CtxLabelRequestID),
		customContext.CtxLabelChannelCode, customContext.GetCtxHeader(r.Ctx, customContext.CtxLabelChannelCode),
	)

	var (
		settings gobreaker.Settings

		defaultRequests             = 10
		defaultFailureRatio float64 = 60
	)

	if config.Instance().CircuitBreakerRequests != "" {
		defaultRequests, _ = strconv.Atoi(config.Instance().CircuitBreakerRequests)
	}

	if config.Instance().CircuitBreakerFailureRatio != "" {
		defaultFailureRatio, _ = strconv.ParseFloat(config.Instance().CircuitBreakerFailureRatio, 64)
	}

	defaultFailureRatio = (defaultFailureRatio / 100)

	settings.Name = name
	// When to flush the internal counts in the Closed state
	settings.Interval = 5 * time.Second
	// Describes how often we should recheck the service health and switch to Half-Open
	settings.Timeout = 3 * time.Second
	// Checks when to switch from Closed to Open
	settings.ReadyToTrip = func(counts gobreaker.Counts) bool {
		// circuit breaker will trip when 60% of requests failed
		// and at least the number of configured requests were made
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		tobeOpened := counts.Requests >= uint32(defaultRequests) && failureRatio >= defaultFailureRatio

		logger.Log(
			"Count of Requests Received", counts.Requests,
			"Default Number of Requests", defaultRequests,
			"Count of Total Failures", counts.TotalFailures,
			"Default Number of Failure", defaultFailureRatio,
			"Calculated Failure Ratio", failureRatio,
			"Is it to be opened", tobeOpened,
		)

		return tobeOpened
	}

	// Handler for every state change
	settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		// for debugging purposes only
		logger.Log(
			"State changed from", from.String(),
			"State changed to", to.String(),
		)
	}

	var cb = gobreaker.NewCircuitBreaker(settings)
	return cb
}

func copyBody(src io.ReadCloser) ([]byte, error) {
	b, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	src.Close()
	return b, nil
}

func resetBody(request *http.Request, originalBody []byte) {
	request.Body = io.NopCloser(bytes.NewBuffer(originalBody))
	request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewBuffer(originalBody)), nil
	}
}
