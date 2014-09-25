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

This file implements a struct for a cookbook, corresponding to how one is
represented by the Supermarket API, e.g.

https://supermarket.getchef.com/api/v1/cookbooks/chef-dk =>

{
	"name": "chef-dk",
	"maintainer": "roboticcheese",
	"description": "Installs/configures the Chef-DK",
	"category": "Other",
	"latest_version": "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1",
	"external_url": "https://github.com/RoboticCheese/chef-dk-chef",
	"average_rating": null,
	"created_at": "2014-06-24T01:14:49.000Z",
	"updated_at": "2014-09-20T04:46:00.780Z",
	"deprecated": false,
	"foodcritic_failure": false,
	"versions": [
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.1",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/1.0.2",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/1.0.0",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/0.3.2",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/0.3.0",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/0.2.0",
		"https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/0.1.0"
	],
	"metrics": {
		"downloads": {
			"total":2150582,
			"versions": {
				"0.1.0": 376076,
				"0.2.0": 376073,
				"0.3.0": 376101,
				"0.3.2": 376236,
				"1.0.0": 333166,
				"1.0.2": 265139,
				"2.0.0": 32520,
				"2.0.1": 15271
			}
		},
		"followers": 7
	}
}
*/
package cookbook

import (
	"encoding/json"
	"github.com/RoboticCheese/goulash/api_instance"
	"github.com/RoboticCheese/goulash/common"
	"io"
	"net/http"
)

// Cookbook implements a data structure for a single Chef cookbook.
type Cookbook struct {
	common.Component
	Name              string   `json:"name"`
	Maintainer        string   `json:"maintainer"`
	Description       string   `json:"description"`
	Category          string   `json:"category"`
	LatestVersion     string   `json:"latest_version"`
	ExternalURL       string   `json:"external_url"`
	AverageRating     int      `json:"average_rating"` // TODO: How to distinguish nil and 0?
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	Deprecated        bool     `json:"deprecated"`
	FoodcriticFailure bool     `json:"foodcritic_failure"` // TODO: How to distinguish nil and false?
	Versions          []string `json:"versions"`
	Metrics           struct {
		Downloads struct {
			Total    int            `json:"total"`
			Versions map[string]int `json:"versions"`
		} `json:"downloads"`
		Followers int `json:"followers"`
	} `json:"metrics"`
}

// New initializes and returns a new Cookbook struct based on a Supermarket
// struct and cookbook name.
func New(i *api_instance.APIInstance, name string) (c *Cookbook, err error) {
	c = new(Cookbook)
	c.Endpoint = i.Endpoint + "/cookbooks/" + name

	resp, err := http.Get(c.Endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = decodeJSON(resp.Body, c)
	return
}

// decodeJSON accepts an IO reader and a Cookbook struct and populates that
// struct with the JSON data.
func decodeJSON(r io.Reader, c *Cookbook) (err error) {
	decoder := json.NewDecoder(r)
	return decoder.Decode(c)
}
