package cookbook_version

import (
	"fmt"
	"github.com/RoboticCheese/goulash/common"
	"github.com/RoboticCheese/goulash/cookbook"
	"net/http"
	"net/http/httptest"
	"testing"
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

var json_data = map[string]string{
	"license":           "Apache v2.0",
	"tarball_file_size": "5913",
	"version":           "2.0.0",
	"average_rating":    "null",
	"cookbook":          "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk",
	"file":              "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
	"dependencies":      `{"dmg": "~> 2.2"}`,
}

func jsonified() (res string) {
	res = `{"license": "` + json_data["license"] + `",` +
		`"tarball_file_size": ` + json_data["tarball_file_size"] + `,` +
		`"version": "` + json_data["version"] + `",` +
		`"average_rating": ` + json_data["average_rating"] + `,` +
		`"cookbook": "` + json_data["cookbook"] + `",` +
		`"file": "` + json_data["file"] + `",` +
		`"dependencies": ` + json_data["dependencies"] + `}`
	return
}

func start_http() (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, jsonified())
			},
		),
	)
	return
}

func Test_Equals_1_Equal(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Equals_2_DifferentEndpoints(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.Endpoint = "https://somewhereelse.com"
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_3_DifferentLicense(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.License = "closedsource"
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_4_DifferentFileSize(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.TarballFileSize = 1
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_5_DifferentDependencies(t *testing.T) {
	data1 := cvdata1()
	data2 := cvdata2()
	data2.Dependencies["thing2"] = ">= 0.0.0"
	res, err := data1.Equals(data2)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
	res, err = data2.Equals(data1)
	if err != nil {
		t.Fatalf("Expected no errors, got: %v", err)
	}
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_New_1_NoError(t *testing.T) {
	ts := start_http()
	defer ts.Close()

	cb := new(cookbook.Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	cv, err := New(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		cv.Endpoint: ts.URL + "/api/v1/cookbooks/chef-dk/versions/2.0.0",
		cv.License:  json_data["license"],
		cv.Version:  json_data["version"],
		cv.Cookbook: json_data["cookbook"],
		cv.File:     json_data["file"],
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
	json_data["average_rating"] = "20"
	ts := start_http()
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
	ts := start_http()
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
