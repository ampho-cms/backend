// Author:  Alexander Shepetko
// Email:   a@shepetko.com
// License: MIT

package servicetest_test

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"

	"ampho.xyz/config"
	"ampho.xyz/httputil"
	"ampho.xyz/service"
	"ampho.xyz/servicetest"
	"ampho.xyz/util"
)

func TestDoRequest(t *testing.T) {
	// Create a service instance
	svc, err := service.New(config.NewTesting("hello"))
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Map a handler to a path
	svc.Router().HandleFunc("/hello/{name}", func(writer http.ResponseWriter, req *http.Request) {
		_, _ = httputil.WriteJSON(writer, map[string]interface{}{
			"name": mux.Vars(req)["name"],
		})
	})

	req, err := http.NewRequest("GET", "/hello/world", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := servicetest.DoRequest(svc, req)
	if resp.Result().StatusCode != 200 {
		t.Errorf("bad status: %d", resp.Result().StatusCode)
	}

	body := string(util.ReadHTTPResponseBodyNoErr(resp.Result()))
	if body != "{\"name\":\"world\"}" {
		t.Errorf("bad response body: %s", body)
	}
}
