// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service provides base service structures and functions.
package service

import "ampho/routing"

// serverSignatureMiddleware adds a service signature to the Server HTTP header.
func (s *Service) serverSignatureMiddleware(req *routing.Request, resp *routing.Response) {
	resp.Writer().Header().Add("Server", s.name)
}

// loggingMiddleware logs an HTTP request
func (s *Service) loggingMiddleware(req *routing.Request, resp *routing.Response) {
	s.log.DebugF("%s %s", req.Request().Method, req.Request().RequestURI)
}
