package universe

import (
	"fmt"
	"github.com/RoboticCheese/goulash/api_instance"
	"net/http"
	"net/http/httptest"
	"testing"
)

var test_data = map[string]map[string]map[string]string{
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

func jsonified() (res string) {
	res = `
		{"chef": {"0.12.0": {` +
		`"location_type": "` + test_data["chef"]["0.12.0"]["location_type"] + `",` +
		`"location_path": "` + test_data["chef"]["0.12.0"]["location_path"] + `",` +
		`"download_url": "` + test_data["chef"]["0.12.0"]["download_url"] + `",` +
		`"dependencies": ` + test_data["chef"]["0.12.0"]["dependencies"] +
		`}, "0.20.0": {` +
		`"location_type": "` + test_data["chef"]["0.20.0"]["location_type"] + `",` +
		`"location_path": "` + test_data["chef"]["0.20.0"]["location_path"] + `",` +
		`"download_url": "` + test_data["chef"]["0.20.0"]["download_url"] + `",` +
		`"dependencies": ` + test_data["chef"]["0.20.0"]["dependencies"] +
		`}}, "djbdns": {"0.7.0": {` +
		`"location_type": "` + test_data["djbdns"]["0.7.0"]["location_type"] + `",` +
		`"location_path": "` + test_data["djbdns"]["0.7.0"]["location_path"] + `",` +
		`"download_url": "` + test_data["djbdns"]["0.7.0"]["download_url"] + `",` +
		`"dependencies": ` + test_data["djbdns"]["0.7.0"]["dependencies"] +
		`}, "0.8.2": {` +
		`"location_type": "` + test_data["djbdns"]["0.8.2"]["location_type"] + `",` +
		`"location_path": "` + test_data["djbdns"]["0.8.2"]["location_path"] + `",` +
		`"download_url": "` + test_data["djbdns"]["0.8.2"]["download_url"] + `",` +
		`"dependencies": ` + test_data["djbdns"]["0.8.2"]["dependencies"] +
		`}}}`
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

//func Test_Diff_1_AddRemoveCookbooks(t *testing.T) {
//	data1 := new(Universe)
//	data1.Endpoint = "https://example1.com"
//	data1.Cookbooks["chef"] = Cookbook{}
//	data1.Cookbooks["chef"]["0.1.0"] = CookbookVersion{
//		"opscode",
//		"https://example1.com",
//		"https://example1.com/1",
//		map[string]string{"thing1": ">= 0.0.0"},
//	}
//	data2 := new(Universe)
//	data2.Endpoint = "https://example2.com"
//	data2.Cookbooks["puppet"] = Cookbook{}
//	data2.Cookbooks["puppet"]["0.0.1"] = CookbookVersion{
//		"opscode",
//		"https://example2.com",
//		"https://example2.com/2",
//		map[string]string{"thing2": ">= 0.0.0"},
//	}
//	pos_diff, neg_diff, err := Diff(data1, data2)
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//	for k, v := range map[int]int{
//		len(pos_diff.Cookbooks):           1,
//		len(neg_diff.Cookbooks):           1,
//		len(pos_diff.Cookbooks["puppet"]): 1,
//		len(neg_diff.Cookbooks["chef"]):   1,
//	} {
//		if k != v {
//			t.Fatalf("Expected: %v, got: %v", v, k)
//		}
//	}
//	for k, v := range map[string]string{
//		pos_diff.Endpoint:                                  "https://example2.com",
//		neg_diff.Endpoint:                                  "",
//		pos_diff.Cookbooks["puppet"]["0.0.1"].LocationPath: "https://example2.com",
//		neg_diff.Cookbooks["chef"]["0.1.0"].LocationPath:   "https://example1.com",
//	} {
//		if k != v {
//			t.Fatalf("Expected: %v, got %v", v, k)
//		}
//	}
//}
//
//func Test_Diff_2_AddRemoveVersions(t *testing.T) {
//	data1 := new(Universe)
//	data1.Endpoint = "https://example1.com"
//	data1.Cookbooks["chef"] = Cookbook{}
//	data1.Cookbooks["chef"]["0.1.0"] = CookbookVersion{
//		"opscode",
//		"https://example1.com",
//		"https://example1.com/1",
//		map[string]string{"thing1": ">= 0.0.0"},
//	}
//	data1.Cookbooks["chef"]["0.2.0"] = CookbookVersion{
//		"opscode",
//		"https://example1.com",
//		"https://example1.com/1",
//		map[string]string{"thing1": ">= 0.0.0"},
//	}
//	data1.Cookbooks["puppet"] = Cookbook{}
//	data1.Cookbooks["puppet"]["0.0.1"] = CookbookVersion{
//		"opscode",
//		"https://example2.com",
//		"https://example2.com/2",
//		map[string]string{"thing2": ">= 0.0.0"},
//	}
//	data2 := new(Universe)
//	data2.Endpoint = "https://example2.com"
//	data2.Cookbooks["chef"] = Cookbook{}
//	data2.Cookbooks["chef"]["0.1.0"] = CookbookVersion{
//		"opscode",
//		"https://example1.com",
//		"https://example1.com/1",
//		map[string]string{"thing1": ">= 0.0.0"},
//	}
//	data2.Cookbooks["puppet"] = Cookbook{}
//	data2.Cookbooks["puppet"]["0.0.1"] = CookbookVersion{
//		"opscode",
//		"https://example2.com",
//		"https://example2.com/2",
//		map[string]string{"thing2": ">= 0.0.0"},
//	}
//	data2.Cookbooks["puppet"]["0.1.0"] = CookbookVersion{
//		"opscode",
//		"https://example2.com",
//		"https://example2.com/2",
//		map[string]string{"thing2": ">= 0.0.0"},
//	}
//
//	pos_diff, neg_diff, err := Diff(data1, data2)
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//	for k, v := range map[int]int{
//		len(pos_diff.Cookbooks):           1,
//		len(neg_diff.Cookbooks):           1,
//		len(pos_diff.Cookbooks["puppet"]): 1,
//		len(neg_diff.Cookbooks["chef"]):   1,
//	} {
//		if k != v {
//			t.Fatalf("Expected: %v, got: %v", v, k)
//		}
//	}
//	for k, v := range map[string]string{
//		pos_diff.Endpoint:                                  "https://example2.com",
//		neg_diff.Endpoint:                                  "",
//		pos_diff.Cookbooks["puppet"]["0.1.0"].LocationPath: "https://example2.com",
//		neg_diff.Cookbooks["chef"]["0.2.0"].LocationPath:   "https://example1.com",
//	} {
//		if k != v {
//			t.Fatalf("Expected: %v, got %v", v, k)
//		}
//	}
//}
//
//func Test_Diff_3_AlterAttributes(t *testing.T) {
//	data1 := new(Universe)
//	data1.Endpoint = "https://example1.com"
//	data1.Cookbooks["chef"] = Cookbook{}
//	data1.Cookbooks["chef"]["0.1.0"] = CookbookVersion{
//		"opscode",
//		"https://example1.com",
//		"https://example1.com/1",
//		map[string]string{"thing1": ">= 0.0.0"},
//	}
//	data2 := new(Universe)
//	data2.Endpoint = "https://example2.com"
//	data2.Cookbooks["chef"] = Cookbook{}
//	data2.Cookbooks["chef"]["0.1.0"] = CookbookVersion{
//		"opscode2",
//		"https://example2.com",
//		"https://example2.com/2",
//		map[string]string{"thing1": ">= 0.0.0"},
//	}
//
//	pos_diff, neg_diff, err := Diff(data1, data2)
//	if err != nil {
//		t.Fatalf("Expected no error, got %v", err)
//	}
//	for k, v := range map[int]int{
//		len(pos_diff.Cookbooks):         1,
//		len(neg_diff.Cookbooks):         0,
//		len(neg_diff.Cookbooks["chef"]): 1,
//	} {
//		if k != v {
//			t.Fatalf("Expected: %v, got: %v", v, k)
//		}
//	}
//	for k, v := range map[string]string{
//		pos_diff.Endpoint:                                "https://example2.com",
//		neg_diff.Endpoint:                                "",
//		neg_diff.Cookbooks["chef"]["0.1.0"].LocationPath: "https://example2.com",
//	} {
//		if k != v {
//			t.Fatalf("Expected: %v, got %v", v, k)
//		}
//	}
//}

func Test_New_1_NoError(t *testing.T) {
	ts := start_http()
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
	for k, v := range map[string]string{
		u.Endpoint: ts.URL + "/universe",
		u.Cookbooks["chef"]["0.12.0"].LocationType:                     test_data["chef"]["0.12.0"]["location_type"],
		u.Cookbooks["chef"]["0.12.0"].LocationPath:                     test_data["chef"]["0.12.0"]["location_path"],
		u.Cookbooks["chef"]["0.12.0"].DownloadURL:                      test_data["chef"]["0.12.0"]["download_url"],
		u.Cookbooks["chef"]["0.12.0"].Dependencies["runit"]:            ">= 0.0.0",
		u.Cookbooks["chef"]["0.12.0"].Dependencies["couchdb"]:          ">= 0.0.0",
		u.Cookbooks["chef"]["0.20.0"].LocationType:                     test_data["chef"]["0.20.0"]["location_type"],
		u.Cookbooks["chef"]["0.20.0"].LocationPath:                     test_data["chef"]["0.20.0"]["location_path"],
		u.Cookbooks["chef"]["0.20.0"].DownloadURL:                      test_data["chef"]["0.20.0"]["download_url"],
		u.Cookbooks["chef"]["0.20.0"].Dependencies["zlib"]:             ">= 0.0.0",
		u.Cookbooks["chef"]["0.20.0"].Dependencies["xml"]:              ">= 0.0.0",
		u.Cookbooks["djbdns"]["0.7.0"].LocationType:                    test_data["djbdns"]["0.7.0"]["location_type"],
		u.Cookbooks["djbdns"]["0.7.0"].LocationPath:                    test_data["djbdns"]["0.7.0"]["location_path"],
		u.Cookbooks["djbdns"]["0.7.0"].DownloadURL:                     test_data["djbdns"]["0.7.0"]["download_url"],
		u.Cookbooks["djbdns"]["0.7.0"].Dependencies["runit"]:           ">= 0.0.0",
		u.Cookbooks["djbdns"]["0.7.0"].Dependencies["build-essential"]: ">= 0.0.0",
		u.Cookbooks["djbdns"]["0.8.2"].LocationType:                    test_data["djbdns"]["0.8.2"]["location_type"],
		u.Cookbooks["djbdns"]["0.8.2"].LocationPath:                    test_data["djbdns"]["0.8.2"]["location_path"],
		u.Cookbooks["djbdns"]["0.8.2"].DownloadURL:                     test_data["djbdns"]["0.8.2"]["download_url"],
		u.Cookbooks["djbdns"]["0.8.2"].Dependencies["runit"]:           ">= 0.0.0",
		u.Cookbooks["djbdns"]["0.8.2"].Dependencies["build-essential"]: ">= 0.0.0",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
}

func Test_New_2_ConnError(t *testing.T) {
	ts := start_http()
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
		u.Cookbooks["chef-dk"]["2.0.0"].LocationType:        "opscode",
		u.Cookbooks["chef-dk"]["2.0.0"].LocationPath:        "https://supermarket.getchef.com/api/v1",
		u.Cookbooks["chef-dk"]["2.0.0"].DownloadURL:         "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
		u.Cookbooks["chef-dk"]["2.0.0"].Dependencies["dmg"]: "~> 2.2",
	} {
		if k != v {
			t.Fatalf("Expected: %v, got: %v", v, k)
		}
	}
}
