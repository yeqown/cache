package utils

import (
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// CopyRequest copy request from a http.Request
func CopyRequest(req *http.Request) *http.Request {
	var (
		body       []byte
		readCloser io.ReadCloser
		cpyReq     *http.Request
		err        error
	)

	if req.Body != nil {
		body, _ = ioutil.ReadAll(req.Body)
	}
	readCloser = ioutil.NopCloser(bytes.NewBuffer(body))
	cpyReq, err = http.NewRequest(req.Method, req.URL.String(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	cpyReq.Header = req.Header
	req.Body = readCloser
	return cpyReq
}

// ParseRequestForm parse request and get form form body or url
// nomarlly support "application/www-form-urlencodeded" header
// TODO: ParseMultipartForm with Content-Type multipart/form-data
func ParseRequestForm(cpyReq *http.Request) url.Values {
	cpyReq.ParseForm()

	switch cpyReq.Method {
	case http.MethodPost, http.MethodPut:
		cpyReq.ParseMultipartForm(32 << 20)
		return cpyReq.PostForm
	default:
		return cpyReq.Form
	}
}

// EncodeFormToString ... you must copy http.Request manually
func EncodeFormToString(req *http.Request) string {
	form := ParseRequestForm(req)
	buffer := bytes.NewBufferString(form.Encode())
	return hex.EncodeToString(buffer.Bytes())
}
