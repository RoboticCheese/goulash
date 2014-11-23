package cookbookversion

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RoboticCheese/goulash/component"
	"github.com/RoboticCheese/goulash/cookbook"
)

func cvdata() (data CookbookVersion) {
	data = CookbookVersion{
		Component:       component.Component{Endpoint: "https://example1.com"},
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

func Test_New_1_NoError(t *testing.T) {
	ts := startHTTP()
	defer ts.Close()

	cb := new(cookbook.Cookbook)
	cb.Endpoint = ts.URL + "/api/v1/cookbooks/chef-dk"
	cv, err := New(cb, "2.0.0")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for _, i := range [][]interface{}{
		{cv.Endpoint, ts.URL + "/api/v1/cookbooks/chef-dk/versions/2.0.0"},
		{cv.License, jsonData["license"]},
		{cv.TarballFileSize, 5913},
		{cv.Version, jsonData["version"]},
		{cv.AverageRating, 0},
		{cv.Cookbook, jsonData["cookbook"]},
		{cv.File, jsonData["file"]},
		{cv.Dependencies["dmg"], "~> 2.2"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
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
	for _, i := range [][]interface{}{
		{cv.Endpoint, "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0"},
		{cv.License, "Apache v2.0"},
		{cv.TarballFileSize, 5913},
		{cv.Version, "2.0.0"},
		{cv.AverageRating, 0},
		{cv.Cookbook, "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk"},
		{cv.File, "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0/download"},
		{len(cv.Dependencies), 1},
		{cv.Dependencies["dmg"], "~> 2.2"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func Test_NewCookbookVersion_1_EmptyStruct(t *testing.T) {
	cv := NewCookbookVersion()
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

func Test_Equals_2_DifferentEndpoints(t *testing.T) {
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

func Test_Equals_3_DifferentLicense(t *testing.T) {
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

func Test_Equals_4_DifferentFileSize(t *testing.T) {
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

func Test_Equals_5_DifferentDependencies(t *testing.T) {
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

func Test_Diff_1_Equal(t *testing.T) {
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

func Test_Diff_2_DataAddedAndDeleted(t *testing.T) {
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
