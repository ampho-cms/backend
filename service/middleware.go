// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import (
	"ampho.xyz/ampho/logging"
	"ampho.xyz/ampho/routing"
)

// ServerSignatureMiddleware adds a Server HTTP header to server responses.
func (s *Base) ServerSignatureMiddleware(req *routing.Request, resp *routing.Response) bool {
	resp.Writer().Header().Add("Server", s.Signature())
	return true
}

// RequestLogDebugMiddleware logs requests at debug level.
func (s *Base) RequestLogDebugMiddleware(req *routing.Request, resp *routing.Response) bool {
	logging.DebugF("%s %s", req.Request().Method, req.Request().RequestURI)
	return true
}
