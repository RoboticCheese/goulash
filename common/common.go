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
	empty = true
	v := reflect.ValueOf(s)
	if !v.IsValid() {
		return
	}
	r := v
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		r = v.Elem()
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

// Equals does a deep comparison on two reflect.Values.
func Equals(s1 Supermarketer, s2 Supermarketer) (equal bool) {
	equal = reflect.DeepEqual(s1, s2)
	return
}

// Diff returns any attributes that have been changed from one reflect.Value
// to another.
func Diff(s1 Supermarketer, s2 Supermarketer, pos Supermarketer, neg Supermarketer) (Supermarketer, Supermarketer) {
	// TODO: Return an error if v1 and v2 are different types?
	v1 := reflect.ValueOf(s1)
	if v1.Kind() == reflect.Interface || v1.Kind() == reflect.Ptr {
		v1 = reflect.ValueOf(s1).Elem()
	}
	v2 := reflect.ValueOf(s2)
	if v2.Kind() == reflect.Interface || v2.Kind() == reflect.Ptr {
		v2 = reflect.ValueOf(s2).Elem()
	}
	vpos := reflect.ValueOf(pos)
	if vpos.Kind() == reflect.Interface || vpos.Kind() == reflect.Ptr {
		vpos = reflect.ValueOf(pos).Elem()
	}
	vneg := reflect.ValueOf(neg)
	if vneg.Kind() == reflect.Interface || vneg.Kind() == reflect.Ptr {
		vneg = reflect.ValueOf(neg).Elem()
	}
	if v1.Type() != v2.Type() {
		pos = s2
		neg = s1
		return pos, neg
	}
	if Equals(s1, s2) {
		pos = nil
		neg = nil
		return pos, neg
	}
	if !v1.IsValid() {
		pos = s2
		neg = nil
		return pos, neg
	}
	if !v2.IsValid() {
		pos = nil
		neg = s1
		return pos, neg
	}

	for i := 0; i < v1.NumField(); i++ {
		f1 := v1.Field(i)
		f2 := v2.Field(i)

		pos_diff, neg_diff := diffValue(f1, f2)
		if pos_diff.IsValid() {
			vpos.Field(i).Set(pos_diff)
		}
		if neg_diff.IsValid() {
			vneg.Field(i).Set(neg_diff)
		}
	}
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
	if !v1.IsValid() {
		vpos = v1
		return
	}
	if !v2.IsValid() {
		vneg = v2
		return
	}
	vpos = reflect.New(v1.Type()).Elem()
	vneg = reflect.New(v1.Type()).Elem()
	switch v1.Kind() {
	case reflect.String:
		if v1.String() != v2.String() {
			vpos.Set(v2)
			vneg.Set(v1)
		}
	case reflect.Struct:
		method := v1.Addr().MethodByName("Equals")
		arg := []reflect.Value{v2}
		if !method.Call(arg)[0].Bool() {
			method := v1.Addr().MethodByName("Diff")
			arg := []reflect.Value{v2}
			diffs := method.Call(arg)
			vpos.Set(diffs[0])
			vneg.Set(diffs[1])
		}
	case reflect.Ptr:
		p, n := diffValue(reflect.Indirect(v1), reflect.Indirect(v2))
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
		if !emptyValue(reflect.Indirect(v)) {
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
