// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import (
	"encoding/json"
	"net/http"
)

// Response is an HTTP response
type Response struct {
	writer http.ResponseWriter
}

// Writer returns an underlying writer
func (r *Response) Writer() http.ResponseWriter {
	return r.writer
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
