package wraptest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanderhahn/htpasswd/wrap"
)

// MockResponseWriter enables testing api responses
type MockResponseWriter struct {
	Output      interface{}
	ErrorOutput error
	ErrorCode   int
}

// JSON remembers the output response
func (w *MockResponseWriter) JSON(data interface{}) (int, error) {
	w.Output = data
	return 0, nil
}

// JSON remembers the error response
func (w *MockResponseWriter) Error(err error, code int) (int, error) {
	w.ErrorOutput = err
	w.ErrorCode = code
	return 0, nil
}

// MustMarshal enforce json marshalling
func MustMarshal(t *testing.T, data interface{}) []byte {
	buf, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	return buf
}

// NewRequest for testing
func NewRequest(t *testing.T, name string, input interface{}, sessionVariables wrap.SessionVariables) *wrap.Request {
	req := wrap.RequestInput{
		Action: wrap.ActionName{
			Name: name,
		},
		Input:            MustMarshal(t, input),
		SessionVariables: sessionVariables,
	}
	body := bytes.NewReader(MustMarshal(t, req))
	_r := httptest.NewRequest(http.MethodPost, "/api", body)
	r, err := wrap.NewRequest(_r)
	if err != nil {
		t.Fail()
	}
	return r
}
