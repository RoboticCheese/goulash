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
Package apiinstance implements the API connection consumed by other Goulash
functionality.

This file implements the APIInstance struct.
*/
package apiinstance

import (
	"errors"
	"net/http"

	"github.com/RoboticCheese/goulash/common"
	"github.com/RoboticCheese/goulash/component"
)

// APIInstance implements a struct for the API connection.
type APIInstance struct {
	component.Component
	BaseURL string
	Version string
}

// New initializes and returns a new API instance based on a Supermarket URL.
func New(url string) (i *APIInstance, err error) {
	i = NewAPIInstance()
	i.BaseURL = url
	i.Component, err = component.New(i.BaseURL)
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

// NewAPIInstance generates an empty APIInstance struct.
func NewAPIInstance() (i *APIInstance) {
	i = new(APIInstance)
	return
}

// Empty checks whether an APIInstance struct has been populated with anything
// or still holds all the base defaults.
func (a *APIInstance) Empty() (empty bool) {
	empty = common.Empty(a)
	return
}
