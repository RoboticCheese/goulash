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
	"reflect"
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
func (cv *CookbookVersion) Empty() (empty bool) {
	empty = true
	for _, i := range []string{
		cv.Version,
		cv.LocationType,
		cv.LocationPath,
		cv.DownloadURL,
	} {
		if i != "" {
			empty = false
			return
		}
	}
	if len(cv.Dependencies) != 0 {
		empty = false
		return
	}
	return
}

// Equals implements an equality test for a CookbookVersion struct
func (cv1 *CookbookVersion) Equals(cv2 *CookbookVersion) (res bool) {
	res = reflect.DeepEqual(cv1, cv2)
	return
}

// Diff returns any attributes that have been changed from one CookbookVersion
// struct to another.
func (cv1 *CookbookVersion) Diff(cv2 *CookbookVersion) (pos, neg *CookbookVersion) {
	if cv1.Equals(cv2) {
		return
	}
	pos = cv1.positiveDiff(cv2)
	neg = cv1.negativeDiff(cv2)
	return
}

// positiveDiff returns any attributes that have been added or changed from one
// CookbookVersion struct to another.
func (cv1 *CookbookVersion) positiveDiff(cv2 *CookbookVersion) (pos *CookbookVersion) {
	if cv1.Equals(cv2) {
		return
	}
	pos = NewCookbookVersion()

	r1 := reflect.ValueOf(cv1).Elem()
	r2 := reflect.ValueOf(cv2).Elem()
	rres := reflect.ValueOf(pos).Elem()
	for i := 0; i < r1.NumField(); i++ {
		f1 := r1.Field(i)
		f2 := r2.Field(i)

		switch f1.Kind() {
		case reflect.String:
			if f1.String() != f2.String() && f2.String() != "" {
				rres.Field(i).Set(f2)
			}
		case reflect.Map:
			for _, k := range f2.MapKeys() {
				if f1.MapIndex(k).Kind() == reflect.Invalid || f1.MapIndex(k).String() != f2.MapIndex(k).String() {
					rres.Field(i).SetMapIndex(k, f2.MapIndex(k))
				}
			}
		}
	}
	if pos.Empty() {
		pos = nil
	}
	return
}

// negativeDiff returns any attributes that have been deleted from one
// CookbookVersion struct to another.
func (cv1 *CookbookVersion) negativeDiff(cv2 *CookbookVersion) (neg *CookbookVersion) {
	if cv1.Equals(cv2) {
		return
	}
	neg = NewCookbookVersion()

	r1 := reflect.ValueOf(cv1).Elem()
	r2 := reflect.ValueOf(cv2).Elem()
	rres := reflect.ValueOf(neg).Elem()
	for i := 0; i < r1.NumField(); i++ {
		f1 := r1.Field(i)
		f2 := r2.Field(i)

		switch f1.Kind() {
		case reflect.String:
			if f2.String() == "" && f1.String() != "" {
				rres.Field(i).Set(f1)
			}
		case reflect.Map:
			for _, k := range f1.MapKeys() {
				if f2.MapIndex(k).Kind() == reflect.Invalid || f2.MapIndex(k).String() == "" {
					rres.Field(i).SetMapIndex(k, f1.MapIndex(k))
				}
			}
		}
	}
	if neg.Empty() {
		neg = nil
	}
	return
}
