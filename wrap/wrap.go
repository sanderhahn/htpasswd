package wrap

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// ResponseWriter adds json rendering
type ResponseWriter interface {
	JSON(data interface{}) (int, error)
	Error(err error, code int) (int, error)
}

// wrapResponseWriter adds json rendering
type wrapResponseWriter struct {
	http.ResponseWriter
}

// NewResponseWriter constructor
func NewResponseWriter(_w http.ResponseWriter) ResponseWriter {
	return &wrapResponseWriter{ResponseWriter: _w}
}

// JSON renders json data
func (w *wrapResponseWriter) JSON(data interface{}) (int, error) {
	w.Header().Set("content-type", "application/json")
	out, err := json.Marshal(data)
	if err != nil {
		return w.Error(err, http.StatusInternalServerError)
	}
	return w.Write(out)
}

// Error renders an error response for an action
func (w *wrapResponseWriter) Error(err error, code int) (int, error) {
	out, err := json.Marshal(struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}{
		Message: err.Error(),
		Code:    strconv.Itoa(code),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return w.Write([]byte(err.Error()))
	}
	w.WriteHeader(code)
	return w.Write(out)
}

// ActionName wraps name
type ActionName struct {
	Name string `json:"name"`
}

// SessionVariables type
type SessionVariables map[string]string

// RequestInput mimics the action structure
type RequestInput struct {
	SessionVariables SessionVariables `json:"session_variables"`
	Input            json.RawMessage  `json:"input"`
	Action           ActionName       `json:"action"`
}

// Request adds action parsing
type Request struct {
	Payload RequestInput
	*http.Request
}

// NewRequest constructor
func NewRequest(_r *http.Request) (*Request, error) {
	r := &Request{Request: _r}
	buf, err := ioutil.ReadAll(_r.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(buf, &r.Payload); err != nil {
		return nil, err
	}
	if r.HasHeader("x-debug") {
		input, err := json.Marshal(r.Payload.Input)
		if err != nil {
			return nil, err
		}
		session, err := json.Marshal(r.Payload.SessionVariables)
		if err != nil {
			return nil, err
		}
		log.Printf("action %s(%s) session variables %s\n", r.Payload.Action.Name, input, session)
	}
	return r, nil
}

// ParseInput allows parsing the payload input
func (r *Request) ParseInput(v interface{}) error {
	return json.Unmarshal(r.Payload.Input, v)
}

// HasHeader tests if a header input is not empty
func (r *Request) HasHeader(header string) bool {
	return r.Header.Get(header) != ""
}

// ParseHeader allows parsing the a header input
func (r *Request) ParseHeader(header string, v interface{}) error {
	return json.Unmarshal([]byte(r.Header.Get(header)), v)
}
