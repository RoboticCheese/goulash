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

This file defines a set of common interfaces and structs.
*/
package common

import (
	"net/http"
	"reflect"
)

// Supermarketer implements an interface shared by all the Goulash structs.
type Supermarketer interface {
	Empty() bool
	Equals(Supermarketer) bool
	Diff(Supermarketer) (Supermarketer, Supermarketer)
}

// Empty can be passed any implementer of the Supermarketer interface and
// determines whether it's been populated with anything or still holds all the
// base defaults.
func Empty(s Supermarketer) (empty bool) {
	empty = true
	r := reflect.ValueOf(s).Elem()
	if s == nil || !r.IsValid() {
		return
	}
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		if !emptyValue(f) {
			empty = false
			break
		}
	}
	return
}

// Equals does a deep comparison on two implementers of the Supermarketer
// interface to determine whether they're equal or not.
func Equals(s1 Supermarketer, s2 Supermarketer) (equal bool) {
	equal = reflect.DeepEqual(s1, s2)
	return
}

// Diff returns any attributes that have been changed from one implementer of
// the Supermarketer interface to another by filling in two empty structs for
// a positive and negative diff.
func Diff(s1 Supermarketer, s2 Supermarketer, pos Supermarketer, neg Supermarketer) (Supermarketer, Supermarketer) {
	if s1.Equals(s2) {
		pos = nil
		neg = nil
		return pos, neg
	}
	r1 := reflect.Indirect(reflect.ValueOf(s1))
	r2 := reflect.Indirect(reflect.ValueOf(s2))
	if !r1.IsValid() {
		pos = s2
		neg = nil
		return pos, neg
	}
	if !r2.IsValid() {
		pos = nil
		neg = s1
		return pos, neg
	}

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
					// TODO: Needs to support more than map[string]string k:v
				} else if f1.MapIndex(k).String() != f2.MapIndex(k).String() {
					rpos.Field(i).SetMapIndex(k, f2.MapIndex(k))
					rneg.Field(i).SetMapIndex(k, f1.MapIndex(k))
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
	return pos, neg
}

// emptyValue implements an emptiness check for a reflect.Value.
func emptyValue(v reflect.Value) (empty bool) {
	empty = true
	switch v.Kind() {
	case reflect.String:
		if v.String() != "" {
			empty = false
		}
	case reflect.Struct:
		method := v.Addr().MethodByName("Empty")
		if !method.Call([]reflect.Value{})[0].Bool() {
			empty = false
		}
	case reflect.Ptr:
		method := v.MethodByName("Empty")
		if !method.Call([]reflect.Value{})[0].Bool() {
			empty = false
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if !emptyValue(v.MapIndex(k)) {
				empty = false
				break
			}
		}
	}
	return
}

// Component defines variables to be shared by all the Goulash structs.
type Component struct {
	Endpoint string
	ETag     string
}

func New(endpoint string) (c Component, err error) {
	c = Component{}
	c.Endpoint = endpoint
	c.ETag, err = getETag(c.Endpoint)
	return
}

// Empty checks whether a Component struct has been populated with anything
// or still holds all the base defaults.
func (c *Component) Empty() (empty bool) {
	empty = Empty(c)
	return
}

// Equals checks whether one Component struct is equal to another.
func (c1 *Component) Equals(c2 Supermarketer) (equal bool) {
	equal = Equals(c1, c2)
	return
}

// Diff returns any attributes that have been changed from one Component struct
// to another.
func (c1 *Component) Diff(c2 Supermarketer) (pos, neg Supermarketer) {
	pos, neg = Diff(c1, c2, &Component{}, &Component{})
	return
}

func getETag(url string) (etag string, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return
	}
	etag = resp.Header.Get("etag")
	return
}
