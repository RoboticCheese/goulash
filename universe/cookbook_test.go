package universe

import (
	"testing"
)

func cdata() (data *Cookbook) {
	data = &Cookbook{
		Name: "something",
		Versions: map[string]*CookbookVersion{
			"0.1.0": &CookbookVersion{
				Version:      "0.1.0",
				LocationType: "opscode",
				LocationPath: "https://example1.com",
				DownloadURL:  "https://example1.com/dl1",
				Dependencies: map[string]string{
					"thing1": ">= 0.0.0",
				},
			},
		},
	}
	return
}

func TestNewCookbook(t *testing.T) {
	c := NewCookbook()
	for _, i := range [][]interface{}{
		{c.Name, ""},
		{len(c.Versions), 0},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookEmptyIsEmpty(t *testing.T) {
	c := new(Cookbook)
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestCookbookEmptyStillEmpty(t *testing.T) {
	c := NewCookbook()
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestCookbookEmptyHasName(t *testing.T) {
	c := NewCookbook()
	c.Name = "thing"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookEmptyHasEmptyVersions(t *testing.T) {
	c := NewCookbook()
	c.Versions["0.1.0"] = NewCookbookVersion()
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestCookbookEmptyHasNonEmptyVersions(t *testing.T) {
	c := NewCookbook()
	c.Versions["0.1.0"] = NewCookbookVersion()
	c.Versions["0.1.0"].LocationType = "opscode"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookEqualsEqual(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != true {
			t.Fatalf("Expected true, got: %v", res)
		}
	}
}

func TestCookbookEqualsChangedName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "otherthing"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestCookbookEqualsMoreVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["0.2.0"] = &CookbookVersion{}
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, git: %v", res)
		}
	}
}

func TestCookbookEqualsFewerVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions = map[string]*CookbookVersion{}
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, git: %v", res)
		}
	}
}

func TestCookbookEqualsDifferentVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["0.1.0"] = &CookbookVersion{}
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, git: %v", res)
		}
	}
}

func TestCookbookDiffEqual(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
}

func TestCookbookDiffDataAddedAndDeleted(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	delete(data2.Versions, "0.1.0")
	data2.Versions["9.9.9"] = &CookbookVersion{Version: "9.9.9"}
	pos, neg := data1.Diff(data2)
	for _, res := range []int{
		len(pos.Versions),
		len(neg.Versions),
	} {
		if res != 1 {
			t.Fatalf("Expected 1 version, got: %v", res)
		}
	}
}

func TestCookbookDiffChangedName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "somethingelse"
	pos, neg := data1.Diff(data2)
	for _, i := range [][]interface{}{
		{pos.Name, "somethingelse"},
		{len(pos.Versions), 0},
		{neg.Name, "something"},
		{len(neg.Versions), 0},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffNewVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["9.9.9"] = &CookbookVersion{Version: "9.9.9"}
	pos, neg := data1.Diff(data2)
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
	for _, i := range [][]interface{}{
		{pos.Name, ""},
		{len(pos.Versions), 1},
		{pos.Versions["9.9.9"].Version, "9.9.9"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffChangedVersion(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["0.1.0"].LocationType = "elsewhere"
	pos, neg := data1.Diff(data2)
	for _, i := range [][]interface{}{
		{pos.Name, ""},
		{len(pos.Versions), 1},
		{pos.Versions["0.1.0"].LocationType, "elsewhere"},
		{pos.Versions["0.1.0"].LocationPath, ""},
		{neg.Name, ""},
		{len(neg.Versions), 1},
		{neg.Versions["0.1.0"].LocationType, "opscode"},
		{neg.Versions["0.1.0"].LocationPath, ""},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffRemovedVersions(t *testing.T) {
	data1 := cdata()
	data2 := &Cookbook{Name: "something"}
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, i := range [][]interface{}{
		{neg.Name, ""},
		{len(neg.Versions), 1},
		{neg.Versions["0.1.0"].LocationType, "opscode"},
		{neg.Versions["0.1.0"].LocationPath, "https://example1.com"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}
