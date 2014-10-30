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
package api_instance

import (
	"errors"
	"github.com/RoboticCheese/goulash/common"
	"net/http"
	"reflect"
)

// Goulash implements a data structure for a Supermarket instance.
type APIInstance struct {
	common.Component
	BaseURL string
	Version string
}

// New initializes and returns a new API instance based on a Supermarket URL.
func New(url string) (i *APIInstance, err error) {
	i = new(APIInstance)
	i.BaseURL = url
	i.Component, err = common.New(i.BaseURL)
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

// Empty checks whether an APIInstance struct has been populated with anything
// or still holds all the base defaults.
func (a *APIInstance) Empty() (empty bool) {
	empty = true
	if a == nil {
		return
	}
	r := reflect.ValueOf(a).Elem()
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		switch f.Kind() {
		case reflect.String:
			if f.String() != "" {
				empty = false
				break
			}
		case reflect.Struct:
			meth := f.Addr().MethodByName("Empty")
			if !meth.Call([]reflect.Value{})[0].Bool() {
				empty = false
				break
			}
		}
	}
	return
}
