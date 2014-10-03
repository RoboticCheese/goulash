package universe

import (
	"testing"
)

func cvdata() (data CookbookVersion) {
	data = CookbookVersion{
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

func Test_CVDiff_1_Equal(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	pos, neg := data1.Diff(&data2)
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
	pos, neg := data1.Diff(&data2)
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
		{neg.LocationType, ""},
		{neg.LocationPath, "https://example1.com"},
		{neg.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func Test_CVEquals_1_Equal(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	for _, res := range []bool{
		data1.Equals(&data2),
		data2.Equals(&data1),
	} {
		if res != true {
			t.Fatalf("Expected true, got: %v", res)
		}
	}
}

func Test_CVEquals_2_DifferentLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationType = "copsode"
	for _, res := range []bool{
		data1.Equals(&data2),
		data2.Equals(&data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func Test_CVEquals_3_DifferentLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.LocationPath = "https://example2.com"
	for _, res := range []bool{
		data1.Equals(&data2),
		data2.Equals(&data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func Test_CVEquals_4_DifferentDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.DownloadURL = "https://example2.com/dl2"
	for _, res := range []bool{
		data1.Equals(&data2),
		data2.Equals(&data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func Test_CVEquals_5_DifferentDependenciesKeys(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies["thing3"] = ">= 0.0.0"
	for _, res := range []bool{
		data1.Equals(&data2),
		data2.Equals(&data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func Test_CVEquals_6_DifferentDependenciesValues(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies["thing1"] = "~> 1.0.0"
	for _, res := range []bool{
		data1.Equals(&data2),
		data2.Equals(&data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func Test_CVpositiveDiff_1_Equal(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	res := data1.positiveDiff(&data2)
	if res != nil {
		t.Fatalf("Expected nil, got: %v", res)
	}
}

func Test_CVpositiveDiff_2_ChangedLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://exam.ple",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
		},
	}
	res := data1.positiveDiff(&data2)
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, "https://exam.ple"},
		{res.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(res.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(res.Dependencies))
	}
}

func Test_CVpositiveDiff_3_ChangedLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationPath: "https://example1.com",
		LocationType: "souschef",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
		},
	}
	res := data1.positiveDiff(&data2)
	for _, v := range [][]string{
		{res.LocationType, "souschef"},
		{res.LocationPath, ""},
		{res.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(res.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(res.Dependencies))
	}
}

func Test_CVpositiveDiff_4_ChangedDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://thing.co",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
		},
	}
	res := data1.positiveDiff(&data2)
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, ""},
		{res.DownloadURL, "https://thing.co"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(res.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(res.Dependencies))
	}
}

func Test_CVpositiveDiff_5_ChangeOneDependency(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": "~> 0.0.1",
		},
	}
	res := data1.positiveDiff(&data2)
	if len(res.Dependencies) != 1 {
		t.Fatalf("Expected 1 dep, got: %v", len(res.Dependencies))
	}
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, ""},
		{res.DownloadURL, ""},
		{res.Dependencies["thing2"], "~> 0.0.1"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func Test_CVpositiveDiff_6_AddedNewDepndency(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
			"thing3": ">= 0.0.0",
		},
	}
	res := data1.positiveDiff(&data2)
	if len(res.Dependencies) != 1 {
		t.Fatalf("Expected 1 dep, got: %v", len(res.Dependencies))
	}
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, ""},
		{res.DownloadURL, ""},
		{res.Dependencies["thing3"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func Test_CVnegativeDiff_1_Equal(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	res := data1.negativeDiff(&data2)
	if res != nil {
		t.Fatalf("Expected nil, got: %v", res)
	}
}

func Test_CVnegativeDiff_2_RemovedLocationPath(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
		},
	}
	res := data1.negativeDiff(&data2)
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, "https://example1.com"},
		{res.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(res.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(res.Dependencies))
	}
}

func Test_CVnegativeDiff_3_RemovedLocationType(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
		},
	}
	res := data1.negativeDiff(&data2)
	for _, v := range [][]string{
		{res.LocationType, "opscode"},
		{res.LocationPath, ""},
		{res.DownloadURL, ""},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(res.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(res.Dependencies))
	}
}

func Test_CVnegativeDiff_4_RemovedDownloadURL(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
			"thing2": ">= 0.0.0",
		},
	}
	res := data1.negativeDiff(&data2)
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, ""},
		{res.DownloadURL, "https://example1.com/dl1"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
	if len(res.Dependencies) != 0 {
		t.Fatalf("Expected 0 deps, got: %v", len(res.Dependencies))
	}
}

func Test_CVnegativeDiff_5_RemovedOneDependency(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
		},
	}
	res := data1.negativeDiff(&data2)
	if len(res.Dependencies) != 1 {
		t.Fatalf("Expected 1 dep, got: %v", len(res.Dependencies))
	}
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, ""},
		{res.DownloadURL, ""},
		{res.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}

func Test_CVnegativeDiff_6_RemovedAllDepndencies(t *testing.T) {
	data1 := cvdata()
	data2 := CookbookVersion{
		LocationType: "opscode",
		LocationPath: "https://example1.com",
		DownloadURL:  "https://example1.com/dl1",
	}
	res := data1.negativeDiff(&data2)
	if len(res.Dependencies) != 2 {
		t.Fatalf("Expected 2 deps, got: %v", len(res.Dependencies))
	}
	for _, v := range [][]string{
		{res.LocationType, ""},
		{res.LocationPath, ""},
		{res.DownloadURL, ""},
		{res.Dependencies["thing1"], ">= 0.0.0"},
		{res.Dependencies["thing2"], ">= 0.0.0"},
	} {
		if v[0] != v[1] {
			t.Fatalf("Expected: %v, got: %v", v[1], v[0])
		}
	}
}
