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

import (
	"github.com/RoboticCheese/goulash/common"
)

// CookbookVersion implements a struct for each cookbook version underneath a
// Universe.
type CookbookVersion struct {
	Version      string
	LocationType string            `json:"location_type"`
	LocationPath string            `json:"location_path"`
	DownloadURL  string            `json:"download_url"`
	Dependencies map[string]string `json:"dependencies"`
}

// NewCookbookVersion generates an empty CookbookVersion struct.
func NewCookbookVersion() (cv *CookbookVersion) {
	cv = new(CookbookVersion)
	cv.Dependencies = map[string]string{}
	return
}

// Empty checks whether a CookbookVersion struct has been populated with
// anything or still holds all the base defaults.
func (cv CookbookVersion) Empty() (empty bool) {
	empty = common.Empty(cv)
	return
}

// Equals implements an equality test for a CookbookVersion struct
func (cv1 CookbookVersion) Equals(cv2 *CookbookVersion) (res bool) {
	res = common.Equals(cv1, cv2)
	return
}

// Diff returns any attributes that have been changed from one CookbookVersion
// struct to another.
func (cv1 CookbookVersion) Diff(cv2 *CookbookVersion) (pos, neg *CookbookVersion) {
	ipos, ineg := common.Diff(cv1, *cv2, CookbookVersion{}, CookbookVersion{})
	if ipos != nil {
		cpos := ipos.(CookbookVersion)
		pos = &cpos
	} else {
		pos = nil
	}
	if ineg != nil {
		cneg := ineg.(CookbookVersion)
		neg = &cneg
	} else {
		neg = nil
	}
	return
}
