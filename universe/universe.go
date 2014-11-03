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

// Universe contains a Cookbooks map of cookbook name strings to Cookbook items.
type Universe struct {
	common.Component
	APIInstance *api_instance.APIInstance
	Cookbooks   map[string]*Cookbook
}

// New accepts a pointer to an APIInstance struct and uses it to initialize
// and return a pointer to a new Universe struct.
func New(i *api_instance.APIInstance) (u *Universe, err error) {
	u = NewUniverse()
	u.APIInstance = i
	u.Component, err = common.New(u.APIInstance.BaseURL + "/universe")
	if err != nil {
		return
	}

	resp, err := http.Get(u.Endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create a temporary map that corresponds more closely to what the
	// universe JSON data looks like
	temp_u := map[string]map[string]*CookbookVersion{}

	err = decodeJSON(resp.Body, &temp_u)
	if err != nil {
		return
	}
	// Fill in the Universe struct with the JSON data gathered above
	for cb_name, cb := range temp_u {
		u.Cookbooks[cb_name] = NewCookbook()
		u.Cookbooks[cb_name].Name = cb_name
		for cv_name, cv := range cb {
			cv.Version = cv_name
			u.Cookbooks[cb_name].Versions[cv_name] = cv
		}
	}
	return
}

// NewUniverse generates an empty Universe struct.
func NewUniverse() (u *Universe) {
	u = new(Universe)
	u.Cookbooks = map[string]*Cookbook{}
	return
}

// Empty checks whether a Universe struct has been populated with anything or
// still holds all the base defaults.
func (u *Universe) Empty() (empty bool) {
	empty = common.Empty(u)
	return
}

// Equals implements an equality test for a Universe.
func (u1 *Universe) Equals(u2 common.Supermarketer) (res bool) {
	res = common.Equals(u1, u2)
	return
}

// Update refreshes a Universe struct and returns the diff of the original
// Universe and the updated one.
func (u *Universe) Update() (pos_diff, neg_diff *Universe, err error) {
	// Try to use the HTTP ETag header first; don't download the entire
	// universe JSON if we don't need to.
	if u.ETag != "" {
		// Fall through to the regular compare if there's an error
		tmp, _ := common.New(u.Endpoint)
		if tmp.ETag != "" && tmp.ETag == u.ETag {
			return
		}
	}

	cur_u, err := New(u.APIInstance)
	if err != nil {
		return
	}
	pos_diff, neg_diff = u.Diff(cur_u)
	*u = *cur_u
	return
}

// Diff returns any attributes that have changed from one Universe struct to
// another.
func (u1 *Universe) Diff(u2 *Universe) (pos, neg *Universe) {
	if u1.Equals(u2) {
		return
	}
	pos = NewUniverse()
	neg = NewUniverse()

	for k, _ := range u1.Cookbooks {
		if u2.Cookbooks[k] == nil {
			neg.Cookbooks[k] = u1.Cookbooks[k]
		} else if !u1.Cookbooks[k].Equals(u2.Cookbooks[k]) {
			pos.Cookbooks[k], neg.Cookbooks[k] = u1.Cookbooks[k].Diff(u2.Cookbooks[k])
		}
	}
	for k, _ := range u2.Cookbooks {
		if u1.Cookbooks[k] == nil {
			pos.Cookbooks[k] = u2.Cookbooks[k]
		}
	}
	if pos.Empty() {
		pos = nil
	}
	if neg.Empty() {
		neg = nil
	}
	return
}

// decodeJSON accepts an IO reader and a Universe struct and populates that
// struct with the JSON data, after doing some extra parsing to account for the
// variant cookbook name and version number keys.
func decodeJSON(r io.Reader, u *map[string]map[string]*CookbookVersion) (err error) {
	decoder := json.NewDecoder(r)
	return decoder.Decode(u)
}
