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
	// Lack of "covariant return types" for interfaces makes these args
	// and return types too non-sensical.
	// Equals(Supermarketer) bool
	// Diff(Supermarketer) (Supermarketer, Supermarketer)
}

// Empty can be passed any reflect.Value and determines whether it's been
// populated with anything or still holds all the base defaults.
func Empty(s Supermarketer) (empty bool) {
	empty = emptyValue(reflect.ValueOf(s))
	return
}

// Equals does a deep comparison on two reflect.Values.
func Equals(s1 Supermarketer, s2 Supermarketer) (equal bool) {
	equal = reflect.DeepEqual(s1, s2)
	return
}

// Diff returns any attributes that have been changed from one reflect.Value
// to another.
func Diff(s1 Supermarketer, s2 Supermarketer, pos Supermarketer, neg Supermarketer) (Supermarketer, Supermarketer) {
	if Equals(s1, s2) {
		pos = nil
		neg = nil
		return pos, neg
	}

	v1 := reflect.ValueOf(s1)
	v2 := reflect.ValueOf(s2)
	vpos := reflect.ValueOf(&pos).Elem()
	vneg := reflect.ValueOf(&neg).Elem()
	p, n := diffValue(v1, v2)
	vpos.Set(p)
	vneg.Set(n)

	if Empty(pos) {
		pos = nil
	}
	if Empty(neg) {
		neg = nil
	}
	return pos, neg
}

// diffValue implements a diff check on two reflect.Values so the check can be
// iterable.
func diffValue(v1 reflect.Value, v2 reflect.Value) (vpos reflect.Value, vneg reflect.Value) {
	vpos = reflect.New(v1.Type()).Elem()
	vneg = reflect.New(v1.Type()).Elem()

	if !v1.IsValid() {
		vpos.Set(v1)
		return
	}
	if !v2.IsValid() {
		vneg.Set(v2)
		return
	}
	if v1.Type() != v2.Type() {
		vpos = v2
		vneg = v1
		return
	}

	switch v1.Kind() {
	case reflect.String:
		if v1.String() != v2.String() {
			vpos.Set(v2)
			vneg.Set(v1)
		}
	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			f1 := v1.Field(i)
			f2 := v2.Field(i)
			p, n := diffValue(f1, f2)
			vpos.Field(i).Set(p)
			vneg.Field(i).Set(n)
		}
	case reflect.Ptr:
		p, n := diffValue(v1.Elem(), v2.Elem())
		vpos.Set(p.Addr())
		vneg.Set(n.Addr())
	case reflect.Interface:
		p, n := diffValue(v1.Elem(), v2.Elem())
		vpos.Set(p)
		vneg.Set(n)
	case reflect.Map:
		for _, k := range v1.MapKeys() {
			sub_pos, sub_neg := diffValue(v1.MapIndex(k), v2.MapIndex(k))
			vpos.SetMapIndex(k, sub_pos)
			vneg.SetMapIndex(k, sub_neg)
		}
	}
	return
}

// emptyValue splits out the iterable portion of an emptiness check.
func emptyValue(v reflect.Value) (empty bool) {
	empty = true
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.String:
		if v.String() != "" {
			empty = false
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if !emptyValue(f) {
				empty = false
				break
			}
		}
	case reflect.Ptr, reflect.Interface:
		if !emptyValue(v.Elem()) {
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

func getETag(url string) (etag string, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return
	}
	etag = resp.Header.Get("etag")
	return
}
