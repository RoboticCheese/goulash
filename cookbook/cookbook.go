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
Package cookbook implements individual cookbook functionality.

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
	"io"
	"net/http"
	"reflect"

	"github.com/RoboticCheese/goulash/apiinstance"
	"github.com/RoboticCheese/goulash/common"
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
func New(i *apiinstance.APIInstance, name string) (c *Cookbook, err error) {
	c = NewCookbook()
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

// NewCookbook generates an empty Cookbook struct.
func NewCookbook() (c *Cookbook) {
	c = new(Cookbook)
	c.Versions = []string{}
	c.Metrics = Metrics{Downloads: Downloads{}}
	return
}

// Empty checks whether a Cookbook struct has been populated with anything or
// still holds all the base defaults.
func (c *Cookbook) Empty() (empty bool) {
	empty = common.Empty(c)
	return
}

// Equals implements an equality test for a Cookbook.
func (c *Cookbook) Equals(c2 common.Supermarketer) (res bool) {
	res = common.Equals(c, c2)
	return
}

// Diff returns any attributes added/changed/removed from one Cookbook struct
// to another
func (c *Cookbook) Diff(c2 *Cookbook) (pos, neg *Cookbook) {
	if c.Equals(c2) {
		return
	}
	r1 := reflect.ValueOf(c).Elem()
	r2 := reflect.ValueOf(c2).Elem()

	if !r1.IsValid() {
		pos = c2
		return
	}
	if !r2.IsValid() {
		neg = c
		return
	}

	pos = NewCookbook()
	neg = NewCookbook()
	rpos := reflect.ValueOf(pos).Elem()
	rneg := reflect.ValueOf(neg).Elem()
	for i := 0; i < r1.NumField(); i++ {
		f1 := r1.Field(i)
		f2 := r2.Field(i)

		switch f1.Kind() {
		case reflect.String:
			if f1.String() != f2.String() {
				rpos.Field(i).Set(f2)
				rneg.Field(i).Set(f1)
			}
		case reflect.Int:
			if f1.Int() != f2.Int() {
				rpos.Field(i).Set(f2)
				rneg.Field(i).Set(f1)
			}
		case reflect.Bool:
			if f1.Bool() != f2.Bool() {
				rpos.Field(i).Set(f2)
				rneg.Field(i).Set(f1)
			}
		case reflect.Slice:
			for j := 0; j < f1.Len(); j++ {
				found := false
				for k := 0; k < f2.Len(); k++ {
					if f2.Index(k) == f1.Index(j) {
						found = true
						break
					}
				}
				if found == false {
					rneg.Field(i).Set(reflect.Append(rneg.Field(i), f1.Index(j)))
				}
			}
			for j := 0; j < f2.Len(); j++ {
				found := false
				for k := 0; k < f1.Len(); k++ {
					if f1.Index(k) == f2.Index(j) {
						found = true
						break
					}
				}
				if found == false {
					rpos.Field(i).Set(reflect.Append(rpos.Field(i), f2.Index(j)))
				}
			}
		case reflect.Struct:
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

// decodeJSON accepts an IO reader and a Cookbook struct and populates that
// struct with the JSON data.
func decodeJSON(r io.Reader, c *Cookbook) (err error) {
	decoder := json.NewDecoder(r)
	return decoder.Decode(c)
}
