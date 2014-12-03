package goulash

import (
	"net/http"
	"testing"
)

func cvdata() (data CookbookVersion) {
	data = CookbookVersion{
		Component:       Component{Endpoint: "https://example1.com"},
		License:         "oss",
		TarballFileSize: 123,
		Version:         "1.2.3",
		AverageRating:   0,
		Cookbook:        "https://example1.com/cookbook1",
		File:            "https://example1.com/cookbook1/file",
		Dependencies: map[string]string{
			"thing1": ">= 0.0.0",
		},
	}
	return
}

var cvjsonData = map[string]string{
	"license":           "Apache v2.0",
	"tarball_file_size": "5913",
	"version":           "2.0.0",
	"average_rating":    "null",
	"cookbook":          "https://supermarket.chef.io/api/v1/cookbooks/chef-dk",
	"file":              "https://supermarket.chef.io/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
	"dependencies":      `{"dmg": "~> 2.2"}`,
}

func cvjsonified() (res string) {
	res = `{"license": "` + cvjsonData["license"] + `",` +
		`"tarball_file_size": ` + cvjsonData["tarball_file_size"] + `,` +
		`"version": "` + cvjsonData["version"] + `",` +
		`"average_rating": ` + cvjsonData["average_rating"] + `,` +
		`"cookbook": "` + cvjsonData["cookbook"] + `",` +
		`"file": "` + cvjsonData["file"] + `",` +
		`"dependencies": ` + cvjsonData["dependencies"] + `}`
	return
}

func TestNewCookbookVersionNoError(t *testing.T) {
	ts := StartHTTP(cvjsonified(), nil)
	defer ts.Close()

	cb := new(Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	cv, err := NewCookbookVersion(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for _, i := range [][]interface{}{
		{cv.Endpoint, ts.URL + "/api/v1/cookbooks/chef-dk/versions/2.0.0"},
		{cv.License, cvjsonData["license"]},
		{cv.TarballFileSize, 5913},
		{cv.Version, cvjsonData["version"]},
		{cv.AverageRating, 0},
		{cv.Cookbook, cvjsonData["cookbook"]},
		{cv.File, cvjsonData["file"]},
		{cv.Dependencies["dmg"], "~> 2.2"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}

func TestNewCookbookVersionAverageRating(t *testing.T) {
	cvjsonData["average_rating"] = "20"
	ts := StartHTTP(cvjsonified(), nil)
	defer ts.Close()

	cb := new(Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	cv, err := NewCookbookVersion(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cv.AverageRating != 20 {
		t.Fatalf("Expected: 20, got: %v", cv.AverageRating)
	}
}

func TestNewCookbookVersionConnError(t *testing.T) {
	ts := StartHTTP(cvjsonified(), nil)
	ts.Close()

	cb := new(Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	_, err := NewCookbookVersion(cb, "2.0.0")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewCookbookVersion404Error(t *testing.T) {
	ts := StartHTTP(http.NotFound, nil)
	defer ts.Close()

	cb := new(Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	_, err := NewCookbookVersion(cb, "2.0.0")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewCookbookVersionRealData(t *testing.T) {
	cb := new(Cookbook)
	cb.Endpoint = "https://supermarket.chef.io/api/v1/cookbooks/chef-dk"
	cv, err := NewCookbookVersion(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for _, i := range [][]interface{}{
		{cv.Endpoint, "https://supermarket.chef.io/api/v1/cookbooks/chef-dk/versions/2.0.0"},
		{cv.License, "Apache v2.0"},
		{cv.TarballFileSize, 5913},
		{cv.Version, "2.0.0"},
		{cv.AverageRating, 0},
		{cv.Cookbook, "https://supermarket.chef.io/api/v1/cookbooks/chef-dk"},
		{cv.File, "https://supermarket.chef.io/api/v1/cookbooks/chef-dk/versions/2.0.0/download"},
		{len(cv.Dependencies), 1},
		{cv.Dependencies["dmg"], "~> 2.2"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestInitCookbookVersionEmptyStruct(t *testing.T) {
	cv := InitCookbookVersion()
	for _, i := range [][]interface{}{
		{cv.License, ""},
		{cv.TarballFileSize, 0},
		{cv.Version, ""},
		{cv.AverageRating, 0},
		{cv.Cookbook, ""},
		{cv.File, ""},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
	if len(cv.Dependencies) != 0 {
		t.Fatalf("Expected no dependencies, got: %v", cv.Dependencies)
	}
}

func TestCookbookVersionEmptyIsEmpty(t *testing.T) {
	data := new(CookbookVersion)
	res := data.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasEndpoint(t *testing.T) {
	data := new(CookbookVersion)
	data.Endpoint = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasLicense(t *testing.T) {
	data := new(CookbookVersion)
	data.License = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasTarballFileSize(t *testing.T) {
	data := new(CookbookVersion)
	data.TarballFileSize = 21
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasVersion(t *testing.T) {
	data := new(CookbookVersion)
	data.Version = "1.2.3"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasAverageRating(t *testing.T) {
	data := new(CookbookVersion)
	data.AverageRating = 1
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasCookbook(t *testing.T) {
	data := new(CookbookVersion)
	data.Cookbook = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasFile(t *testing.T) {
	data := new(CookbookVersion)
	data.File = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEmptyHasDependencies(t *testing.T) {
	data := new(CookbookVersion)
	data.Dependencies = map[string]string{"thing1": "1.2.3"}
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEqualsEqual(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	res := data1.Equals(&data2)
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestCookbookVersionEqualsDifferentEndpoints(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Endpoint = "https://somewhereelse.com"
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEqualsDifferentLicense(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.License = "closedsource"
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEqualsDifferentFileSize(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.TarballFileSize = 1
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionEqualsDifferentDependencies(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.Dependencies["thing2"] = ">= 0.0.0"
	res := data1.Equals(&data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(&data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestCookbookVersionDiffEqual(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	for _, i := range []*CookbookVersion{
		pos1,
		neg1,
		pos2,
		neg2,
	} {
		if i != nil {
			t.Fatalf("Expected nil, got: %v", i)
		}
	}
}

func TestCookbookVersionDiffDataAddedAndDeleted(t *testing.T) {
	data1 := cvdata()
	data2 := cvdata()
	data2.License = ""
	data2.File = "otherfile"
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	for _, i := range [][]interface{}{
		{pos1.File, "otherfile"},
		{neg1.License, "oss"},
		{pos2.License, "oss"},
		{pos2.File, "https://example1.com/cookbook1/file"},
		{neg2.License, ""},
		{neg2.File, "otherfile"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}
