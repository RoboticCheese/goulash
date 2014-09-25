package cookbook_version

import (
	"fmt"
	"github.com/RoboticCheese/goulash/cookbook"
	"net/http"
	"net/http/httptest"
	"testing"
)

var test_data = map[string]string{
	"license":           "Apache v2.0",
	"tarball_file_size": "5913",
	"version":           "2.0.0",
	"average_rating":    "null",
	"cookbook":          "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk",
	"file":              "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
	"dependencies":      `{"dmg": "~> 2.2"}`,
}

func jsonified() (res string) {
	res = `{"license": "` + test_data["license"] + `",` +
		`"tarball_file_size": ` + test_data["tarball_file_size"] + `,` +
		`"version": "` + test_data["version"] + `",` +
		`"average_rating": ` + test_data["average_rating"] + `,` +
		`"cookbook": "` + test_data["cookbook"] + `",` +
		`"file": "` + test_data["file"] + `",` +
		`"dependencies": ` + test_data["dependencies"] + `}`
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
		cv.License:  test_data["license"],
		cv.Version:  test_data["version"],
		cv.Cookbook: test_data["cookbook"],
		cv.File:     test_data["file"],
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
	test_data["average_rating"] = "20"
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
