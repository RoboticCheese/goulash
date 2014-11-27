package goulash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoboticCheese/goulash/apiinstance"
	"github.com/RoboticCheese/goulash/component"
	"github.com/RoboticCheese/goulash/cookbook"
	"github.com/RoboticCheese/goulash/cookbookversion"
	"github.com/RoboticCheese/goulash/universe"
)

func startHTTP(data interface{}) (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				j, _ := json.Marshal(data)
				fmt.Fprint(w, string(j))
			},
		),
	)
	return
}

func Test_NewAPIInstance_1(t *testing.T) {
	data := map[string]string{"doesnot": "matter"}
	ts := startHTTP(data)
	i, err := NewAPIInstance(ts.URL)
	ts.Close()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if i.BaseURL == "" {
		t.Fatalf("Expected non-empty BaseURL, got: %v", i.BaseURL)
	}
}

func Test_NewCookbook_1(t *testing.T) {
	data := cookbook.Cookbook{
		Name:        "chef-dk",
		Maintainer:  "roboticcheese",
		Description: "Installs/configures the Chef-DK",
	}
	ts := startHTTP(data)
	defer ts.Close()
	i := apiinstance.APIInstance{
		Component: component.Component{Endpoint: ts.URL},
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

func Test_NewCookbookVersion_1(t *testing.T) {
	data := cookbookversion.CookbookVersion{
		License: "Apache v2.0",
	}
	ts := startHTTP(data)
	defer ts.Close()
	c := cookbook.Cookbook{
		Component: component.Component{Endpoint: ts.URL},
	}
	cv, err := NewCookbookVersion(&c, "1.2.3")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cv.License != "Apache v2.0" {
		t.Fatalf("Expected: Apache v2.0, got: %v", cv.License)
	}
}

func Test_NewUniverse_1(t *testing.T) {
	data := map[string]map[string]*universe.CookbookVersion{
		"chef": {
			"0.12.0": &universe.CookbookVersion{
				LocationType: "opscode",
				LocationPath: "https://supermarket.getchef.com/api/v1",
				DownloadURL:  "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.12.0/download",
				Dependencies: map[string]string{"runit": ">= 0.0.0", "couchdb": ">= 0.0.0"},
			},
		},
	}
	ts := startHTTP(data)
	defer ts.Close()
	i := apiinstance.APIInstance{
		Component: component.Component{Endpoint: ts.URL},
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
