// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package routing

import (
	"encoding/json"
	"net/http"
)

// Response is the http.ResponseWriter wrapper with several helper methods.
type Response struct {
	writer http.ResponseWriter
}

// Writer returns an underlying writer
func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

// SetStatus writes a status string to the response header.
func (r *Response) SetStatus(code int) {
	r.writer.WriteHeader(code)
}

// SetHeader sets header's value.
func (r *Response) SetHeader(key, value string) {
	r.writer.Header().Set(key, value)
}

// WriteString writes a string to a response.
func (r *Response) WriteString(v string) (int, error) {
	return r.writer.Write([]byte(v))
}

// WriteJSON writes a JSON structure to a response and sets Content-Type header.
func (r *Response) WriteJSON(v interface{}) (int, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	r.writer.Header().Set("Content-Type", "application/json")

	return r.writer.Write(bytes)
}
