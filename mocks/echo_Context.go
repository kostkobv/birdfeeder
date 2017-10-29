package mocks

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

// EchoContextMock is echo.Context mock
type EchoContextMock struct {
	mock.Mock
}

// Request mock
func (e *EchoContextMock) Request() *http.Request {
	args := e.Called()
	return args.Get(0).(*http.Request)
}

// SetRequest mock
func (e *EchoContextMock) SetRequest(r *http.Request) {
	e.Called(r)
}

// Response mock
func (e *EchoContextMock) Response() *echo.Response {
	args := e.Called()
	return args.Get(0).(*echo.Response)
}

// IsTLS mock
func (e *EchoContextMock) IsTLS() bool {
	args := e.Called()
	return args.Bool(0)
}

// IsWebSocket mock
func (e *EchoContextMock) IsWebSocket() bool {
	args := e.Called()
	return args.Bool(0)
}

// Scheme mock
func (e *EchoContextMock) Scheme() string {
	args := e.Called()
	return args.String(0)
}

// RealIP mock
func (e *EchoContextMock) RealIP() string {
	args := e.Called()
	return args.String(0)
}

// Path mock
func (e *EchoContextMock) Path() string {
	args := e.Called()
	return args.String(0)
}

// SetPath mock
func (e *EchoContextMock) SetPath(p string) {
	e.Called(p)
}

// Param mock
func (e *EchoContextMock) Param(name string) string {
	args := e.Called(name)
	return args.String(0)
}

// ParamNames mock
func (e *EchoContextMock) ParamNames() []string {
	args := e.Called()
	return args.Get(0).([]string)
}

// SetParamNames mock
func (e *EchoContextMock) SetParamNames(names ...string) {
	a := make([]interface{}, len(names))

	for i, v := range names {
		a[i] = v
	}

	e.Called(a...)
}

// ParamValues mock
func (e *EchoContextMock) ParamValues() []string {
	args := e.Called()
	return args.Get(0).([]string)
}

// SetParamValues mock
func (e *EchoContextMock) SetParamValues(values ...string) {
	a := make([]interface{}, len(values))

	for i, v := range values {
		a[i] = v
	}

	e.Called(a...)
}

// QueryParam mock
func (e *EchoContextMock) QueryParam(name string) string {
	args := e.Called(name)
	return args.String(0)
}

// QueryParams mock
func (e *EchoContextMock) QueryParams() url.Values {
	args := e.Called()
	return args.Get(0).(url.Values)
}

// QueryString mock
func (e *EchoContextMock) QueryString() string {
	args := e.Called()
	return args.String(0)
}

// FormValue mock
func (e *EchoContextMock) FormValue(name string) string {
	args := e.Called(name)
	return args.String(0)
}

// FormParams mock
func (e *EchoContextMock) FormParams() (url.Values, error) {
	args := e.Called()
	return args.Get(0).(url.Values), args.Error(1)
}

// FormFile mock
func (e *EchoContextMock) FormFile(name string) (*multipart.FileHeader, error) {
	args := e.Called(name)
	return args.Get(0).(*multipart.FileHeader), args.Error(1)
}

// MultipartForm mock
func (e *EchoContextMock) MultipartForm() (*multipart.Form, error) {
	args := e.Called()
	return args.Get(0).(*multipart.Form), args.Error(1)
}

// Cookie mock
func (e *EchoContextMock) Cookie(name string) (*http.Cookie, error) {
	args := e.Called(name)
	return args.Get(0).(*http.Cookie), args.Error(1)
}

// SetCookie mock
func (e *EchoContextMock) SetCookie(cookie *http.Cookie) {
	e.Called(cookie)
}

// Cookies mock
func (e *EchoContextMock) Cookies() []*http.Cookie {
	args := e.Called()
	return args.Get(0).([]*http.Cookie)
}

// Get mock
func (e *EchoContextMock) Get(key string) interface{} {
	args := e.Called(key)
	return args.Get(0)
}

// Set mock
func (e *EchoContextMock) Set(key string, val interface{}) {
	e.Called(key, val)
}

