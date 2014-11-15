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

This file defines a common, shared Component struct.
*/
package common

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
func (c Component) Empty() (empty bool) {
	empty = Empty(c)
	return
}

// Equals checks whether one Component struct is equal to another.
func (c1 Component) Equals(c2 *Component) (equal bool) {
	equal = Equals(c1, c2)
	return
}

// Diff returns any attributes that have been changed from one Component struct
// to another.
func (c1 Component) Diff(c2 *Component) (pos, neg *Component) {
	// TODO: How do we handle when there's a struct/pointer type mismatch
	// between c1 and c2?
	ipos, ineg := Diff(c1, *c2, Component{}, Component{})
	if ipos != nil {
		cpos := ipos.(Component)
		pos = &cpos
	} else {
		pos = nil
	}
	if ineg != nil {
		cneg := ineg.(Component)
		neg = &cneg
	} else {
		neg = nil
	}
	return
}
