package service

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Request is the service HTTP request.
type Request struct {
	request *http.Request
}

// Request returns an underlying http.Request object.
func (r *Request) Request() *http.Request {
	return r.request
}

// Vars returns request variables.
func (r *Request) Vars() map[string]string {
	return mux.Vars(r.request)
}

// Var returns a request variable value.
func (r *Request) Var(k string) string {
	return r.Vars()[k]
}
