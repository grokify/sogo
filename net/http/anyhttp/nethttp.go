package anyhttp

import (
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/grokify/mogo/net/http/httputilmore"
)

type RequestNetHTTP struct {
	Raw      *http.Request
	allArgs  *ArgsURLValues
	postArgs *ArgsURLValues
	// multipartForm       *multipart.Form
	parsedMultipartForm bool
	parsedFormArgs      bool
}

func NewRequestNetHTTP(req *http.Request) *RequestNetHTTP {
	return &RequestNetHTTP{
		Raw:      req,
		allArgs:  &ArgsURLValues{Raw: req.Form},
		postArgs: &ArgsURLValues{Raw: req.PostForm}}
}

func (r *RequestNetHTTP) ParseForm() error {
	if r.parsedFormArgs {
		return nil
	}
	r.parsedFormArgs = true
	if err := r.Raw.ParseForm(); err != nil {
		return err
	}
	r.allArgs = &ArgsURLValues{r.Raw.Form}
	r.postArgs = &ArgsURLValues{r.Raw.PostForm}
	return nil
}

func (r RequestNetHTTP) Header(s string) []byte       { return []byte(r.Raw.Header.Get(s)) }
func (r RequestNetHTTP) HeaderString(s string) string { return r.Raw.Header.Get(s) }
func (r RequestNetHTTP) RemoteAddr() net.Addr {
	return Addr{Protocol: "tcp", Address: r.Raw.RemoteAddr}
}
func (r RequestNetHTTP) RemoteAddress() string     { return r.Raw.RemoteAddr }
func (r RequestNetHTTP) UserAgent() []byte         { return []byte(r.Raw.UserAgent()) }
func (r RequestNetHTTP) AllArgs() Args             { return r.allArgs }
func (r RequestNetHTTP) QueryArgs() Args           { return &ArgsURLValues{r.Raw.URL.Query()} }
func (r RequestNetHTTP) PostArgs() Args            { return r.postArgs }
func (r RequestNetHTTP) Method() []byte            { return []byte(r.Raw.Method) }
func (r RequestNetHTTP) Headers() http.Header      { return r.Raw.Header }
func (r RequestNetHTTP) Form() url.Values          { return r.Raw.Form }
func (r RequestNetHTTP) RequestURI() []byte        { return []byte(r.Raw.RequestURI) }
func (r RequestNetHTTP) PostBody() ([]byte, error) { return io.ReadAll(r.Raw.Body) }

func (r RequestNetHTTP) RequestURIVar(s string) string {
	if r.Raw == nil {
		return ""
	}
	vars := mux.Vars(r.Raw)
	if val, ok := vars[s]; ok {
		return val
	}
	return ""
}

func (r *RequestNetHTTP) MultipartForm() (*multipart.Form, error) {
	if !r.parsedMultipartForm {
		r.parsedMultipartForm = true
		if err := r.Raw.ParseMultipartForm(100000); err != nil {
			return nil, err
		}
	}
	return r.Raw.MultipartForm, nil
}

type ResponseNetHTTP struct{ Raw http.ResponseWriter }

func NewResponseNetHTTP(w http.ResponseWriter) ResponseNetHTTP { return ResponseNetHTTP{Raw: w} }

func (w ResponseNetHTTP) GetHeader(k string) []byte { return []byte(w.Raw.Header().Get(k)) }
func (w ResponseNetHTTP) SetHeader(k, v string)     { w.Raw.Header().Set(k, v) }
func (w ResponseNetHTTP) SetStatusCode(code int)    { w.Raw.WriteHeader(code) }
func (w ResponseNetHTTP) SetContentType(ct string) {
	w.Raw.Header().Set(httputilmore.HeaderContentType, ct)
}

func (w ResponseNetHTTP) SetBodyBytes(body []byte) (int, error) {
	return w.Raw.Write(body)
}

// SetBodyStream takes an `io.Reader`. `bodySize` is accepted but
// ignored to fulfill the `Response` interface requirement.
func (w ResponseNetHTTP) SetBodyStream(bodyStream io.Reader, bodySize int) error {
	bytes, err := io.ReadAll(bodyStream)
	if err != nil {
		return err
	}
	_, err = w.Raw.Write(bytes)
	return err
}

func (w ResponseNetHTTP) SetCookie(cookie *Cookie) {
	c := cookie.ToNetHTTP()
	c.Secure = true
	http.SetCookie(w.Raw, c)
}

func NewResReqNetHTTP(res http.ResponseWriter, req *http.Request) (ResponseNetHTTP, *RequestNetHTTP) {
	return NewResponseNetHTTP(res), NewRequestNetHTTP(req)
}
