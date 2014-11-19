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
Package common implements a shared set of Goulash functionality.

This file defines a shared Component struct for any API endpoint.
*/
package common

// Component defines variables to be shared by all the Goulash structs.
type Component struct {
	Endpoint string
	ETag     string
}

// New creates a new Component struct from a given endpoint string and returns
// that struct and any error.
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
func (c *Component) Equals(c2 *Component) (equal bool) {
	equal = Equals(c, c2)
	return
}

// Diff returns any attributes that have been changed from one Component struct
// to another.
func (c *Component) Diff(c2 *Component) (pos, neg *Component) {
	// TODO: How do we handle when there's a struct/pointer type mismatch
	// between c1 and c2?
	ipos, ineg := Diff(c, c2, &Component{}, &Component{})
	if ipos != nil {
		cpos := ipos.(*Component)
		pos = cpos
	} else {
		pos = nil
	}
	if ineg != nil {
		cneg := ineg.(*Component)
		neg = cneg
	} else {
		neg = nil
	}
	return
}
