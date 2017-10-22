package mocks

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/mock"
)

type EchoContextMock struct {
	mock.Mock
}

func (e *EchoContextMock) Request() *http.Request {
	args := e.Called()
	return args.Get(0).(*http.Request)
}

func (e *EchoContextMock) SetRequest(r *http.Request) {
	e.Called(r)
}

func (e *EchoContextMock) Response() *echo.Response {
	args := e.Called()
	return args.Get(0).(*echo.Response)
}

func (e *EchoContextMock) IsTLS() bool {
	args := e.Called()
	return args.Bool(0)
}

func (e *EchoContextMock) IsWebSocket() bool {
	args := e.Called()
	return args.Bool(0)
}

func (e *EchoContextMock) Scheme() string {
	args := e.Called()
	return args.String(0)
}

func (e *EchoContextMock) RealIP() string {
	args := e.Called()
	return args.String(0)
}

func (e *EchoContextMock) Path() string {
	args := e.Called()
	return args.String(0)
}

func (e *EchoContextMock) SetPath(p string) {
	e.Called(p)
}

func (e *EchoContextMock) Param(name string) string {
	args := e.Called(name)
	return args.String(0)
}

func (e *EchoContextMock) ParamNames() []string {
	args := e.Called()
	return args.Get(0).([]string)
}

func (e *EchoContextMock) SetParamNames(names ...string) {
	a := make([]interface{}, len(names))

	for i, v := range names {
		a[i] = v
	}

	e.Called(a...)
}
func (e *EchoContextMock) ParamValues() []string {
	args := e.Called()
	return args.Get(0).([]string)
}

func (e *EchoContextMock) SetParamValues(values ...string) {
	a := make([]interface{}, len(values))

	for i, v := range values {
		a[i] = v
	}

	e.Called(a...)
}

func (e *EchoContextMock) QueryParam(name string) string {
	args := e.Called(name)
	return args.String(0)
}

func (e *EchoContextMock) QueryParams() url.Values {
	args := e.Called()
	return args.Get(0).(url.Values)
}

func (e *EchoContextMock) QueryString() string {
	args := e.Called()
	return args.String(0)
}

func (e *EchoContextMock) FormValue(name string) string {
	args := e.Called(name)
	return args.String(0)
}

func (e *EchoContextMock) FormParams() (url.Values, error) {
	args := e.Called()
	return args.Get(0).(url.Values), args.Error(1)
}

func (e *EchoContextMock) FormFile(name string) (*multipart.FileHeader, error) {
	args := e.Called(name)
	return args.Get(0).(*multipart.FileHeader), args.Error(1)
}

func (e *EchoContextMock) MultipartForm() (*multipart.Form, error) {
	args := e.Called()
	return args.Get(0).(*multipart.Form), args.Error(1)
}

func (e *EchoContextMock) Cookie(name string) (*http.Cookie, error) {
	args := e.Called(name)
	return args.Get(0).(*http.Cookie), args.Error(1)
}

func (e *EchoContextMock) SetCookie(cookie *http.Cookie) {
	e.Called(cookie)
}

func (e *EchoContextMock) Cookies() []*http.Cookie {
	args := e.Called()
	return args.Get(0).([]*http.Cookie)
}

func (e *EchoContextMock) Get(key string) interface{} {
	args := e.Called(key)
	return args.Get(0)
}

func (e *EchoContextMock) Set(key string, val interface{}) {
	e.Called(key, val)
}

func (e *EchoContextMock) Bind(i interface{}) error {
	args := e.Called(i)
	return args.Error(0)
}

func (e *EchoContextMock) Validate(i interface{}) error {
	args := e.Called(i)
	return args.Error(0)
}

func (e *EchoContextMock) Render(code int, name string, data interface{}) error {
	args := e.Called(code, name, data)
	return args.Error(0)
}

func (e *EchoContextMock) HTML(code int, html string) error {
	args := e.Called(code, html)
	return args.Error(0)
}

func (e *EchoContextMock) HTMLBlob(code int, b []byte) error {
	args := e.Called(code, b)
	return args.Error(0)
}

func (e *EchoContextMock) String(code int, s string) error {
	args := e.Called(code, s)
	return args.Error(0)
}

func (e *EchoContextMock) JSON(code int, i interface{}) error {
	args := e.Called(code, i)
	return args.Error(0)
}

func (e *EchoContextMock) JSONPretty(code int, i interface{}, indent string) error {
	args := e.Called(code, i)
	return args.Error(0)
}

func (e *EchoContextMock) JSONBlob(code int, b []byte) error {
	args := e.Called(code, b)
	return args.Error(0)
}

func (e *EchoContextMock) JSONP(code int, callback string, i interface{}) error {
	args := e.Called(code, callback)
	return args.Error(0)
}

func (e *EchoContextMock) JSONPBlob(code int, callback string, b []byte) error {
	args := e.Called(code, callback, b)
	return args.Error(0)
}

func (e *EchoContextMock) XML(code int, i interface{}) error {
	args := e.Called(code, i)
	return args.Error(0)
}

func (e *EchoContextMock) XMLPretty(code int, i interface{}, indent string) error {
	args := e.Called(code, i, indent)
	return args.Error(0)
}

func (e *EchoContextMock) XMLBlob(code int, b []byte) error {
	args := e.Called(code, b)
	return args.Error(0)
}

func (e *EchoContextMock) Blob(code int, contentType string, b []byte) error {
	args := e.Called(code, contentType, b)
	return args.Error(0)
}

func (e *EchoContextMock) Stream(code int, contentType string, r io.Reader) error {
	args := e.Called(code, contentType, r)
	return args.Error(0)
}

func (e *EchoContextMock) File(file string) error {
	args := e.Called(file)
	return args.Error(0)
}

func (e *EchoContextMock) Attachment(file string, name string) error {
	args := e.Called(file, name)
	return args.Error(0)
}

func (e *EchoContextMock) Inline(file string, name string) error {
	args := e.Called(file, name)
	return args.Error(0)
}

func (e *EchoContextMock) NoContent(code int) error {
	args := e.Called(code)
	return args.Error(0)
}

func (e *EchoContextMock) Redirect(code int, url string) error {
	args := e.Called(code, url)
	return args.Error(0)
}

func (e *EchoContextMock) Error(err error) {
	e.Called(err)
}

func (e *EchoContextMock) Handler() echo.HandlerFunc {
	args := e.Called()
	return args.Get(0).(echo.HandlerFunc)
}

func (e *EchoContextMock) SetHandler(h echo.HandlerFunc) {
	e.Called(h)
}

func (e *EchoContextMock) Logger() echo.Logger {
	args := e.Called()
	return args.Get(0).(echo.Logger)

}

func (e *EchoContextMock) Echo() *echo.Echo {
	args := e.Called()
	return args.Get(0).(*echo.Echo)
}

func (e *EchoContextMock) Reset(r *http.Request, w http.ResponseWriter) {
	e.Called(r, w)
}
