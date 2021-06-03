// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package util

import (
	"io"
	"net/http"
)

// ReadHTTPResponseBodyNoErr reads entire HTTP response body into a byte slice, silently skipping errors.
func ReadHTTPResponseBodyNoErr(resp *http.Response) []byte {
	r, _ := io.ReadAll(resp.Body)
	if r == nil {
		return make([]byte, 0)
	}

	return r
}
