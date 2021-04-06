package service

// RequestHandler is the service HTTP request function
type RequestHandler func(*Service, *Request, *Response)
