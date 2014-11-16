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

This file implements a struct for a cookbook as described by a Berkshelf-style
universe endpoint, e.g.

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
	...
*/
package universe

import (
	"github.com/RoboticCheese/goulash/common"
)

// Cookbook is just a map of version strings to Version structs
type Cookbook struct {
	Name     string
	Versions map[string]*CookbookVersion
}

// NewCookbook generates an empty Cookbook struct.
func NewCookbook() (c *Cookbook) {
	c = new(Cookbook)
	c.Versions = map[string]*CookbookVersion{}
	return
}

// Empty checks whether a Cookbook struct has been populated with anything or
// still holds all the base defaults.
func (c Cookbook) Empty() (empty bool) {
	empty = common.Empty(c)
	return
}

// Equals implements an equality test for a Cookbook.
func (c1 Cookbook) Equals(c2 *Cookbook) (res bool) {
	res = common.Equals(c1, c2)
	return
}

// Diff returns any attributes that have changed from one Cookbook struct to
// another.
func (c1 Cookbook) Diff(c2 *Cookbook) (pos, neg *Cookbook) {
	ipos, ineg := common.Diff(c1, *c2, Cookbook{}, Cookbook{})
	if ipos != nil {
		cpos := ipos.(Cookbook)
		pos = &cpos
	} else {
		pos = nil
	}
	if ineg != nil {
		cneg := ineg.(Cookbook)
		neg = &cneg
	} else {
		neg = nil
	}
	return
}