// Bind mock
func (e *EchoContextMock) Bind(i interface{}) error {
	args := e.Called(i)
	return args.Error(0)
}

// Validate mock
func (e *EchoContextMock) Validate(i interface{}) error {
	args := e.Called(i)
	return args.Error(0)
}

// Render mock
func (e *EchoContextMock) Render(code int, name string, data interface{}) error {
	args := e.Called(code, name, data)
	return args.Error(0)
}

// HTML mock
func (e *EchoContextMock) HTML(code int, html string) error {
	args := e.Called(code, html)
	return args.Error(0)
}

// HTMLBlob mock
func (e *EchoContextMock) HTMLBlob(code int, b []byte) error {
	args := e.Called(code, b)
	return args.Error(0)
}

// String mock
func (e *EchoContextMock) String(code int, s string) error {
	args := e.Called(code, s)
	return args.Error(0)
}

// JSON mock
func (e *EchoContextMock) JSON(code int, i interface{}) error {
	args := e.Called(code, i)
	return args.Error(0)
}

// JSONPretty mock
func (e *EchoContextMock) JSONPretty(code int, i interface{}, indent string) error {
	args := e.Called(code, i)
	return args.Error(0)
}

// JSONBlob mock
func (e *EchoContextMock) JSONBlob(code int, b []byte) error {
	args := e.Called(code, b)
	return args.Error(0)
}

// JSONP mock
func (e *EchoContextMock) JSONP(code int, callback string, i interface{}) error {
	args := e.Called(code, callback)
	return args.Error(0)
}

// JSONPBlob mock
func (e *EchoContextMock) JSONPBlob(code int, callback string, b []byte) error {
	args := e.Called(code, callback, b)
	return args.Error(0)
}

// XML mock
func (e *EchoContextMock) XML(code int, i interface{}) error {
	args := e.Called(code, i)
	return args.Error(0)
}

// XMLPretty mock
func (e *EchoContextMock) XMLPretty(code int, i interface{}, indent string) error {
	args := e.Called(code, i, indent)
	return args.Error(0)
}

// XMLBlob mock
func (e *EchoContextMock) XMLBlob(code int, b []byte) error {
	args := e.Called(code, b)
	return args.Error(0)
}

// Blob mock
func (e *EchoContextMock) Blob(code int, contentType string, b []byte) error {
	args := e.Called(code, contentType, b)
	return args.Error(0)
}

// Stream mock
func (e *EchoContextMock) Stream(code int, contentType string, r io.Reader) error {
	args := e.Called(code, contentType, r)
	return args.Error(0)
}

// File mock
func (e *EchoContextMock) File(file string) error {
	args := e.Called(file)
	return args.Error(0)
}

// Attachment mock
func (e *EchoContextMock) Attachment(file string, name string) error {
	args := e.Called(file, name)
	return args.Error(0)
}

// Inline mock
func (e *EchoContextMock) Inline(file string, name string) error {
	args := e.Called(file, name)
	return args.Error(0)
}

// NoContent mock
func (e *EchoContextMock) NoContent(code int) error {
	args := e.Called(code)
	return args.Error(0)
}

// Redirect mock
func (e *EchoContextMock) Redirect(code int, url string) error {
	args := e.Called(code, url)
	return args.Error(0)
}

// Error mock
func (e *EchoContextMock) Error(err error) {
	e.Called(err)
}

// Handler mock
func (e *EchoContextMock) Handler() echo.HandlerFunc {
	args := e.Called()
	return args.Get(0).(echo.HandlerFunc)
}

// SetHandler mock
func (e *EchoContextMock) SetHandler(h echo.HandlerFunc) {
	e.Called(h)
}

// Logger mock
func (e *EchoContextMock) Logger() echo.Logger {
	args := e.Called()
	return args.Get(0).(echo.Logger)

}

// Echo mock
func (e *EchoContextMock) Echo() *echo.Echo {
	args := e.Called()
	return args.Get(0).(*echo.Echo)
}

// Reset mock
func (e *EchoContextMock) Reset(r *http.Request, w http.ResponseWriter) {
	e.Called(r, w)
}
