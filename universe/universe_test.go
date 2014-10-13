package universe

import (
	"fmt"
	"github.com/RoboticCheese/goulash/api_instance"
	"github.com/RoboticCheese/goulash/common"
	"net/http"
	"net/http/httptest"
	"testing"
)

func udata() (data *Universe) {
	data = &Universe{
		Component:   common.Component{Endpoint: "https://example.com"},
		APIInstance: &api_instance.APIInstance{},
		Cookbooks: map[string]*Cookbook{
			"test1": &Cookbook{
				Name: "test1",
				Versions: map[string]*CookbookVersion{
					"0.1.0": &CookbookVersion{
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

func json_data() (json_data map[string]map[string]map[string]string) {
	json_data = map[string]map[string]map[string]string{
		"chef": {
			"0.12.0": {
				"location_type": "opscode",
				"location_path": "https://supermarket.getchef.com/api/v1",
				"download_url":  "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.12.0/download",
				"dependencies":  `{"runit": ">= 0.0.0","couchdb": ">= 0.0.0"}`,
			},
			"0.20.0": {
				"location_type": "opscode",
				"location_path": "https://supermarket.getchef.com/api/v1",
				"download_url":  "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.20.0/download",
				"dependencies":  `{"zlib": ">= 0.0.0","xml": ">= 0.0.0"}`,
			},
		},
		"djbdns": {
			"0.7.0": {
				"location_type": "opscode",
				"location_path": "https://supermarket.getchef.com/api/v1",
				"download_url":  "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.7.0/download",
				"dependencies":  `{"runit": ">= 0.0.0","build-essential": ">= 0.0.0"}`,
			},
			"0.8.2": {
				"location_type": "opscode",
				"location_path": "https://supermarket.getchef.com/api/v1",
				"download_url":  "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.8.2/download",
				"dependencies":  `{"runit": ">= 0.0.0","build-essential": ">= 0.0.0"}`,
			},
		},
	}
	return
}

func http_headers() (res map[string]string) {
	res = map[string]string{}
	return
}

func http_body(json_data map[string]map[string]map[string]string) (res string) {
	res = `
		{"chef": {"0.12.0": {` +
		`"location_type": "` + json_data["chef"]["0.12.0"]["location_type"] + `",` +
		`"location_path": "` + json_data["chef"]["0.12.0"]["location_path"] + `",` +
		`"download_url": "` + json_data["chef"]["0.12.0"]["download_url"] + `",` +
		`"dependencies": ` + json_data["chef"]["0.12.0"]["dependencies"] +
		`}, "0.20.0": {` +
		`"location_type": "` + json_data["chef"]["0.20.0"]["location_type"] + `",` +
		`"location_path": "` + json_data["chef"]["0.20.0"]["location_path"] + `",` +
		`"download_url": "` + json_data["chef"]["0.20.0"]["download_url"] + `",` +
		`"dependencies": ` + json_data["chef"]["0.20.0"]["dependencies"] +
		`}}, "djbdns": {"0.7.0": {` +
		`"location_type": "` + json_data["djbdns"]["0.7.0"]["location_type"] + `",` +
		`"location_path": "` + json_data["djbdns"]["0.7.0"]["location_path"] + `",` +
		`"download_url": "` + json_data["djbdns"]["0.7.0"]["download_url"] + `",` +
		`"dependencies": ` + json_data["djbdns"]["0.7.0"]["dependencies"] +
		`}, "0.8.2": {` +
		`"location_type": "` + json_data["djbdns"]["0.8.2"]["location_type"] + `",` +
		`"location_path": "` + json_data["djbdns"]["0.8.2"]["location_path"] + `",` +
		`"download_url": "` + json_data["djbdns"]["0.8.2"]["download_url"] + `",` +
		`"dependencies": ` + json_data["djbdns"]["0.8.2"]["dependencies"] +
		`}}}`
	return
}

func start_http(http_headers func() map[string]string, json_data func() map[string]map[string]map[string]string) (ts *httptest.Server) {
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				for k, v := range http_headers() {
					w.Header().Set(k, v)
				}
				fmt.Fprint(w, http_body(json_data()))
			},
		),
	)
	return
}

func Test_New_1_NoError(t *testing.T) {
	ts := start_http(http_headers, json_data)
	defer ts.Close()

	i := new(api_instance.APIInstance)
	i.BaseURL = ts.URL
	u, err := New(i)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(u.Cookbooks) != 2 {
		t.Fatalf("Expected 2 cookbooks, got: %v", len(u.Cookbooks))
	}
	json_data := json_data()
	for k, v := range map[string]string{
		u.Endpoint:                                                              ts.URL + "/universe",
		u.Cookbooks["chef"].Name:                                                "chef",
		u.Cookbooks["chef"].Versions["0.12.0"].Version:                          "0.12.0",
		u.Cookbooks["chef"].Versions["0.12.0"].LocationType:                     json_data["chef"]["0.12.0"]["location_type"],
		u.Cookbooks["chef"].Versions["0.12.0"].LocationPath:                     json_data["chef"]["0.12.0"]["location_path"],
		u.Cookbooks["chef"].Versions["0.12.0"].DownloadURL:                      json_data["chef"]["0.12.0"]["download_url"],
		u.Cookbooks["chef"].Versions["0.12.0"].Dependencies["runit"]:            ">= 0.0.0",
		u.Cookbooks["chef"].Versions["0.12.0"].Dependencies["couchdb"]:          ">= 0.0.0",
		u.Cookbooks["chef"].Versions["0.20.0"].Version:                          "0.20.0",
		u.Cookbooks["chef"].Versions["0.20.0"].LocationType:                     json_data["chef"]["0.20.0"]["location_type"],
		u.Cookbooks["chef"].Versions["0.20.0"].LocationPath:                     json_data["chef"]["0.20.0"]["location_path"],
		u.Cookbooks["chef"].Versions["0.20.0"].DownloadURL:                      json_data["chef"]["0.20.0"]["download_url"],
		u.Cookbooks["chef"].Versions["0.20.0"].Dependencies["zlib"]:             ">= 0.0.0",
		u.Cookbooks["chef"].Versions["0.20.0"].Dependencies["xml"]:              ">= 0.0.0",
		u.Cookbooks["djbdns"].Name:                                              "djbdns",
		u.Cookbooks["djbdns"].Versions["0.7.0"].Version:                         "0.7.0",
		u.Cookbooks["djbdns"].Versions["0.7.0"].LocationType:                    json_data["djbdns"]["0.7.0"]["location_type"],
		u.Cookbooks["djbdns"].Versions["0.7.0"].LocationPath:                    json_data["djbdns"]["0.7.0"]["location_path"],
		u.Cookbooks["djbdns"].Versions["0.7.0"].DownloadURL:                     json_data["djbdns"]["0.7.0"]["download_url"],
		u.Cookbooks["djbdns"].Versions["0.7.0"].Dependencies["runit"]:           ">= 0.0.0",
		u.Cookbooks["djbdns"].Versions["0.7.0"].Dependencies["build-essential"]: ">= 0.0.0",
		u.Cookbooks["djbdns"].Versions["0.8.2"].Version:                         "0.8.2",
		u.Cookbooks["djbdns"].Versions["0.8.2"].LocationType:                    json_data["djbdns"]["0.8.2"]["location_type"],
		u.Cookbooks["djbdns"].Versions["0.8.2"].LocationPath:                    json_data["djbdns"]["0.8.2"]["location_path"],
		u.Cookbooks["djbdns"].Versions["0.8.2"].DownloadURL:                     json_data["djbdns"]["0.8.2"]["download_url"],
		u.Cookbooks["djbdns"].Versions["0.8.2"].Dependencies["runit"]:           ">= 0.0.0",
		u.Cookbooks["djbdns"].Versions["0.8.2"].Dependencies["build-essential"]: ">= 0.0.0",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
}

func Test_New_2_ConnError(t *testing.T) {
	ts := start_http(http_headers, json_data)
	ts.Close()

	i := new(api_instance.APIInstance)
	i.BaseURL = ts.URL
	_, err := New(i)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_3_404Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer ts.Close()

	i := new(api_instance.APIInstance)
	i.BaseURL = ts.URL
	_, err := New(i)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_4_RealData(t *testing.T) {
	i := new(api_instance.APIInstance)
	i.BaseURL = "https://supermarket.getchef.com"
	u, err := New(i)
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

func Test_Empty_1_Empty(t *testing.T) {
	u := new(Universe)
	res := u.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Empty_2_StillEmpty(t *testing.T) {
	u := NewUniverse()
	res := u.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Empty_3_HasEmptyCookbooks(t *testing.T) {
	u := NewUniverse()
	u.Cookbooks["thing1"] = NewCookbook()
	res := u.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Empty_4_HasEndpoint(t *testing.T) {
	u := NewUniverse()
	u.Endpoint = "http://example.com"
	res := u.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_5_HasAPIInstance(t *testing.T) {
	u := NewUniverse()
	i := new(api_instance.APIInstance)
	u.APIInstance = i
	res := u.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_6_HasCookbooks(t *testing.T) {
	u := NewUniverse()
	u.Cookbooks["nginx"] = NewCookbook()
	u.Cookbooks["nginx"].Versions["0.1.0"] = NewCookbookVersion()
	u.Cookbooks["nginx"].Versions["0.1.0"].LocationType = "opscode"
	res := u.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Equals_1_Equal(t *testing.T) {
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

// TODO
//func Test_Equals_2_DifferentEndpoints(t *testing.T) {
//	data1 := udata()
//	data2 := udata()
//	data2.Endpoint = "otherexample.com"
//	for _, res := range []bool{
//		data1.Equals(data2),
//		data2.Equals(data1),
//	} {
//		if res != false {
//			t.Fatalf("Expected false, got: %v", res)
//		}
//	}
//}

func Test_Equals_2_MoreCookbooks(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["test2"] = &Cookbook{}
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func Test_Equals_3_FewerCookbooks(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks = map[string]*Cookbook{}
	for _, res := range []bool{
		data1.Equals(data2),
		data2.Equals(data1),
	} {
		if res != false {
			t.Fatalf("Expected false, got: %v", res)
		}
	}
}

func Test_Equals_4_DifferentCookbooks(t *testing.T) {
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

func Test_Update_1_NoChanges(t *testing.T) {
	ts := start_http(http_headers, json_data)
	defer ts.Close()

	a, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := New(a)
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

func Test_Update_2_SomeChanges(t *testing.T) {
	json := json_data()
	json_data := func() map[string]map[string]map[string]string {
		return json
	}

	ts := start_http(http_headers, json_data)
	defer ts.Close()

	a, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := New(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	json["chef"]["0.12.0"]["location_type"] = "elsewhere"
	json["chef"]["0.12.0"]["location_path"] = "https://example.com"

	pos, neg, err := u.Update()
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
	chk := pos.Cookbooks["chef"].Versions["0.12.0"]
	if chk.LocationType != "elsewhere" {
		t.Fatalf("Expected 'elsewhere', got: %v", chk.LocationType)
	}
	if chk.LocationPath != "https://example.com" {
		t.Fatalf("Expected 'https://example.com', got: %v",
			chk.LocationPath)
	}
}

func Test_Update_3_ETagSomeChanges(t *testing.T) {
	headers := http_headers()
	headers["ETag"] = "tag1"
	http_headers := func() map[string]string {
		return headers
	}
	json := json_data()
	json_data := func() map[string]map[string]map[string]string {
		return json
	}

	ts := start_http(http_headers, json_data)
	defer ts.Close()

	a, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := New(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	json["chef"]["0.12.0"]["location_type"] = "elsewhere"
	json["chef"]["0.12.0"]["location_path"] = "https://example.com"
	headers["ETag"] = "tag2"

	pos, neg, err := u.Update()
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
	chk := pos.Cookbooks["chef"].Versions["0.12.0"]
	if chk.LocationType != "elsewhere" {
		t.Fatalf("Expected 'elsewhere', got: %v", chk.LocationType)
	}
	if chk.LocationPath != "https://example.com" {
		t.Fatalf("Expected 'https://example.com', got: %v",
			chk.LocationPath)
	}
}

func Test_Update_4_ETagNoChanges(t *testing.T) {
	headers := http_headers()
	headers["ETag"] = "tag1"
	http_headers := func() map[string]string {
		return headers
	}
	json := json_data()
	json_data := func() map[string]map[string]map[string]string {
		return json
	}

	ts := start_http(http_headers, json_data)
	defer ts.Close()

	a, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := New(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	json["chef"]["0.12.0"]["location_type"] = "elsewhere"
	json["chef"]["0.12.0"]["location_path"] = "https://example.com"

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
}

func Test_Update_5_Error(t *testing.T) {
	ts := start_http(http_headers, json_data)

	a, err := api_instance.New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}
	u, err := New(a)
	if err != nil {
		t.Fatalf("Expected no err, got: %v", err)
	}

	ts.Close()

	_, _, err = u.Update()
	if err == nil {
		t.Fatalf("Expected non-nil, got: %v", err)
	}
}

func Test_Diff_1_Equal(t *testing.T) {
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

func Test_Diff_2_AddedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["nginx"] = &Cookbook{
		Name: "nginx",
		Versions: map[string]*CookbookVersion{
			"0.1.0": &CookbookVersion{
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

func Test_Diff_3_DeletedCookbook(t *testing.T) {
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

func Test_Diff_4_UpdatedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["test1"].Versions["0.1.0"].LocationType = "elsewhere"
	pos, neg := data1.Diff(data2)
	if len(pos.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(pos.Cookbooks))
	}
	if pos.Cookbooks["test1"].Versions["0.1.0"].LocationType != "elsewhere" {
		t.Fatalf("Expected 'elsewhere', got: %v",
			pos.Cookbooks["test1"].Versions["0.1.0"].LocationType)
	}
	if neg != nil {
		t.Fatalf("Expected nil, got: %v", neg)
	}
}

func Test_positiveDiff_1_Equal(t *testing.T) {
	data1 := udata()
	data2 := udata()
	diff := data1.positiveDiff(data2)
	if diff != nil {
		t.Fatalf("Expected nil, got: %v", diff)
	}
}

func Test_positiveDiff_2_AddedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["nginx"] = &Cookbook{
		Name: "nginx",
		Versions: map[string]*CookbookVersion{
			"0.1.0": &CookbookVersion{
				Version:      "0.1.0",
				LocationType: "somewhere",
				LocationPath: "https://example.com/nginx",
				DownloadURL:  "https://example.com/nginx/download",
				Dependencies: map[string]string{"thing1": ">= 0.0.0"},
			},
		},
	}
	diff := data1.positiveDiff(data2)
	if len(diff.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(diff.Cookbooks))
	}
	if diff.Cookbooks["nginx"].Versions["0.1.0"].LocationType != "somewhere" {
		t.Fatalf("Expected 'somewhere', got: %v",
			diff.Cookbooks["nginx"].Versions["0.1.0"].LocationType)
	}
}

func Test_positiveDiff_3_UpdatedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	data2.Cookbooks["test1"].Versions["0.1.0"].LocationType = "elsewhere"
	diff := data1.positiveDiff(data2)
	if len(diff.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(diff.Cookbooks))
	}
	if diff.Cookbooks["test1"].Versions["0.1.0"].LocationType != "elsewhere" {
		t.Fatalf("Expected 'elsewhere', got: %v",
			diff.Cookbooks["test1"].Versions["0.1.0"].LocationType)
	}
}

func Test_negativeDiff_1_Equal(t *testing.T) {
	data1 := udata()
	data2 := udata()
	diff := data1.negativeDiff(data2)
	if diff != nil {
		t.Fatalf("Expected nil, got: %v", diff)
	}
}

func Test_negativeDiff_2_DeletedCookbook(t *testing.T) {
	data1 := udata()
	data2 := udata()
	delete(data2.Cookbooks, "test1")
	diff := data1.negativeDiff(data2)
	if len(diff.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(diff.Cookbooks))
	}
	if diff.Cookbooks["test1"].Versions["0.1.0"].LocationType != "opscode" {
		t.Fatalf("Expected 'opscode', got: %v",
			diff.Cookbooks["test1"].Versions["0.1.0"].LocationType)
	}
}

func Test_negativeDiff_3_DeletedCookbookVersion(t *testing.T) {
	data1 := udata()
	data2 := udata()
	delete(data2.Cookbooks["test1"].Versions, "0.1.0")
	diff := data1.negativeDiff(data2)
	if len(diff.Cookbooks) != 1 {
		t.Fatalf("Expected 1 cookbook, got: %v", len(diff.Cookbooks))
	}
	if diff.Cookbooks["test1"].Name != "" {
		t.Fatalf("Expected empty string, got: %v",
			diff.Cookbooks["test1"].Name)
	}
	if diff.Cookbooks["test1"].Versions["0.1.0"].LocationType != "opscode" {
		t.Fatalf("Expected 'opscode', got: %v",
			diff.Cookbooks["test1"].Versions["0.1.0"].LocationType)
	}
}

func Test_decodeJSON_1(t *testing.T) {
	ts := start_http(http_headers, json_data)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("Expected nil, got: %v", err)
	}

	temp_u := map[string]map[string]*CookbookVersion{}
	err = decodeJSON(resp.Body, &temp_u)
	if err != nil {
		t.Fatalf("Expected nil, got: %v", err)
	}
	if len(temp_u) != 2 {
		t.Fatalf("Expected 2 cookbooks, got: %v", len(temp_u))
	}
}
