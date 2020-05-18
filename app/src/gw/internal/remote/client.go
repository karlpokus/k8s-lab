package remote

import (
	"fmt"
	"bytes"
	"net/http"
)

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

// Body implements io.ReadCloser
type Body struct {
	bytes.Buffer
}

func (b *Body) Close() error {
	return nil
}

func NewBody(s string) *Body {
	body := &Body{}
	body.Write([]byte(s))
	return body
}

type Response map[string]*http.Response

// Mock implements Client
type Mock struct {
	Response
}

func (m *Mock) Do(r *http.Request) (*http.Response, error) {
	path := r.URL.Path
	res, ok := m.Response[path]
	if !ok {
		return &http.Response{}, fmt.Errorf("%s missing in mock client", path)
	}
	// make a response copy so path can be reused
	cp := &http.Response{}
	*cp = *res
	// if request contains body then replace it in the response
	if r.Body != nil {
		cp.Body = r.Body
	}
	return cp, nil
}
