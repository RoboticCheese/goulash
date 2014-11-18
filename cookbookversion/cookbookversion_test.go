package cookbookversion

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoboticCheese/goulash/common"
	"github.com/RoboticCheese/goulash/cookbook"
)

func cvdata1() (data1 CookbookVersion) {
	data1 = CookbookVersion{
		Component:       common.Component{Endpoint: "https://example1.com"},
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

func cvdata2() (data2 CookbookVersion) {
	data2 = CookbookVersion{
		Component:       common.Component{Endpoint: "https://example1.com"},
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

var jsonData = map[string]string{
	"license":           "Apache v2.0",
	"tarball_file_size": "5913",
	"version":           "2.0.0",
	"average_rating":    "null",
	"cookbook":          "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk",
	"file":              "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
	"dependencies":      `{"dmg": "~> 2.2"}`,
}

func jsonified() (res string) {
	res = `{"license": "` + jsonData["license"] + `",` +
		`"tarball_file_size": ` + jsonData["tarball_file_size"] + `,` +
		`"version": "` + jsonData["version"] + `",` +
		`"average_rating": ` + jsonData["average_rating"] + `,` +
		`"cookbook": "` + jsonData["cookbook"] + `",` +
		`"file": "` + jsonData["file"] + `",` +
		`"dependencies": ` + jsonData["dependencies"] + `}`
	return
}

func startHTTP() (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, jsonified())
			},
		),
	)
	return
}

func Test_Empty_1_Empty(t *testing.T) {
	data := new(CookbookVersion)
	res := data.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Empty_2_HasEndpoint(t *testing.T) {
	data := new(CookbookVersion)
	data.Endpoint = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_3_HasLicense(t *testing.T) {
	data := new(CookbookVersion)
	data.License = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_4_HasTarballFileSize(t *testing.T) {
	data := new(CookbookVersion)
	data.TarballFileSize = 21
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_5_HasVersion(t *testing.T) {
	data := new(CookbookVersion)
	data.Version = "1.2.3"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_6_HasAverageRating(t *testing.T) {
	data := new(CookbookVersion)
	data.AverageRating = 1
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_7_HasCookbook(t *testing.T) {
	data := new(CookbookVersion)
	data.Cookbook = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_8_HasFile(t *testing.T) {
	data := new(CookbookVersion)
	data.File = "something"
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_8_HasDependencies(t *testing.T) {
	data := new(CookbookVersion)
	data.Dependencies = map[string]string{"thing1": "1.2.3"}
	res := data.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_1_Equal(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	res := data1.Equals(data2)
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
	res = data2.Equals(data1)
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Equals_2_DifferentEndpoints(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.Endpoint = "https://somewhereelse.com"
	res := data1.Equals(data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_3_DifferentLicense(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.License = "closedsource"
	res := data1.Equals(data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_4_DifferentFileSize(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.TarballFileSize = 1
	res := data1.Equals(data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_5_DifferentDependencies(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.Dependencies["thing2"] = ">= 0.0.0"
	res := data1.Equals(data2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res = data2.Equals(data1)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_New_1_NoError(t *testing.T) {
	ts := startHTTP()
	defer ts.Close()

	cb := new(cookbook.Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	cv, err := New(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		cv.Endpoint: ts.URL + "/api/v1/cookbooks/chef-dk/versions/2.0.0",
		cv.License:  jsonData["license"],
		cv.Version:  jsonData["version"],
		cv.Cookbook: jsonData["cookbook"],
		cv.File:     jsonData["file"],
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
	if cv.TarballFileSize != 5913 {
		t.Fatalf("Expected: 5913, got: %v", cv.TarballFileSize)
	}
	if cv.AverageRating != 0 {
		t.Fatalf("Expected: 0, got: %v", cv.AverageRating)
	}
	if len(cv.Dependencies) != 1 {
		t.Fatalf("Expected: 1 dependency, got: %v", len(cv.Dependencies))
	}
	if cv.Dependencies["dmg"] != "~> 2.2" {
		t.Fatalf("Expected: ~> 2.2, got: %v", cv.Dependencies["dmg"])
	}
}

func Test_New_2_AverageRating(t *testing.T) {
	jsonData["average_rating"] = "20"
	ts := startHTTP()
	defer ts.Close()

	cb := new(cookbook.Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	cv, err := New(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cv.AverageRating != 20 {
		t.Fatalf("Expected: 20, got: %v", cv.AverageRating)
	}
}

func Test_New_3_ConnError(t *testing.T) {
	ts := startHTTP()
	ts.Close()

	cb := new(cookbook.Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	_, err := New(cb, "2.0.0")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_4_404Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer ts.Close()

	cb := new(cookbook.Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	_, err := New(cb, "2.0.0")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_5_RealData(t *testing.T) {
	cb := new(cookbook.Cookbook)
	cb.Endpoint = "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk"
	cv, err := New(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		cv.Endpoint: "https://supermarket.getchef.com/api/v1/" +
			"cookbooks/chef-dk/versions/2.0.0",
		cv.License: "Apache v2.0",
		cv.Version: "2.0.0",
		cv.Cookbook: "https://supermarket.getchef.com/api/v1/" +
			"cookbooks/chef-dk",
		cv.File: "https://supermarket.getchef.com/api/v1/" +
			"cookbooks/chef-dk/versions/2.0.0/download",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
	if cv.TarballFileSize != 5913 {
		t.Fatalf("Expected: 5913, got: %v", cv.TarballFileSize)
	}
	if cv.AverageRating != 0 {
		t.Fatalf("Expected: 0, got: %v", cv.AverageRating)
	}
	if len(cv.Dependencies) != 1 {
		t.Fatalf("Expected: 1 dependency, got: %v", len(cv.Dependencies))
	}
	if cv.Dependencies["dmg"] != "~> 2.2" {
		t.Fatalf("Expected: ~> 2.2, got: %v", cv.Dependencies["dmg"])
	}
}
