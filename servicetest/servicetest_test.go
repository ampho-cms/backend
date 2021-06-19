// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package servicetest_test

import (
	"testing"

	"ampho.xyz/ampho/routing"
	"ampho.xyz/ampho/service"
	"ampho.xyz/ampho/servicetest"
	"ampho.xyz/ampho/util"
)

func TestDoRequest(t *testing.T) {
	// Create a service instance
	svc, err := service.NewTesting("hello")
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Map a handler to a path
	svc.Router().Handle("/hello/{name}", func(req *routing.Request, resp *routing.Response) {
		_, _ = resp.WriteJSON(map[string]interface{}{
			"name": req.Var("name"),
		})
	})

	resp := servicetest.DoRequest(svc, "GET", "/hello/world")
	if resp.StatusCode != 200 {
		t.Errorf("bad status")
	}

	body := string(util.ReadHTTPResponseBodyNoErr(resp))
	if body != "{\"name\":\"world\"}" {
		t.Errorf("bad response body: %s", body)
	}
}
