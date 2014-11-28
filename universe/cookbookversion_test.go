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

func Test_NewCookbookVersion_1(t *testing.T) {
	cv := NewCookbookVersion()
	for _, i := range []string{
		cv.Version,
		cv.LocationType,
		cv.LocationPath,
		cv.DownloadURL,
	} {
		if i != "" {
			t.Fatalf("Expected empty string, got: %v", i)
		}
	}
	if len(cv.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(cv.Dependencies))
	}
}

func Test_CVEmpty_1_Empty(t *testing.T) {
	data := new(CookbookVersion)
	res := data.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_CVEmpty_2_StillEmpty(t *testing.T) {
	data := NewCookbookVersion()
	res := data.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_CVEmpty_3_HasVersionStr(t *testing.T) {
	data := NewCookbookVersion()
	data.Version = "0.1.0"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_CVEmpty_4_HasDownloadURL(t *testing.T) {
	data := NewCookbookVersion()
	data.DownloadURL = "http://example.com"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_CVEmpty_5_HasDependencies(t *testing.T) {
	data := NewCookbookVersion()
	data.Dependencies["thing1"] = ">= 0.0.0"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_CVEquals_1_Equal(t *testing.T) {
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

func Test_CVEquals_2_DifferentVersion(t *testing.T) {
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

func Test_CVEquals_3_DifferentLocationType(t *testing.T) {
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

func Test_CVEquals_4_DifferentLocationPath(t *testing.T) {
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

func Test_CVEquals_5_DifferentDownloadURL(t *testing.T) {
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

func Test_CVEquals_6_DifferentDependenciesKeys(t *testing.T) {
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

func Test_CVEquals_7_DifferentDependenciesValues(t *testing.T) {
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

func Test_CVDiff_1_Equal(t *testing.T) {
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

func Test_CVDiff_2_DataAddedAndDeleted(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = "souschef"
	data2.LocationPath = ""
	pos, neg := data1.Diff(data2)
	if len(pos.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(pos.Dependencies))
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
	for _, v := range [][]string{
		{pos.LocationType, "souschef"},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{neg.LocationType, "opscode"},
		{neg.LocationPath, "https://example1.com"},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func Test_CVDiff_3_ChangedVersion(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Version = "9.9.9"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]string{
		{pos.Version, "9.9.9"},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{neg.Version, "0.1.0"},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(pos.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(pos.Dependencies))
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_4_ChangedLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationPath = "https://exam.ple"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]string{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, "https://exam.ple"},
		{pos.DownloadURL, ""},
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, "https://example1.com"},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(pos.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(pos.Dependencies))
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_5_ChangedLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = "souschef"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]string{
		{pos.Version, ""},
		{pos.LocationType, "souschef"},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{neg.Version, ""},
		{neg.LocationType, "opscode"},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(pos.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(pos.Dependencies))
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_6_ChangedDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.DownloadURL = "https://thing.co"
	pos, neg := data1.Diff(data2)
	for _, v := range [][]string{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, "https://thing.co"},
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, "https://example1.com/dl1"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(pos.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(pos.Dependencies))
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_7_ChangedOneDependency(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies = map[string]string{
		"thing1": ">= 0.0.0",
		"thing2": "~> 0.0.1",
	}
	pos, neg := data1.Diff(data2)
	if len(pos.Dependencies) != 1 {
		t.Fatalf("Expected 1 dep, got: %v", len(pos.Dependencies))
	}
	if len(neg.Dependencies) != 1 {
		t.Fatalf("Expected 1 dep, got: %v", len(neg.Dependencies))
	}
	for _, v := range [][]string{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
		{pos.Dependencies["thing2"], "~> 0.0.1"},
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{neg.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func Test_CVDiff_8_AddedNewDepndency(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies["thing3"] = ">= 0.0.0"
	pos, neg := data1.Diff(data2)
	if len(pos.Dependencies) != 1 {
		t.Fatalf("Expected 1 dep, got: %v", len(pos.Dependencies))
	}
	for _, v := range [][]string{
		{pos.Version, ""},
		{pos.LocationType, ""},
		{pos.LocationPath, ""},
		{pos.DownloadURL, ""},
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

func Test_CVDiff_9_RemovedVersion(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Version = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]string{
		{neg.Version, "0.1.0"},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_10_RemovedLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationPath = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]string{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, "https://example1.com"},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_11_RemovedLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]string{
		{neg.Version, ""},
		{neg.LocationType, "opscode"},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_12_RemovedDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.DownloadURL = ""
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	for _, v := range [][]string{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, "https://example1.com/dl1"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(neg.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(neg.Dependencies))
	}
}

func Test_CVDiff_13_RemovedOneDependency(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies = map[string]string{
		"thing1": ">= 0.0.0",
	}
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	if len(neg.Dependencies) != 1 {
		t.Fatalf("Expected 1 dep, got: %v", len(neg.Dependencies))
	}
	for _, v := range [][]string{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{neg.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func Test_CVDiff_14_RemovedAllDepndencies(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies = map[string]string{}
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	if len(neg.Dependencies) != 2 {
		t.Fatalf("Expected 2 deps, got: %v", len(neg.Dependencies))
	}
	for _, v := range [][]string{
		{neg.Version, ""},
		{neg.LocationType, ""},
		{neg.LocationPath, ""},
		{neg.DownloadURL, ""},
		{neg.Dependencies["thing1"], ">= 0.0.0"},
		{neg.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}
