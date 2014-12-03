package goulash

import (
	"net/http"
	"testing"
)

func cdata() (data Cookbook) {
	data = Cookbook{
		Component:         Component{Endpoint: "https://example1.com"},
		Name:              "test1",
		Maintainer:        "someuser",
		Description:       "A cookbook",
		Category:          "Other",
		LatestVersion:     "1.2.3",
		ExternalURL:       "https://extexample1.com",
		AverageRating:     0,
		CreatedAt:         "2014-09-01T01:01:01.123Z",
		UpdatedAt:         "2014-09-02T01:01:01.123Z",
		Deprecated:        false,
		FoodcriticFailure: false,
		Versions:          []string{"1.2.3", "1.2.0", "1.1.0"},
		Metrics: Metrics{
			Downloads: Downloads{
				Total: 99,
				Versions: map[string]int{
					"1.2.3": 32,
					"1.2.0": 33,
					"1.1.0": 34,
				},
			},
			Followers: 123,
		},
	}
	return
}

var cjsonData = map[string]string{
	"name":               "chef-dk",
	"maintainer":         "roboticcheese",
	"description":        "Installs/configures the Chef-DK",
	"category":           "Other",
	"latest_version":     "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1",
	"external_url":       "https://github.com/RoboticCheese/chef-dk-chef",
	"average_rating":     "null",
	"created_at":         "2014-06-24T01:14:49.000Z",
	"updated_at":         "2014-09-20T04:46:00.780Z",
	"deprecated":         "false",
	"foodcritic_failure": "false",
	"versions": `
		[ "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1",
		 "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0"]`,
	"metrics": `{
		"downloads": {
			"total": 100,
			"versions": {
				"2.0.0": 50,
				"2.0.1": 50
			}
		},
		"followers": 20
	}`,
}

func cjsonified() (res string) {
	res = `{"name": "` + cjsonData["name"] + `",` +
		`"maintainer": "` + cjsonData["maintainer"] + `",` +
		`"description": "` + cjsonData["description"] + `",` +
		`"category": "` + cjsonData["category"] + `",` +
		`"latest_version": "` + cjsonData["latest_version"] + `",` +
		`"external_url": "` + cjsonData["external_url"] + `",` +
		`"average_rating": ` + cjsonData["average_rating"] + `,` +
		`"created_at": "` + cjsonData["created_at"] + `",` +
		`"updated_at": "` + cjsonData["updated_at"] + `",` +
		`"deprecated": ` + cjsonData["deprecated"] + `,` +
		`"foodcritic_failure": ` + cjsonData["foodcritic_failure"] + `,` +
		`"versions": ` + cjsonData["versions"] + `,` +
		`"metrics": ` + cjsonData["metrics"] + `}`
	return
}

