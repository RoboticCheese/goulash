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

This file defines an APIInstance struct, representing an API conection to be
consumed by other Goulash components.
*/
package goulash

import (
	"errors"
	"net/http"

	"github.com/RoboticCheese/goulash/common"
)

// APIInstance implements a struct for the API connection.
type APIInstance struct {
	Component
	BaseURL string
	Version string
}

// NewAPIInstance initializes and returns a new API instance based on a
// Supermarket URL.
func NewAPIInstance(url string) (i *APIInstance, err error) {
	i = InitAPIInstance()
	i.BaseURL = url
	i.Component, err = NewComponent(i.BaseURL)
	if err != nil {
		return
	}
	// TODO: Make the version configurable somewhere...
	i.Version = "1"
	i.Endpoint = i.BaseURL + "/api/v" + i.Version
	resp, err := http.Get(i.BaseURL + "/status")
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}
	return
}

// InitAPIInstance generates an empty APIInstance struct.
func InitAPIInstance() (i *APIInstance) {
	i = new(APIInstance)
	return
}

// Empty checks whether an APIInstance struct has been populated with anything
// or still holds all the base defaults.
func (a *APIInstance) Empty() (empty bool) {
	empty = common.Empty(a)
	return
}
