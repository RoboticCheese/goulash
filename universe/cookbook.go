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
	"reflect"
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
func (c *Cookbook) Empty() (empty bool) {
	empty = common.Empty(c)
	return
}

// Equals implements an equality test for a Cookbook.
func (c1 *Cookbook) Equals(c2 common.Supermarketer) (res bool) {
	res = common.Equals(c1, c2)
	return
}

// Diff returns any attributes that have changed from one Cookbook struct to
// another.
func (c1 *Cookbook) Diff(c2 *Cookbook) (pos, neg *Cookbook) {
	if c1.Equals(c2) {
		return
	}
	r1 := reflect.ValueOf(c1).Elem()
	r2 := reflect.ValueOf(c2).Elem()

	if !r1.IsValid() {
		pos = c2
		return
	}
	if !r2.IsValid() {
		neg = c1
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
		case reflect.Map:
			for _, k := range f1.MapKeys() {
				if f2.MapIndex(k).Kind() == reflect.Invalid {
					rneg.Field(i).SetMapIndex(k, f1.MapIndex(k))
				} else {
					meth := f1.MapIndex(k).MethodByName("Equals")
					arg := []reflect.Value{f2.MapIndex(k)}
					if !meth.Call(arg)[0].Bool() {
						meth := f1.MapIndex(k).MethodByName("Diff")
						arg := []reflect.Value{f2.MapIndex(k)}
						diffs := meth.Call(arg)
						rpos.Field(i).SetMapIndex(k, diffs[0])
						rneg.Field(i).SetMapIndex(k, diffs[1])
					}
				}
			}
			for _, k := range f2.MapKeys() {
				if f1.MapIndex(k).Kind() == reflect.Invalid {
					rpos.Field(i).SetMapIndex(k, f2.MapIndex(k))
				}
			}
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
