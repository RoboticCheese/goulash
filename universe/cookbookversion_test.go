package universe

import (
	"testing"
)

func cvdata() (data *CookbookVersion) {
	data = &CookbookVersion{
		Version:      "0.1.0",
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
		},
	}
	return
}

func TestNewCookbookVersion(t *testing.T) {
	cv := NewCookbookVersion()
	for _, i := range [][]interface{}{
		{cv.Version, ""},
		{cv.LocationType, ""},
		{cv.LocationPath, ""},
		{cv.DownloadURL, ""},
		{len(cv.Dependencies), 0},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookVersionEmptyIsEmpty(t *testing.T) {
	data := new(CookbookVersion)
	res := data.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestCookbookVersionEmptyStillEmpty(t *testing.T) {
	data := NewCookbookVersion()
	res := data.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasVersionStr(t *testing.T) {
	data := NewCookbookVersion()
	data.Version = "0.1.0"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasDownloadURL(t *testing.T) {
	data := NewCookbookVersion()
	data.DownloadURL = "http://example.com"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasDependencies(t *testing.T) {
	data := NewCookbookVersion()
	data.Dependencies["thing1"] = ">= 0.0.0"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEqualsEqual(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != true {
			t.Fatalf("Expected true, got: %v", res)
		}
	}
}

func TestCookbookVersionEqualsDifferentVersion(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Version = "9.9.9"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestCookbookVersionEqualsDifferentLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = "copsode"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestCookbookVersionEqualsDifferentLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationPath = "https://example2.com"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestCookbookVersionEqualsDifferentDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.DownloadURL = "https://example2.com/dl2"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestCookbookVersionEqualsDifferentDependenciesKeys(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies["thing3"] = ">= 0.0.0"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestCookbookVersionEqualsDifferentDependenciesValues(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies["thing1"] = "~> 1.0.0"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestCookbookVersionDiffEqual(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
}

func TestCookbookVersionDiffDataAddedAndDeleted(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = "souschef"
	data2.LocationPath = ""
	pos, neg := data1.Diff(data2)
	for _, v := range [][]interface{}{
		{pos.LocationType, "souschef"},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{len(pos.Dependencies), 0},
		{neg.LocationType, "opscode"},
		{neg.LocationPath, "https://example1.com"},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffChangedVersion(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Version = "9.9.9"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]interface{}{
		{pos.Version, "9.9.9"},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{len(pos.Dependencies), 0},
		{neg.Version, "0.1.0"},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffChangedLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationPath = "https://exam.ple"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]interface{}{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, "https://exam.ple"},
		{pos.DownloadURL, ""},
		{len(pos.Dependencies), 0},
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, "https://example1.com"},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffChangedLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = "souschef"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]interface{}{
		{pos.Version, ""},
		{pos.LocationType, "souschef"},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{len(pos.Dependencies), 0},
		{neg.Version, ""},
		{neg.LocationType, "opscode"},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffChangedDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.DownloadURL = "https://thing.co"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]interface{}{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, "https://thing.co"},
		{len(pos.Dependencies), 0},
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, "https://example1.com/dl1"},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffChangedOneDependency(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies = map[string]string{
		"thing1": ">= 0.0.0",
		"thing2": "~> 0.0.1",
	}
	pos, neg := data1.Diff(data2)
	for _, v := range [][]interface{}{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{len(pos.Dependencies), 1},
		{pos.Dependencies["thing2"], "~> 0.0.1"},
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 1},
		{neg.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffAddedNewDepndency(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies["thing3"] = ">= 0.0.0"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]interface{}{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{len(pos.Dependencies), 1},
		{pos.Dependencies["thing3"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
}

func TestCookbookVersionDiffRemovedVersion(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Version = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]interface{}{
		{neg.Version, "0.1.0"},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffRemovedLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationPath = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]interface{}{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, "https://example1.com"},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffRemovedLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]interface{}{
		{neg.Version, ""},
		{neg.LocationType, "opscode"},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffRemovedDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.DownloadURL = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]interface{}{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, "https://example1.com/dl1"},
		{len(neg.Dependencies), 0},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffRemovedOneDependency(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies = map[string]string{
		"thing1": ">= 0.0.0",
	}
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]interface{}{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 1},
		{neg.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func TestCookbookVersionDiffRemovedAllDepndencies(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies = map[string]string{}
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]interface{}{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{len(neg.Dependencies), 2},
		{neg.Dependencies["thing1"], ">= 0.0.0"},
		{neg.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}
