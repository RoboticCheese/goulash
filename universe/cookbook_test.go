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

func Test_NewCookbook_1(t *testing.T) {
	c := NewCookbook()
	if c.Name != "" {
		t.Fatalf("Expected empty string, got: %v", c.Name)
	}
	if len(c.Versions) != 0 {
		t.Fatalf("Expected 0 versions, got: %v", len(c.Versions))
	}
}

func Test_CEmpty_1_Empty(t *testing.T) {
	c := new(Cookbook)
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_CEmpty_2_StillEmpty(t *testing.T) {
	c := NewCookbook()
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_CEmpty_3_HasName(t *testing.T) {
	c := NewCookbook()
	c.Name = "thing"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_CEmpty_3_HasEmptyVersions(t *testing.T) {
	c := NewCookbook()
	c.Versions["0.1.0"] = NewCookbookVersion()
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_CEmpty_3_HasNonEmptyVersions(t *testing.T) {
	c := NewCookbook()
	c.Versions["0.1.0"] = NewCookbookVersion()
	c.Versions["0.1.0"].LocationType = "opscode"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_CEquals_1_Equal(t *testing.T) {
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

func Test_CEquals_2_ChangedName(t *testing.T) {
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

func Test_CEquals_3_MoreVersions(t *testing.T) {
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

func Test_CEquals_4_FewerVersions(t *testing.T) {
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

func Test_CEquals_5_DifferentVersions(t *testing.T) {
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

func Test_CDiff_1_Equal(t *testing.T) {
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

func Test_CDiff_2_DataAddedAndDeleted(t *testing.T) {
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

func Test_CDiff_3_ChangedName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "somethingelse"
	pos, neg := data1.Diff(data2)
	for _, i := range [][]string{
		{pos.Name, "somethingelse"},
		{neg.Name, "something"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
	if len(pos.Versions) != 0 {
		t.Fatalf("Expected 0 versions, got: %v", len(pos.Versions))
	}
	if len(neg.Versions) != 0 {
		t.Fatalf("Expected 0 versions, got: %v", len(neg.Versions))
	}
}

func Test_CDiff_4_NewVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["9.9.9"] = &CookbookVersion{Version: "9.9.9"}
	pos, neg := data1.Diff(data2)
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
	for _, i := range [][]string{
		{pos.Name, ""},
		{pos.Versions["9.9.9"].Version, "9.9.9"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
	if len(pos.Versions) != 1 {
		t.Fatalf("Expected 1 versions, got: %v", len(pos.Versions))
	}
}

func Test_CDiff_5_ChangedVersion(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["0.1.0"].LocationType = "elsewhere"
	pos, neg := data1.Diff(data2)
	for _, i := range [][]string{
		{pos.Name, ""},
		{pos.Versions["0.1.0"].LocationType, "elsewhere"},
		{pos.Versions["0.1.0"].LocationPath, ""},
		{neg.Name, ""},
		{neg.Versions["0.1.0"].LocationType, "opscode"},
		{neg.Versions["0.1.0"].LocationPath, ""},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
	if len(pos.Versions) != 1 {
		t.Fatalf("Expected 1 versions, got: %v", len(pos.Versions))
	}
	if len(neg.Versions) != 1 {
		t.Fatalf("Expected 1 versions, got: %v", len(neg.Versions))
	}
}

func Test_CDiff_6_RemovedVersions(t *testing.T) {
	data1 := cdata()
	data2 := &Cookbook{Name: "something"}
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, i := range [][]string{
		{neg.Name, ""},
		{neg.Versions["0.1.0"].LocationType, "opscode"},
		{neg.Versions["0.1.0"].LocationPath, "https://example1.com"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
	if len(neg.Versions) != 1 {
		t.Fatalf("Expected 1 versions, got: %v", len(neg.Versions))
	}
}
