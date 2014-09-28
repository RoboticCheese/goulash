// Author:: Jonathan Hartman (<j@p4nt5.com>)
//
// Copyright (C) 2014, Jonathan Hartman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package goulash implements an API client for the Chef Supermarket.

This file implements a struct for a cookbook version as described by a
Berkshelf-style universe endpoint, e.g.

https://supermarket.getchef.com/universe =>

{
	"chef": {
		"0.12.0": {
			"location_type": "opscode",
			"location_path": "https://supermarket.getchef.com/api/v1",
			"download_url": "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.12.0/download",
			"dependencies": {
				"runit":">= 0.0.0",
				"couchdb":">= 0.0.0",
				...
			}
		},
		...
	},
	...
*/
package universe

// CookbookVersion implements a struct for each cookbook version underneath a
// Universe.
type CookbookVersion struct {
	LocationType string            `json:"location_type"`
	LocationPath string            `json:"location_path"`
	DownloadURL  string            `json:"download_url"`
	Dependencies map[string]string `json:"dependencies"`
}

// Equas implements an equality test for a CookbookVersion struct
func (cv1 CookbookVersion) Equals(cv2 CookbookVersion) (res bool, err error) {
	res = false
	for _, i := range [][]string{
		{cv1.LocationType, cv2.LocationType},
		{cv1.LocationPath, cv2.LocationPath},
		{cv1.DownloadURL, cv2.DownloadURL},
	} {
		if i[0] != i[1] {
			return
		}
	}
	if len(cv1.Dependencies) != len(cv2.Dependencies) {
		return
	}
	for k, v := range cv1.Dependencies {
		if v != cv2.Dependencies[k] {
			return
		}
	}
	res = true
	return
}
