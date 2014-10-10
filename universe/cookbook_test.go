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
	data2.Versions = map[string]*CookbookVersion{
		"9.9.9": &CookbookVersion{},
	}
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

func Test_CpositiveDiff_1_Equal(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	res := data1.positiveDiff(data2)
	if res != nil {
		t.Fatalf("Expected nil, got: %v", res)
	}
}

func Test_CpositiveDiff_2_ChangedName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "somethingelse"
	res := data1.positiveDiff(data2)
	if res.Name != "somethingelse" {
		t.Fatalf("Expected changed name, got: %v", res.Name)
	}
	if len(res.Versions) != 0 {
		t.Fatalf("Expected 0 versions, got: %v", len(res.Versions))
	}
}

func Test_CpositiveDiff_3_NewVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["9.9.9"] = &CookbookVersion{Version: "9.9.9"}
	res := data1.positiveDiff(data2)
	if res.Name != "" {
		t.Fatalf("Expected unchanged name, got: %v", res.Name)
	}
	if len(res.Versions) != 1 {
		t.Fatalf("Expected 1 versions, got: %v", len(res.Versions))
	}
	if res.Versions["9.9.9"].Version != "9.9.9" {
		t.Fatalf("Expected ver 9.9.9, got: %v", res.Versions["9.9.9"])
	}
}

func Test_CpositiveDiff_4_ChangedVersion(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["0.1.0"].LocationType = "elsewhere"
	res := data1.positiveDiff(data2)
	if res.Name != "" {
		t.Fatalf("Expected unchanged name, got: %v", res.Name)
	}
	if len(res.Versions) != 1 {
		t.Fatalf("Expected 1 versions, got: %v", len(res.Versions))
	}
	if res.Versions["0.1.0"].LocationType != "elsewhere" {
		t.Fatalf("Expected 'elsewhere', got: %v",
			res.Versions["0.1.0"].LocationType)
	}
	if res.Versions["0.1.0"].LocationPath != "" {
		t.Fatalf("Expected empty string, got: %v",
			res.Versions["0.1.0"].LocationPath)
	}
}

func Test_CnegativeDiff_1_Equal(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	res := data1.negativeDiff(data2)
	if res != nil {
		t.Fatalf("Expected nil, got: %v", res)
	}
}

func Test_CnegativeDiff_2_ChangedName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "somethingelse"
	res := data1.negativeDiff(data2)
	if len(res.Versions) != 0 {
		t.Fatalf("Expected 0 versions, got: %v", len(res.Versions))
	}
}

func Test_CnegativeDiff_3_RemovedVersions(t *testing.T) {
	data1 := cdata()
	data2 := &Cookbook{Name: "something"}
	res := data1.negativeDiff(data2)
	if res.Name != "" {
		t.Fatalf("Expected unchanged name, got: %v", res.Name)
	}
	if len(res.Versions) != 1 {
		t.Fatalf("Expected 1 versions, got: %v", len(res.Versions))
	}
	if res.Versions["0.1.0"].LocationType != "opscode" {
		t.Fatalf("Expected 'opscode', got: %v",
			res.Versions["0.1.0"].LocationType)
	}
	if res.Versions["0.1.0"].LocationPath != "https://example1.com" {
		t.Fatalf("Expected 'https://example1.com', got: %v",
			res.Versions["0.1.0"].LocationPath)
	}
}

func Test_CnegativeDiff_4_ChangedVersion(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions["0.1.0"].LocationType = "elsewhere"
	res := data1.negativeDiff(data2)
	if len(res.Versions) != 0 {
		t.Fatalf("Expected 0 versions, got: %v", len(res.Versions))
	}
}
