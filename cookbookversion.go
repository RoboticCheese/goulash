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
Package goulash implements a Go client library for the Chef Supermarket API.

This file defines a CookbookVersion struct, corresponding to how a cookbook
version is represented by the API, e.g.

https://supermarket.chef.io/api/v1/cookbooks/chef-dk/versions/2.0.0 =>

{
	"license": "Apache v2.0",
	"tarball_file_size": 5913,
	"version": "2.0.0",
	"average_rating": null,
	"cookbook": "https://supermarket.chef.io/api/v1/cookbooks/chef-dk",
	"file": "https://supermarket.chef.io/api/v1/cookbooks/chef-dk/versions/2.0.0/download",
	"dependencies": {
		"dmg":"~> 2.2"
	}
}
*/
package goulash

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/RoboticCheese/goulash/common"
)

// CookbookVersion implements a data structure for a specific version of a cookbook.
type CookbookVersion struct {
	Component
	License         string            `json:"license"`
	TarballFileSize int               `json:"tarball_file_size"`
	Version         string            `json:"version"`
	AverageRating   int               `json:"average_rating"` // TODO: How to distinguish nil from 0?
	Cookbook        string            `json:"cookbook"`
	File            string            `json:"file"`
	Dependencies    map[string]string `json:"dependencies"`
}

// NewCookbookVersion initializes and returns a new CookbookVersion struct
// based on a Cookbook.
func NewCookbookVersion(cb *Cookbook, v string) (cv *CookbookVersion, err error) {
	cv = InitCookbookVersion()
	cv.Endpoint = cb.Endpoint + "/versions/" + v
	cv.Component, err = NewComponent(cv.Endpoint)
	if err != nil {
		return
	}

	resp, err := http.Get(cv.Endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = cv.decodeJSON(resp.Body)
	return
}

// InitCookbookVersion generates an empty CookbookVersion struct.
func InitCookbookVersion() (cv *CookbookVersion) {
	cv = new(CookbookVersion)
	cv.Dependencies = map[string]string{}
	return
}

// Empty checks whether a CookbookVersion struct has been populated with
// anything or still holds all the base defaults.
func (cv *CookbookVersion) Empty() (empty bool) {
	empty = common.Empty(cv)
	return
}

// Equals implements an equality test for a CookbookVersion.
func (cv *CookbookVersion) Equals(cv2 common.Supermarketer) (res bool) {
	res = common.Equals(cv, cv2)
	return
}

// Diff returns any attributes added/changed/removed from one CookbookVersion
// struct to another, represented by a positive and negative diff
// CookbookVersion.
func (cv *CookbookVersion) Diff(cv2 *CookbookVersion) (pos, neg *CookbookVersion) {
	ipos, ineg := common.Diff(cv, cv2, &CookbookVersion{}, &CookbookVersion{})
	if ipos != nil {
		pos = ipos.(*CookbookVersion)
	} else {
		pos = nil
	}
	if ineg != nil {
		neg = ineg.(*CookbookVersion)
	} else {
		neg = nil
	}
	return
}

// decodeJSON accepts an IO reader and a CookbookVersion struct and populates
// that struct with the JSON data.
func (cv *CookbookVersion) decodeJSON(r io.Reader) (err error) {
	decoder := json.NewDecoder(r)
	return decoder.Decode(cv)
}
