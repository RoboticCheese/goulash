package goulash

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/RoboticCheese/goulash/universe"
)

func udata() (data *Universe) {
	data = &Universe{
		Component:   Component{Endpoint: "https://example.com"},
		APIInstance: &APIInstance{},
		Cookbooks: map[string]*universe.Cookbook{
			"test1": &universe.Cookbook{
				Name: "test1",
				Versions: map[string]*universe.CookbookVersion{
					"0.1.0": &universe.CookbookVersion{
						Version:      "0.1.0",
						LocationType: "opscode",
						LocationPath: "https://example.com",
						DownloadURL:  "https://example.com/1",
						Dependencies: map[string]string{
							"thing1": ">= 0.0.0",
							"thing2": ">= 0.0.0",
						},
					},
				},
			},
		},
	}
	return
}

func ujsonData() (jsonData map[string]map[string]*universe.CookbookVersion) {
	jsonData = map[string]map[string]*universe.CookbookVersion{
		"chef": {
			"0.12.0": &universe.CookbookVersion{
				LocationType: "opscode",
				LocationPath: "https://supermarket.getchef.com/api/v1",
				DownloadURL:  "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.12.0/download",
				Dependencies: map[string]string{"runit": ">= 0.0.0", "couchdb": ">= 0.0.0"},
			},
			"0.20.0": &universe.CookbookVersion{
				LocationType: "opscode",
				LocationPath: "https://supermarket.getchef.com/api/v1",
				DownloadURL:  "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.20.0/download",
				Dependencies: map[string]string{"zlib": ">= 0.0.0", "xml": ">= 0.0.0"},
			},
		},
		"djbdns": {
			"0.7.0": &universe.CookbookVersion{
				LocationType: "opscode",
				LocationPath: "https://supermarket.getchef.com/api/v1",
				DownloadURL:  "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.7.0/download",
				Dependencies: map[string]string{"runit": ">= 0.0.0", "build-essential": ">= 0.0.0"},
			},
			"0.8.2": &universe.CookbookVersion{
				LocationType: "opscode",
				LocationPath: "https://supermarket.getchef.com/api/v1",
				DownloadURL:  "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.8.2/download",
				Dependencies: map[string]string{"runit": ">= 0.0.0", "build-essential": ">= 0.0.0"},
			},
		},
	}
	return
}

func uhttpBody(ujsonData map[string]map[string]*universe.CookbookVersion) (res string) {
	bres, _ := json.Marshal(ujsonData)
	res = string(bres)
	return
}

