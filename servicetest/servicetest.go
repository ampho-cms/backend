// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package servicetest

import (
	"net/http"
	"net/http/httptest"

	"ampho.xyz/ampho/service"
)

type Tester struct {
	svc *service.Service
}


func (w *Tester) DoRequest(method, target string) *http.Response {
	req := httptest.NewRequest(method, target, nil)
	respW := httptest.NewRecorder()
	w.svc.Server().Handler.ServeHTTP(respW, req)

	return respW.Result()
}

func (w *Tester) Stop() {
	w.svc.Stop()
}

func NewTester(name string, setup func(*service.Service) *service.Service) *Tester {
	svc := setup(service.NewTesting(name))
	go svc.Start()

	return &Tester{svc}
}
