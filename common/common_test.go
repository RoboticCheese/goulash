package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var http_etag = ""
var http_data = "SOME HTTP DATA"

func start_http() (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if http_etag != "" {
					w.Header().Set("ETag", http_etag)
				}
				fmt.Fprint(w, http_data)
			},
		),
	)
	return
}

func Test_Supermarketer_1(t *testing.T) {
	// Doesn't do anything just yet
}

func Test_Component_1(t *testing.T) {
	type Thing struct {
		Component
	}

	res := Thing{
		Component: Component{
			Endpoint: "something",
			ETag:     "anotherthing",
		},
	}
	if res.Endpoint != "something" {
		t.Fatalf("Expected 'something', got: %v", res.Endpoint)
	}
	if res.ETag != "anotherthing" {
		t.Fatalf("Expected 'anotherthing', got: %v", res.ETag)
	}
}

func Test_New_1_NoETag(t *testing.T) {
	ts := start_http()
	defer ts.Close()

	c, err := New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	if c.Endpoint != ts.URL {
		t.Fatalf("Expected '%v', got: %v", ts.URL, c.Endpoint)
	}
	if c.ETag != "" {
		t.Fatalf("Expected empty str, got: %v", c.ETag)
	}
}

func Test_New_2_ETag(t *testing.T) {
	http_etag = "hellothere"
	ts := start_http()
	defer ts.Close()

	c, err := New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	if c.Endpoint != ts.URL {
		t.Fatalf("Expected '%v', got: %v", ts.URL, c.Endpoint)
	}
	if c.ETag != "hellothere" {
		t.Fatalf("Expected 'hellothere', got: %v", c.ETag)
	}
}

func Test_Empty_1_Empty(t *testing.T) {
	c := new(Component)
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Empty_2_HasEndpoint(t *testing.T) {
	c := new(Component)
	c.Endpoint = "https://example.com"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_3_HasETag(t *testing.T) {
	c := new(Component)
	c.ETag = "thing"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}