func TestNewUniverseNoError(t *testing.T) {
	ts := StartHTTP(uhttpBody(ujsonData()), nil)
	defer ts.Close()

	i := new(APIInstance)
	i.BaseURL = ts.URL
	u, err := NewUniverse(i)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(u.Cookbooks) != 2 {
		t.Fatalf("Expected 2 cookbooks, got: %v", len(u.Cookbooks))
	}
	for k, v := range map[string]string{
		u.Endpoint:                                                              ts.URL + "/universe",
		u.Cookbooks["chef"].Name:                                                "chef",
		u.Cookbooks["chef"].Versions["0.12.0"].Version:                          "0.12.0",
		u.Cookbooks["chef"].Versions["0.12.0"].LocationType:                     "opscode",
		u.Cookbooks["chef"].Versions["0.12.0"].LocationPath:                     "https://supermarket.getchf.com/api/v1",
		u.Cookbooks["chef"].Versions["0.12.0"].DownloadURL:                      "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.12.0/download",
		u.Cookbooks["chef"].Versions["0.12.0"].Dependencies["runit"]:            ">= 0.0.0",
		u.Cookbooks["chef"].Versions["0.12.0"].Dependencies["couchdb"]:          ">= 0.0.0",
		u.Cookbooks["chef"].Versions["0.20.0"].Version:                          "0.20.0",
		u.Cookbooks["chef"].Versions["0.20.0"].LocationType:                     "opscode",
		u.Cookbooks["chef"].Versions["0.20.0"].LocationPath:                     "https://supermarket.getchef.com/api/v1",
		u.Cookbooks["chef"].Versions["0.20.0"].DownloadURL:                      "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.20.0/download",
		u.Cookbooks["chef"].Versions["0.20.0"].Dependencies["zlib"]:             ">= 0.0.0",
		u.Cookbooks["chef"].Versions["0.20.0"].Dependencies["xml"]:              ">= 0.0.0",
		u.Cookbooks["djbdns"].Name:                                              "djbdns",
		u.Cookbooks["djbdns"].Versions["0.7.0"].Version:                         "0.7.0",
		u.Cookbooks["djbdns"].Versions["0.7.0"].LocationType:                    "opscode",
		u.Cookbooks["djbdns"].Versions["0.7.0"].LocationPath:                    "https://supermarket.getchef.com/api/v1",
		u.Cookbooks["djbdns"].Versions["0.7.0"].DownloadURL:                     "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.7.0/download",
		u.Cookbooks["djbdns"].Versions["0.7.0"].Dependencies["runit"]:           ">= 0.0.0",
		u.Cookbooks["djbdns"].Versions["0.7.0"].Dependencies["build-essential"]: ">= 0.0.0",
		u.Cookbooks["djbdns"].Versions["0.8.2"].Version:                         "0.8.2",
		u.Cookbooks["djbdns"].Versions["0.8.2"].LocationType:                    "opscode",
		u.Cookbooks["djbdns"].Versions["0.8.2"].LocationPath:                    "https://supermarket.getchef.com/api/v1",
		u.Cookbooks["djbdns"].Versions["0.8.2"].DownloadURL:                     "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.8.2/download",
		u.Cookbooks["djbdns"].Versions["0.8.2"].Dependencies["runit"]:           ">= 0.0.0",
		u.Cookbooks["djbdns"].Versions["0.8.2"].Dependencies["build-essential"]: ">= 0.0.0",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
}

func TestNewUniverseConnError(t *testing.T) {
	ts := StartHTTP(uhttpBody(ujsonData()), nil)
	ts.Close()

	i := new(APIInstance)
	i.BaseURL = ts.URL
	_, err := NewUniverse(i)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewUniverse404Error(t *testing.T) {
	ts := StartHTTP(http.NotFound, nil)
	defer ts.Close()

	i := new(APIInstance)
	i.BaseURL = ts.URL
	_, err := NewUniverse(i)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewUniverseRealData(t *testing.T) {
	i := new(APIInstance)
	i.BaseURL = "https://supermarket.getchef.com"
	u, err := NewUniverse(i)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	for k, v := range map[string]string{
		u.Cookbooks["chef-dk"].Name:                                  "chef-dk",
		u.Cookbooks["chef-dk"].Versions["2.0.0"].Version:             "2.0.0",
		u.Cookbooks["chef-dk"].Versions["2.0.0"].LocationType:        "opscode",
		u.Cookbooks["chef-dk"].Versions["2.0.0"].LocationPath:        "https://supermarket.getchef.com/api/v1",
		u.Cookbooks["chef-dk"].Versions["2.0.0"].DownloadURL:         "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
		u.Cookbooks["chef-dk"].Versions["2.0.0"].Dependencies["dmg"]: "~> 2.2",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
}

func TestUniverseEmptyEmpty(t *testing.T) {
	u := new(Universe)
	res := u.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestUniverseEmptyStillEmpty(t *testing.T) {
	u := InitUniverse()
	res := u.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestUniverseEmptyHasEmptyCookbooks(t *testing.T) {
	u := InitUniverse()
	u.Cookbooks["thing1"] = universe.NewCookbook()
	res := u.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestUniverseEmptyHasEndpoint(t *testing.T) {
	u := InitUniverse()
	u.Endpoint = "http://example.com"
	res := u.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestUniverseEmptyHasAPIInstance(t *testing.T) {
	u := InitUniverse()
	i := new(APIInstance)
	i.BaseURL = "https://example.com"
	u.APIInstance = i
	res := u.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestUniverseEmptyHasCookbooks(t *testing.T) {
	u := InitUniverse()
	u.Cookbooks["nginx"] = universe.NewCookbook()
	u.Cookbooks["nginx"].Versions["0.1.0"] = universe.NewCookbookVersion()
	u.Cookbooks["nginx"].Versions["0.1.0"].LocationType = "opscode"
	res := u.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestUniverseEqualsEqual(t *testing.T) {
	data1 := udata()
	data2 := udata()
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != true {
			t.Fatalf("Expected true, got: %v", res)
		}
	}
}

func TestUniverseEqualsMoreCookbooks(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["test2"] = &universe.Cookbook{}
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestUniverseEqualsFewerCookbooks(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks = map[string]*universe.Cookbook{}
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestUniverseEqualsDifferentCookbooks(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["test1"].Versions["0.1.0"].LocationType = "other"
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func TestUniverseUpdateNoChanges(t *testing.T) {
	ts := StartHTTP(uhttpBody(ujsonData()), nil)
	defer ts.Close()

	a, err := NewAPIInstance(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := NewUniverse(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	pos, neg, err := u.Update()
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
}

func TestUniverseUpdateSomeChanges(t *testing.T) {
	data := ujsonData()
	body := func() string {
		res, _ := json.Marshal(data)
		return string(res)
	}

	ts := StartHTTP(body, nil)
	defer ts.Close()

	a, err := NewAPIInstance(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := NewUniverse(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	data["chef"]["0.12.0"].LocationType = "elsewhere"
	data["chef"]["0.12.0"].LocationPath = "https://example.com"

	pos, neg, err := u.Update()
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	for _, i := range [][]string{
		{u.Cookbooks["chef"].Versions["0.12.0"].LocationType, "elsewhere"},
		{u.Cookbooks["chef"].Versions["0.12.0"].LocationPath, "https://example.com"},
		{pos.Cookbooks["chef"].Versions["0.12.0"].LocationType, "elsewhere"},
		{pos.Cookbooks["chef"].Versions["0.12.0"].LocationPath, "https://example.com"},
		{neg.Cookbooks["chef"].Versions["0.12.0"].LocationType, "opscode"},
		{neg.Cookbooks["chef"].Versions["0.12.0"].LocationPath, "https://supermarket.getchef.com/api/v1"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}

func TestUniverseUpdateNewVersionReleased(t *testing.T) {
	data := ujsonData()
	body := func() string {
		res, _ := json.Marshal(data)
		return string(res)
	}

	ts := StartHTTP(body, nil)
	defer ts.Close()

	a, err := NewAPIInstance(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := NewUniverse(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	data["chef"]["9.9.9"] = &universe.CookbookVersion{
		LocationType: "opsplode",
		LocationPath: "https://example.com",
		DownloadURL:  "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/9.9.9/download",
		Dependencies: map[string]string{"otherthing": ">= 0.0.0"},
	}
	pos, neg, err := u.Update()
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
	if len(pos.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(pos.Cookbooks))
	}
	chk := pos.Cookbooks["chef"].Versions["9.9.9"]
	if chk == nil {
		t.Fatalf("Expected non-nil, got: %v", chk)
	}
	if chk.LocationType != "opsplode" {
		t.Fatalf("Expected 'opsplode', got: %v", chk.LocationType)
	}
	if chk.LocationPath != "https://example.com" {
		t.Fatalf("Expected 'https://example.com', got: %v",
			chk.LocationPath)
	}
	chk = u.Cookbooks["chef"].Versions["9.9.9"]
	if chk == nil {
		t.Fatalf("Expected non-nil, got: %v", chk)
	}
	if chk.LocationType != "opsplode" {
		t.Fatalf("Expected 'opsplode', got: %v", chk.LocationType)
	}
	if chk.LocationPath != "https://example.com" {
		t.Fatalf("Expected 'https://example.com', got: %v",
			chk.LocationPath)
	}
}

func TestUniverseUpdateETagSomeChanges(t *testing.T) {
	headers := map[string]string{"ETag": "tag1"}
	data := ujsonData()
	body := func() string {
		res, _ := json.Marshal(data)
		return string(res)
	}

	ts := StartHTTP(body, func() map[string]string { return headers })
	defer ts.Close()

	a, err := NewAPIInstance(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := NewUniverse(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	data["chef"]["0.12.0"].LocationType = "elsewhere"
	data["chef"]["0.12.0"].LocationPath = "https://example.com"
	headers["ETag"] = "tag2"

	pos, neg, err := u.Update()
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	for _, i := range [][]string{
		{u.Cookbooks["chef"].Versions["0.12.0"].LocationType, "elsewhere"},
		{u.Cookbooks["chef"].Versions["0.12.0"].LocationPath, "https://example.com"},
		{pos.Cookbooks["chef"].Versions["0.12.0"].LocationType, "elsewhere"},
		{pos.Cookbooks["chef"].Versions["0.12.0"].LocationPath, "https://example.com"},
		{neg.Cookbooks["chef"].Versions["0.12.0"].LocationType, "opscode"},
		{neg.Cookbooks["chef"].Versions["0.12.0"].LocationPath, "https://supermarket.getchef.com/api/v1"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}

func TestUniverseUpdateETagNoChanges(t *testing.T) {
	json := ujsonData()
	ujsonData := func() map[string]map[string]*universe.CookbookVersion {
		return json
	}
	ts := StartHTTP(uhttpBody(ujsonData()), map[string]string{"ETag": "tag1"})
	defer ts.Close()

	a, err := NewAPIInstance(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := NewUniverse(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	json["chef"]["0.12.0"].LocationType = "elsewhere"
	json["chef"]["0.12.0"].LocationPath = "https://example.com"

	pos, neg, err := u.Update()
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	chk := u.Cookbooks["chef"].Versions["0.12.0"].LocationType
	if chk != "opscode" {
		t.Fatalf("Expected 'opscode', got: %v", chk)
	}
}

func TestUniverseUpdateError(t *testing.T) {
	ts := StartHTTP(uhttpBody(ujsonData()), nil)

	a, err := NewAPIInstance(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := NewUniverse(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	ts.Close()

	_, _, err = u.Update()
	if err == nil {
		t.Fatalf("Expected non-nil, got: %v", err)
	}
}

func TestUniverseDiffEqual(t *testing.T) {
	data1 := udata()
	data2 := udata()
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
}

func TestUniverseDiffAddedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["nginx"] = &universe.Cookbook{
		Name: "nginx",
		Versions: map[string]*universe.CookbookVersion{
			"0.1.0": &universe.CookbookVersion{
				Version:      "0.1.0",
				LocationType: "somewhere",
				LocationPath: "https://example.com/nginx",
				DownloadURL:  "https://example.com/nginx/download",
				Dependencies: map[string]string{"thing1": ">= 0.0.0"},
			},
		},
	}
	pos, neg := data1.Diff(data2)
	if len(pos.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(pos.Cookbooks))
	}
	if pos.Cookbooks["nginx"].Versions["0.1.0"].LocationType != "somewhere" {
		t.Fatalf("Expected 'somewhere', got: %v",
			pos.Cookbooks["nginx"].Versions["0.1.0"].LocationType)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
}

func TestUniverseDiffDeletedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	delete(data2.Cookbooks, "test1")
	pos, neg := data1.Diff(data2)
	if pos != nil {
		t.Fatalf("Expected nil, got: %v", pos)
	}
	if len(neg.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(neg.Cookbooks))
	}
	if neg.Cookbooks["test1"].Versions["0.1.0"].LocationType != "opscode" {
		t.Fatalf("Expected 'somewhere', got: %v",
			neg.Cookbooks["test1"].Versions["0.1.0"].LocationType)
	}
}

func TestUniverseDiffUpdatedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["test1"].Versions["0.1.0"].LocationType = "elsewhere"
	pos, neg := data1.Diff(data2)
	if len(pos.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(pos.Cookbooks))
	}
	if len(neg.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(neg.Cookbooks))
	}
	for _, i := range [][]string{
		{pos.Cookbooks["test1"].Versions["0.1.0"].LocationType, "elsewhere"},
		{neg.Cookbooks["test1"].Versions["0.1.0"].LocationType, "opscode"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected %v, got: %v", i[1], i[0])
		}
	}
}

func TestdecodeUniverseJSON(t *testing.T) {
	ts := StartHTTP(uhttpBody(ujsonData()), nil)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Expected nil, got: %v", err)
	}

	tempU := map[string]map[string]*universe.CookbookVersion{}
	err = decodeUniverseJSON(resp.Body, &tempU)
	if err != nil {
		t.Fatalf("Expected nil, got: %v", err)
	}
	if len(tempU) != 2 {
		t.Fatalf("Expected 2 cookbooks, got: %v", len(tempU))
	}
}
