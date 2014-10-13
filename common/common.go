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
)

// Supermarketer implements an interface shared by all the Goulash structs.
// type Supermarketer interface{}

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

func getETag(url string) (etag string, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return
	}
	etag = resp.Header.Get("etag")
	return
}
