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

// Downloads represents the Downloads section of the metrics data.
type Downloads struct {
	Total    int            `json:"total"`
	Versions map[string]int `json:"versions"`
}

// Metrics represents the metrics section of the cookbook data.
type Metrics struct {
	Downloads Downloads `json:"downloads"`
	Followers int
}

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
	Metrics           Metrics  `json:"metrics"`
}

// New initializes and returns a new Cookbook struct based on a Supermarket
// struct and cookbook name.
func New(i *api_instance.APIInstance, name string) (c *Cookbook, err error) {
	c = new(Cookbook)
	c.Endpoint = i.Endpoint + "/cookbooks/" + name
	c.Component, err = common.New(c.Endpoint)
	if err != nil {
		return
	}

	resp, err := http.Get(c.Endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = decodeJSON(resp.Body, c)
	return
}

// Equals implements an equality test for a Cookbook.
func (c1 Cookbook) Equals(c2 Cookbook) (res bool, err error) {
	res = false
	for _, i := range [][]string{
		{c1.Endpoint, c2.Endpoint},
		{c1.Name, c2.Name},
		{c1.Maintainer, c2.Maintainer},
		{c1.Description, c2.Description},
		{c1.Category, c2.Category},
		{c1.LatestVersion, c2.LatestVersion},
		{c1.ExternalURL, c2.ExternalURL},
		{c1.CreatedAt, c2.CreatedAt},
		{c1.UpdatedAt, c2.UpdatedAt},
	} {
		if i[0] != i[1] {
			return
		}
	}
	for _, i := range [][]int{
		{c1.AverageRating, c2.AverageRating},
		{len(c1.Versions), len(c2.Versions)},
		{c1.Metrics.Downloads.Total, c2.Metrics.Downloads.Total},
		{c1.Metrics.Followers, c2.Metrics.Followers},
	} {
		if i[0] != i[1] {
			return
		}
	}
	if c1.Deprecated != c2.Deprecated || c1.FoodcriticFailure != c2.FoodcriticFailure {
		return
	}
	for k, v := range c1.Versions {
		if v != c2.Versions[k] {
			return
		}
	}
	for k, v := range c1.Metrics.Downloads.Versions {
		if v != c2.Metrics.Downloads.Versions[k] {
			return
		}
	}
	res = true
	return
}

// Diff returns any attributes added/changed/removed from one Cookbook struct
// to another
func (c1 *Cookbook) Diff(c2 *Cookbook) (pos, neg *Cookbook) {
	if c1.Equals(c2) {
		return
	}
	pos = c1.positiveDiff(c2)
	neg = c1.negativeDiff(c2)
	return
}

// positiveDiff returns any attributes that have been added or modified (a
// positive diff) from one Cookbook struct to another.
func (c1 *Cookbook) positiveDiff(c2 *Cookbook) (diff *Cookbook) {
}

// negativeDiff returns any attributes that have been removed (a negative diff)
// from one Cookbook struct to another.
func (c1 *Cookbook) negativeDiff(c2 *Cookbook) (diff *Cookbook) {
}

// decodeJSON accepts an IO reader and a Cookbook struct and populates that
// struct with the JSON data.
func decodeJSON(r io.Reader, c *Cookbook) (err error) {
	decoder := json.NewDecoder(r)
	return decoder.Decode(c)
}
