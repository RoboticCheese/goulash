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
	//	Equals(*Supermarketer) bool
	//	Diff(*Supermarketer) (*Supermarketer, *Supermarketer)
}

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

func getETag(url string) (etag string, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return
	}
	etag = resp.Header.Get("etag")
	return
}
