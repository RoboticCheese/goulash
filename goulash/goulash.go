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

This file implements the goulash struct, which handles the main, overarching
information about the Supermarket instace to which you're connecting
*/
package goulash

import (
	"errors"
	"github.com/RoboticCheese/goulash/common"
	"net/http"
)

// Goulash implements a data structure for a Supermarket instance.
type Goulash struct {
	common.Component
}

// New initializes and returns a new Goulash struct based on a Supermarket
// URL.
func New(url string) (g *Goulash, err error) {
	g = new(Goulash)
	g.Endpoint = url
	resp, err := http.Get(g.Endpoint + "/status")
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	return
}