func TestNewCookbookNoError(t *testing.T) {
	ts := StartHTTP(cjsonified(), nil)
	defer ts.Close()

	i := new(APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	c, err := NewCookbook(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for _, i := range [][]interface{}{
		{c.Endpoint, ts.URL + "/api/v1/cookbooks/chef-dk"},
		{c.Name, cjsonData["name"]},
		{c.Maintainer, cjsonData["maintainer"]},
		{c.Description, cjsonData["description"]},
		{c.Category, cjsonData["category"]},
		{c.LatestVersion, cjsonData["latest_version"]},
		{c.ExternalURL, cjsonData["external_url"]},
		{c.CreatedAt, cjsonData["created_at"]},
		{c.UpdatedAt, cjsonData["updated_at"]},
		{c.Deprecated, false},
		{c.FoodcriticFailure, false},
		{c.AverageRating, 0},
		{len(c.Versions), 2},
		{c.Versions[0], "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1"},
		{c.Versions[1], "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0"},
		{c.Metrics.Downloads.Total, 100},
		{c.Metrics.Downloads.Versions["2.0.0"], 50},
		{c.Metrics.Downloads.Versions["2.0.1"], 50},
		{c.Metrics.Followers, 20},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestNewCookbookNilFoodcriticFailure(t *testing.T) {
	cjsonData["foodcritic_failure"] = "null"
	ts := StartHTTP(cjsonified(), nil)
	defer ts.Close()

	i := new(APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	c, err := NewCookbook(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if c.FoodcriticFailure != false {
		t.Fatalf("Expected: nil, got: %v", c.FoodcriticFailure)
	}
}

func TestNewCookbookAverageRating(t *testing.T) {
	cjsonData["average_rating"] = "20"
	ts := StartHTTP(cjsonified(), nil)
	defer ts.Close()

	i := new(APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	c, err := NewCookbook(i, "chef-dk")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if c.AverageRating != 20 {
		t.Fatalf("Expected: 20, got: %v", c.AverageRating)
	}
}

func TestNewCookbookConnError(t *testing.T) {
	ts := StartHTTP(cjsonified(), nil)
	ts.Close()

	i := new(APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	_, err := NewCookbook(i, "chef-dk")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewCookbook404Error(t *testing.T) {
	ts := StartHTTP(http.NotFound, nil)
	defer ts.Close()

	i := new(APIInstance)
	i.Endpoint = ts.URL + "/api/v1"
	_, err := NewCookbook(i, "chef-dk")
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewCookbookRealData(t *testing.T) {
	i := new(APIInstance)
	i.Endpoint = "https://supermarket.getchef.com/api/v1"
	c, err := NewCookbook(i, "chef-dk")
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

func TestInitCookbookEmptyStruct(t *testing.T) {
	c := InitCookbook()
	for _, i := range [][]interface{}{
		{c.Endpoint, ""},
		{c.Name, ""},
		{c.Maintainer, ""},
		{c.Description, ""},
		{c.Category, ""},
		{c.LatestVersion, ""},
		{c.ExternalURL, ""},
		{c.AverageRating, 0},
		{c.CreatedAt, ""},
		{c.UpdatedAt, ""},
		{c.Deprecated, false},
		{c.FoodcriticFailure, false},
		{len(c.Versions), 0},
		{c.Metrics.Downloads.Total, 0},
		{len(c.Metrics.Downloads.Versions), 0},
		{c.Metrics.Followers, 0},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookEmptyEmpty(t *testing.T) {
	c := InitCookbook()
	res := c.Empty()
	if res != true {
		t.Fatalf("Expected: true, got: %v", res)
	}
}

func TestCookbookEmptyHasName(t *testing.T) {
	c := InitCookbook()
	c.Name = "thing"
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func TestCookbookEmptyHasRating(t *testing.T) {
	c := InitCookbook()
	c.AverageRating = 10
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func TestCookbookEmptyIsDeprecated(t *testing.T) {
	c := InitCookbook()
	c.Deprecated = true
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func TestCookbookEmptyHasVersions(t *testing.T) {
	c := InitCookbook()
	c.Versions = []string{"0.1.0"}
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func TestCookbookEmptyHasFollowers(t *testing.T) {
	c := InitCookbook()
	c.Metrics.Followers = 20
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func TestCookbookEmptyHasDownloads(t *testing.T) {
	c := InitCookbook()
	c.Metrics.Downloads.Total = 20
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func TestCookbookEmptyHasVersionDownloads(t *testing.T) {
	c := InitCookbook()
	c.Metrics.Downloads.Versions["0.1.0"] = 20
	res := c.Empty()
	if res != false {
		t.Fatalf("Expected: false, got: %v", res)
	}
}

func TestCookbookEqualsEqual(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	res1 := data1.Equals(&data2)
	res2 := data2.Equals(&data1)
	for _, i := range [][]bool{
		{res1, true},
		{res2, true},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookEqualsDifferentEndpoints(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Endpoint = "https://somewherelse.com"
	res1 := data1.Equals(&data2)
	res2 := data2.Equals(&data1)
	for _, i := range [][]bool{
		{res1, false},
		{res2, false},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookEqualsDifferentName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "ansible"
	res1 := data1.Equals(&data2)
	res2 := data2.Equals(&data1)
	for _, i := range [][]bool{
		{res1, false},
		{res2, false},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookEqualsDifferentLatestVersion(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.LatestVersion = "9.9.9"
	res1 := data1.Equals(&data2)
	res2 := data2.Equals(&data1)
	for _, i := range [][]bool{
		{res1, false},
		{res2, false},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookEqualsDifferentVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions = append(data2.Versions, "9.9.9")
	res1 := data1.Equals(&data2)
	res2 := data2.Equals(&data1)
	for _, i := range [][]bool{
		{res1, false},
		{res2, false},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookEqualsDifferentMetrics(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Metrics.Downloads.Versions["1.2.3"] = 999
	res1 := data1.Equals(&data2)
	res2 := data2.Equals(&data1)
	for _, i := range [][]bool{
		{res1, false},
		{res2, false},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffEqual(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	for _, i := range []*Cookbook{
		pos1,
		neg1,
		pos2,
		neg2,
	} {
		if i != nil {
			t.Fatalf("Expected: nil, got: %v", i)
		}
	}
}

func TestCookbookDiffDifferentEndpoints(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Endpoint = "https://somewherelse.com"
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	for _, i := range [][]string{
		{pos1.Name, ""},
		{pos1.Endpoint, "https://somewherelse.com"},
		{neg1.Name, ""},
		{neg1.Endpoint, "https://example1.com"},
		{pos2.Name, ""},
		{pos2.Endpoint, "https://example1.com"},
		{neg2.Name, ""},
		{neg2.Endpoint, "https://somewherelse.com"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffDifferentName(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Name = "ansible"
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	for _, i := range [][]string{
		{pos1.Name, "ansible"},
		{pos1.Endpoint, ""},
		{neg1.Name, "test1"},
		{neg1.Endpoint, ""},
		{pos2.Name, "test1"},
		{pos2.Endpoint, ""},
		{neg2.Name, "ansible"},
		{neg2.Endpoint, ""},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffDifferentRating(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.AverageRating = 99
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	if neg1 != nil {
		t.Fatalf("Expected: nil, got: %v", neg1)
	}
	if pos2 != nil {
		t.Fatalf("Expected: nil, got: %v", pos2)
	}
	for _, i := range [][]interface{}{
		{pos1.AverageRating, 99},
		{neg2.AverageRating, 99},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffDifferentDeprecatedStatus(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Deprecated = true
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	if neg1 != nil {
		t.Fatalf("Expected: nil, got: %v", neg1)
	}
	if pos2 != nil {
		t.Fatalf("Expected: nil, got: %v", pos2)
	}
	for _, i := range [][]interface{}{
		{pos1.Deprecated, true},
		{neg2.Deprecated, true},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffDifferentVersions(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Versions = []string{"1.2.3", "1.1.0", "9.9.9"}
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	for _, i := range [][]interface{}{
		{len(pos1.Versions), 1},
		{pos1.Versions[0], "9.9.9"},
		{len(neg1.Versions), 1},
		{neg1.Versions[0], "1.2.0"},
		{len(pos2.Versions), 1},
		{pos2.Versions[0], "1.2.0"},
		{len(neg2.Versions), 1},
		{neg2.Versions[0], "9.9.9"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestCookbookDiffDifferentMetrics(t *testing.T) {
	data1 := cdata()
	data2 := cdata()
	data2.Metrics.Downloads.Versions = map[string]int{
		"1.2.3": 32,
		"1.1.0": 999,
		"2.0.0": 64,
	}
	data2.Metrics.Followers = 456
	pos1, neg1 := data1.Diff(&data2)
	pos2, neg2 := data2.Diff(&data1)
	for _, i := range [][]interface{}{
		{len(pos1.Metrics.Downloads.Versions), 2},
		{pos1.Metrics.Downloads.Versions["1.1.0"], 999},
		{pos1.Metrics.Downloads.Versions["2.0.0"], 64},
		{pos1.Metrics.Followers, 456},
		{len(neg1.Metrics.Downloads.Versions), 2},
		{neg1.Metrics.Downloads.Versions["1.2.0"], 33},
		{neg1.Metrics.Downloads.Versions["1.1.0"], 34},
		{neg1.Metrics.Followers, 123},
		{len(pos2.Metrics.Downloads.Versions), 2},
		{pos2.Metrics.Downloads.Versions["1.1.0"], 34},
		{pos2.Metrics.Downloads.Versions["1.2.0"], 33},
		{pos2.Metrics.Followers, 123},
		{len(neg2.Metrics.Downloads.Versions), 2},
		{neg2.Metrics.Downloads.Versions["1.1.0"], 999},
		{neg2.Metrics.Downloads.Versions["2.0.0"], 64},
		{neg2.Metrics.Followers, 456},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}
