// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package httputil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// WriteJSON writes a JSON structure to a response and sets Content-Type header.
func WriteJSON(w http.ResponseWriter, v interface{}) (int, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}

	w.Header().Set("Content-Type", "application/json")

	return w.Write(bytes)
}

func WriteStatus(w http.ResponseWriter, code int) (int, error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	return fmt.Fprintln(w, http.StatusText(code))
}

// ReadHTTPResponseBodyNoErr reads entire HTTP response body into a byte slice, silently skipping errors.
func ReadHTTPResponseBodyNoErr(resp *http.Response) []byte {
	r, _ := io.ReadAll(resp.Body)
	if r == nil {
		return make([]byte, 0)
	}

	return r
}
