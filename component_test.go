package goulash

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoboticCheese/goulash/common"
)

var compHTTPETag = ""
var compHTTPData = "SOME HTTP DATA"

func compStartHTTP() (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if compHTTPETag != "" {
					w.Header().Set("ETag", compHTTPETag)
				}
				fmt.Fprint(w, compHTTPData)
			},
		),
	)
	return
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

func Test_NewComponent_1_NoETag(t *testing.T) {
	ts := compStartHTTP()
	defer ts.Close()

	c, err := NewComponent(ts.URL)
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

func Test_NewComponent_2_ETag(t *testing.T) {
	compHTTPETag = "hellothere"
	ts := compStartHTTP()
	defer ts.Close()

	c, err := NewComponent(ts.URL)
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

func Test_InitComponent_1_EmptyStruct(t *testing.T) {
	c := InitComponent()
	for _, k := range []string{
		c.Endpoint,
		c.ETag,
	} {
		if k != "" {
			t.Fatalf("Expected empty string, got: %v", k)
		}
	}
}

func Test_Component_Empty_1_Empty(t *testing.T) {
	c := new(Component)
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Component_Empty_2_HasEndpoint(t *testing.T) {
	c := new(Component)
	c.Endpoint = "https://example.com"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Component_Empty_3_HasETag(t *testing.T) {
	c := new(Component)
	c.ETag = "thing"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Component_Diff_1_Equal(t *testing.T) {
	c1 := Component{Endpoint: "abc", ETag: "def"}
	c2 := Component{Endpoint: "abc", ETag: "def"}
	pos1, neg1 := c1.Diff(&c2)
	pos2, neg2 := c2.Diff(&c1)
	for _, i := range []*Component{pos1, neg1, pos2, neg2} {
		if i != nil {
			t.Fatalf("Expected nil, got: %v", i)
		}
	}
}

func Test_Component_Diff_2_AddedAndDeletedData(t *testing.T) {
	c1 := Component{}
	c2 := Component{Endpoint: "abc", ETag: "def"}
	pos1, neg1 := c1.Diff(&c2)
	pos2, neg2 := c2.Diff(&c1)
	for _, i := range []*Component{neg1, pos2} {
		if i != nil {
			t.Fatalf("Expected nil, got: %v", i)
		}
	}
	for _, k := range [][]*Component{
		{pos1, &Component{Endpoint: "abc", ETag: "def"}},
		{neg2, &Component{Endpoint: "abc", ETag: "def"}},
	} {
		if !common.Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}

func Test_Component_Diff_3_ChangedData(t *testing.T) {
	c1 := Component{Endpoint: "abc", ETag: "def"}
	c2 := Component{Endpoint: "uvw", ETag: "xyz"}
	pos1, neg1 := c1.Diff(&c2)
	pos2, neg2 := c2.Diff(&c1)
	for _, k := range [][]common.Supermarketer{
		{pos1, &Component{Endpoint: "uvw", ETag: "xyz"}},
		{neg1, &Component{Endpoint: "abc", ETag: "def"}},
		{pos2, &Component{Endpoint: "abc", ETag: "def"}},
		{neg2, &Component{Endpoint: "uvw", ETag: "xyz"}},
	} {
		if !common.Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}
