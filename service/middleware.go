// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

// Package service provides base service structures and functions.
package service

// serverSignatureMiddleware adds a service signature to the Server HTTP header.
func serverSignatureMiddleware(s *Service, req *Request, resp *Response) {
	resp.writer.Header().Add("Server", s.name)
}

// loggingMiddleware logs an HTTP request
func loggingMiddleware(s *Service, req *Request, resp *Response) {
	s.log.DebugF("%s %s", req.request.Method, req.request.RequestURI)
}
