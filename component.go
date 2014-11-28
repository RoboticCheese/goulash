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

This file defines a shared Component struct to be used as an anonymous field
by any other Goulash structs.
*/
package goulash

import (
	"net/http"

	"github.com/RoboticCheese/goulash/common"
)

// Component defines variables to be shared by all the Goulash structs.
type Component struct {
	Endpoint string
	ETag     string
}

// NewComponent creates a new Component struct from a given endpoint string and
// returns that struct and any error.
func NewComponent(endpoint string) (c Component, err error) {
	c = InitComponent()
	c.Endpoint = endpoint
	err = c.getETag()
	return
}

// InitComponent generates an empty Component struct.
func InitComponent() (c Component) {
	c = Component{}
	return
}

// Empty checks whether a Component struct has been populated with anything
// or still holds all the base defaults.
func (c *Component) Empty() (empty bool) {
	empty = common.Empty(c)
	return
}

// Equals checks whether one Component struct is equal to another.
func (c *Component) Equals(c2 *Component) (equal bool) {
	equal = common.Equals(c, c2)
	return
}

// Diff returns any attributes that have been changed from one Component struct
// to another.
func (c *Component) Diff(c2 *Component) (pos, neg *Component) {
	ipos, ineg := common.Diff(c, c2, &Component{}, &Component{})
	if ipos != nil {
		pos = ipos.(*Component)
	} else {
		pos = nil
	}
	if ineg != nil {
		neg = ineg.(*Component)
	} else {
		neg = nil
	}
	return
}

// getETag accepts a URL string and returns any ETag header and error returned
// from an HTTP HEAD on that URL.
func (c *Component) getETag() (err error) {
	resp, err := http.Head(c.Endpoint)
	if err != nil {
		return
	}
	c.ETag = resp.Header.Get("etag")
	return
}
