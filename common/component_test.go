package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var httpETag = ""
var httpData = "SOME HTTP DATA"

func startHTTP() (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if httpETag != "" {
					w.Header().Set("ETag", httpETag)
				}
				fmt.Fprint(w, httpData)
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

func Test_CNew_1_NoETag(t *testing.T) {
	ts := startHTTP()
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

func Test_CNew_2_ETag(t *testing.T) {
	httpETag = "hellothere"
	ts := startHTTP()
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

func Test_CEmpty_1_Empty(t *testing.T) {
	c := new(Component)
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_CEmpty_2_HasEndpoint(t *testing.T) {
	c := new(Component)
	c.Endpoint = "https://example.com"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_CEmpty_3_HasETag(t *testing.T) {
	c := new(Component)
	c.ETag = "thing"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_CDiff_1_Equal(t *testing.T) {
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

func Test_CDiff_2_AddedAndDeletedData(t *testing.T) {
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
		if !Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}

func Test_CDiff_3_ChangedData(t *testing.T) {
	c1 := Component{Endpoint: "abc", ETag: "def"}
	c2 := Component{Endpoint: "uvw", ETag: "xyz"}
	pos1, neg1 := c1.Diff(&c2)
	pos2, neg2 := c2.Diff(&c1)
	for _, k := range [][]Supermarketer{
		{pos1, &Component{Endpoint: "uvw", ETag: "xyz"}},
		{neg1, &Component{Endpoint: "abc", ETag: "def"}},
		{pos2, &Component{Endpoint: "abc", ETag: "def"}},
		{neg2, &Component{Endpoint: "uvw", ETag: "xyz"}},
	} {
		if !Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}
