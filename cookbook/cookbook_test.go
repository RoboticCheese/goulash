package cookbook

import (
	"fmt"
	"github.com/RoboticCheese/goulash/api_instance"
	"net/http"
	"net/http/httptest"
	"testing"
)

var name = "chef-dk"
var maintainer = "roboticcheese"
var description = "Installs/configures the Chef-DK"
var category = "Other"
var latest_version = "https://supermarket.getchef.com/api/v1/cookbooks/" +
	"chef-dk/versions/2.0.1"
var external_url = "https://github.com/RoboticCheese/chef-dk-chef"
var average_rating = "null"
var created_at = "2014-06-24T01:14:49.000Z"
var updated_at = "2014-09-20T04:46:00.780Z"
var deprecated = "false"
var foodcritic_failure = "false"
var versions = `
	[
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0"
	]
`
var metrics = `
	{
		"downloads": {
			"total": 100,
			"versions": {
				"2.0.0": 50,
				"2.0.1": 50
			}
		},
		"followers": 20
	}
`

func jsonified() (res string) {
	res = `{"name": "` + name + `",` +
		`"maintainer": "` + maintainer + `",` +
		`"description": "` + description + `",` +
		`"category": "` + category + `",` +
		`"latest_version": "` + latest_version + `",` +
		`"external_url": "` + external_url + `",` +
		`"average_rating": ` + average_rating + `,` +
		`"created_at": "` + created_at + `",` +
		`"updated_at": "` + updated_at + `",` +
		`"deprecated": ` + deprecated + `,` +
		`"foodcritic_failure": ` + foodcritic_failure + `,` +
		`"versions": ` + versions + `,` +
		`"metrics": ` + metrics + `}`
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

	i, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		c.Endpoint:    ts.URL + "/api/v1/cookbooks/chef-dk",
		c.Name:        "chef-dk",
		c.Maintainer:  "roboticcheese",
		c.Description: "Installs/configures the Chef-DK",
		c.Category:    "Other",
		c.LatestVersion: "https://supermarket.getchef.com/api/v1/" +
			"cookbooks/chef-dk/versions/2.0.1",
		c.ExternalURL: "https://github.com/RoboticCheese/chef-dk-chef",
		c.CreatedAt:   "2014-06-24T01:14:49.000Z",
		c.UpdatedAt:   "2014-09-20T04:46:00.780Z",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
	if c.Deprecated != false {
		t.Fatalf("Expected: false, got: %v", c.Deprecated)
	}
	if c.FoodcriticFailure != false {
		t.Fatalf("Expected: false, got: %v", c.FoodcriticFailure)
	}
	if c.AverageRating != 0 {
		t.Fatalf("Expected: 0, got: %v", c.AverageRating)
	}
	if len(c.Versions) != 2 {
		t.Fatalf("Expected: 2 versions, got: %v", len(c.Versions))
	}
	ver := "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1"
	if c.Versions[0] != ver {
		t.Fatalf("Expected: %v, got: %v", ver, c.Versions[0])
	}
	ver = "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0"
	if c.Versions[1] != ver {
		t.Fatalf("Expected: %v, got: %v", ver, c.Versions[1])
	}
	if c.Metrics.Downloads.Total != 100 {
		t.Fatalf("Expected: 100, got: %v", c.Metrics.Downloads.Total)
	}
	if c.Metrics.Downloads.Versions["2.0.0"] != 50 {
		t.Fatalf("Expected: 50, got: %v", c.Metrics.Downloads.Versions["2.0.0"])
	}
	if c.Metrics.Downloads.Versions["2.0.1"] != 50 {
		t.Fatalf("Expected: 50, got: %v", c.Metrics.Downloads.Versions["2.0.1"])
	}
	if c.Metrics.Followers != 20 {
		t.Fatalf("Expected: 20, got: %v", c.Metrics.Followers)
	}
}

func Test_New_2_NilFoodcriticFailure(t *testing.T) {
	foodcritic_failure = "null"
	ts := start_http()
	defer ts.Close()

	i, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if c.FoodcriticFailure != false {
		t.Fatalf("Expected: nil, got: %v", c.FoodcriticFailure)
	}
}

func Test_New_3_AverageRating(t *testing.T) {
	average_rating = "20"
	ts := start_http()
	defer ts.Close()

	i, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if c.AverageRating != 20 {
		t.Fatalf("Expected: 20, got: %v", c.AverageRating)
	}
}

func Test_New_4_ConnError(t *testing.T) {
	ts := start_http()
	ts.Close()

	i, err := api_instance.New(ts.URL)
	_, err = New(i, "chef-dk")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_5_404Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer ts.Close()

	i, err := api_instance.New(ts.URL)
	_, err = New(i, "chef-dk")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_6_RealData(t *testing.T) {
	i, err := api_instance.New("https://supermarket.getchef.com")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	c, err := New(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		c.Name:        "chef-dk",
		c.Maintainer:  "roboticcheese",
		c.Category:    "Other",
		c.ExternalURL: "https://github.com/RoboticCheese/chef-dk-chef",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
}
