package goulash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoboticCheese/goulash/universe"
)

func StartHTTP(r interface{}, h interface{}) (ts *httptest.Server) {
	switch resp := r.(type) {
	case string:
		handler := func(w http.ResponseWriter, req *http.Request) {
			switch headers := h.(type) {
			case map[string]string:
				for k, v := range headers {
					w.Header().Set(k, v)
				}
			case func() map[string]string:
				for k, v := range headers() {
					w.Header().Set(k, v)
				}
			}
			fmt.Fprint(w, resp)
		}
		ts = httptest.NewServer(http.HandlerFunc(handler))
	case func() string:
		handler := func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, resp())
		}
		ts = httptest.NewServer(http.HandlerFunc(handler))
	case func(http.ResponseWriter, *http.Request):
		ts = httptest.NewServer(http.HandlerFunc(resp))
	}
	return
}

func TestNewAPIInstance(t *testing.T) {
	ts := StartHTTP("", nil)
	i, err := NewAPIInstance(ts.URL)
	ts.Close()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if i.BaseURL == "" {
		t.Fatalf("Expected non-empty BaseURL, got: %v", i.BaseURL)
	}
}

func TestNewCookbook(t *testing.T) {
	data, _ := json.Marshal(Cookbook{
		Name:        "chef-dk",
		Maintainer:  "roboticcheese",
		Description: "Installs/configures the Chef-DK",
	})
	ts := StartHTTP(string(data), nil)
	defer ts.Close()
	i := APIInstance{
		Component: Component{Endpoint: ts.URL},
		BaseURL:   ts.URL,
	}
	c, err := NewCookbook(&i, "chef-dk")
	ts.Close()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if c.Maintainer != "roboticcheese" {
		t.Fatalf("Expected: roboticcheese, got: %v", c.Maintainer)
	}
}

func TestNewCookbookVersion(t *testing.T) {
	data, _ := json.Marshal(CookbookVersion{
		License: "Apache v2.0",
	})
	ts := StartHTTP(string(data), nil)
	defer ts.Close()
	c := Cookbook{
		Component: Component{Endpoint: ts.URL},
	}
	cv, err := NewCookbookVersion(&c, "1.2.3")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cv.License != "Apache v2.0" {
		t.Fatalf("Expected: Apache v2.0, got: %v", cv.License)
	}
}

func TestNewUniverse(t *testing.T) {
	data, _ := json.Marshal(map[string]map[string]*universe.CookbookVersion{
		"chef": {
			"0.12.0": &universe.CookbookVersion{
				LocationType: "opscode",
				LocationPath: "https://supermarket.chef.io/api/v1",
				DownloadURL:  "https://supermarket.chef.io/api/v1/cookbooks/chef/versions/0.12.0/download",
				Dependencies: map[string]string{"runit": ">= 0.0.0", "couchdb": ">= 0.0.0"},
			},
		},
	})
	ts := StartHTTP(string(data), nil)
	defer ts.Close()
	i := APIInstance{
		Component: Component{Endpoint: ts.URL},
		BaseURL:   ts.URL,
	}
	u, err := NewUniverse(&i)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(u.Cookbooks) != 1 {
		t.Fatalf("Expected 1 Cookbook, got: %v", len(u.Cookbooks))
	}
}
