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

This file implements a struct for the Berkshelf-style universe endpoint, e.g.

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
		"0.20.0": {
			"location_type": "opscode",
			"location_path": "https://supermarket.getchef.com/api/v1",
			"download_url": "https://supermarket.getchef.com/api/v1/cookbooks/chef/versions/0.20.0/download",
			"dependencies": {
				"zlib":">= 0.0.0",
				"xml": ">= 0.0.0",
				...
			}
		},
		...
	},
	"djbdns": {
		"0.7.0": {
			"location_type": "opscode",
			"location_path": "https://supermarket.getchef.com/api/v1",
			"download_url": "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.7.0/download",
			"dependencies": {
				"runit": ">= 0.0.0",
				"build-essential": ">= 0.0.0",
				...
			}
		},
		"0.8.2": {
			"location_type": "opscode",
			"location_path": "https://supermarket.getchef.com/api/v1",
			"download_url": "https://supermarket.getchef.com/api/v1/cookbooks/djbdns/versions/0.8.2/download",
			"dependencies": {
				"runit": ">= 0.0.0",
				"build-essential": ">= 0.0.0",
				...
			}
		},
		...
	},
	...
*/
package universe

import (
	"encoding/json"
	"github.com/RoboticCheese/goulash/api_instance"
	"github.com/RoboticCheese/goulash/common"
	"io"
	"net/http"
)

// Cookbook is just a map of version strings to Version structs
type Cookbook map[string]CookbookVersion

// Universe contains a Cookbooks map of cookbook name strings to Cookbook items.
type Universe struct {
	common.Component
	Cookbooks map[string]Cookbook
}

// New initializes and returns a new Universe struct based on an API instance.
func New(i *api_instance.APIInstance) (u *Universe, err error) {
	u = new(Universe)
	u.Endpoint = i.BaseURL + "/universe"

	resp, err := http.Get(u.Endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = decodeJSON(resp.Body, u)
	return
}

// decodeJSON accepts an IO reader and a Universe struct and populates that
// struct with the JSON data, after doing some extra parsing to account for the
// variant cookbook name and version number keys.
func decodeJSON(r io.Reader, u *Universe) (err error) {
	decoder := json.NewDecoder(r)
	return decoder.Decode(&u.Cookbooks)
}
