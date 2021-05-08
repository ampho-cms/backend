// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package service

import "ampho/routing"

// ServerSignatureMiddleware adds a Server HTTP header to server responses.
func (s *Service) ServerSignatureMiddleware(req *routing.Request, resp *routing.Response) bool {
	resp.Writer().Header().Add("Server", s.Signature())
	return true
}

// RequestLogDebugMiddleware logs requests at debug level.
func (s *Service) RequestLogDebugMiddleware(req *routing.Request, resp *routing.Response) bool {
	s.log.DebugF("%s %s", req.Request().Method, req.Request().RequestURI)
	return true
}
