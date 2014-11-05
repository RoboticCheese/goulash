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

This file implements a struct for a cookbook version, corresponding to how one
is represented by the Supermarket API, e.g.

https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0 =>

{
	"license": "Apache v2.0",
	"tarball_file_size": 5913,
	"version": "2.0.0",
	"average_rating": null,
	"cookbook": "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk",
	"file": "https://supermarket.getchef.com/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
	"dependencies": {
		"dmg":"~> 2.2"
	}
}
*/
package cookbook_version

import (
	"encoding/json"
	"github.com/RoboticCheese/goulash/common"
	"github.com/RoboticCheese/goulash/cookbook"
	"io"
	"net/http"
)

// CookbookVersion implements a data structure for a specific version of a cookbook.
type CookbookVersion struct {
	common.Component
	License         string            `json:"license"`
	TarballFileSize int               `json:"tarball_file_size"`
	Version         string            `json:"version"`
	AverageRating   int               `json:"average_rating"` // TODO: How to distinguish nil from 0?
	Cookbook        string            `json:"cookbook"`
	File            string            `json:"file"`
	Dependencies    map[string]string `json:"dependencies"`
}

// New initializes and returns a new CookbookVersion struct based on a Cookbook.
func New(cb *cookbook.Cookbook, v string) (cv *CookbookVersion, err error) {
	cv = new(CookbookVersion)
	cv.Endpoint = cb.Endpoint + "/versions/" + v
	cv.Component, err = common.New(cv.Endpoint)
	if err != nil {
		return
	}

	resp, err := http.Get(cv.Endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = decodeJSON(resp.Body, cv)
	return
}

// Empty checks whethera CookbookVersion struct has been populated with
// anything or still holds all the base defaults.
func (cv *CookbookVersion) Empty() (empty bool) {
	empty = common.Empty(cv)
}

// Equals implements an equality test for a CookbookVersion.
func (cv1 *CookbookVersion) Equals(cv2 Supermarketer) (res bool) {
	res = common.Equals(cv1, cv2)
	return
}

// decodeJSON accepts an IO reader and a CookbookVersion struct and populates
// that struct with the JSON data.
func decodeJSON(r io.Reader, cv *CookbookVersion) (err error) {
	decoder := json.NewDecoder(r)
	return decoder.Decode(cv)
}
