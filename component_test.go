package goulash

import (
	"testing"

	"github.com/RoboticCheese/goulash/common"
)

func TestComponent(t *testing.T) {
	type Thing struct {
		Component
	}
	res := Thing{
		Component: Component{
			Endpoint: "something",
			ETag:     "anotherthing",
		},
	}
	for _, i := range [][]string{
		{res.Endpoint, "something"},
		{res.ETag, "anotherthing"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestNewComponentNoETag(t *testing.T) {
	ts := StartHTTP("", nil)
	defer ts.Close()

	c, err := NewComponent(ts.URL)
	for _, i := range [][]interface{}{
		{err, nil},
		{c.Endpoint, ts.URL},
		{c.ETag, ""},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestNewComponentETag(t *testing.T) {
	ts := StartHTTP("", map[string]string{"ETag": "hellothere"})
	defer ts.Close()

	c, err := NewComponent(ts.URL)
	for _, i := range [][]interface{}{
		{err, nil},
		{c.Endpoint, ts.URL},
		{c.ETag, "hellothere"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestInitComponentEmptyStruct(t *testing.T) {
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

func TestComponentEmptyIsEmpty(t *testing.T) {
	c := new(Component)
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestComponentEmptyHasEndpoint(t *testing.T) {
	c := new(Component)
	c.Endpoint = "https://example.com"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestComponentEmptyHasETag(t *testing.T) {
	c := new(Component)
	c.ETag = "thing"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestComponentDiffEqual(t *testing.T) {
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

func TestComponentDiffAddedAndDeletedData(t *testing.T) {
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

func TestComponentDiffChangedData(t *testing.T) {
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
