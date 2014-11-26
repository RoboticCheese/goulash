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
Package goulash wraps all the sub-packages to over consumers of this library
a single package to import and namespace to contend with.
*/
package goulash

import (
	"github.com/RoboticCheese/goulash/apiinstance"
	"github.com/RoboticCheese/goulash/cookbook"
	"github.com/RoboticCheese/goulash/cookbookversion"
	"github.com/RoboticCheese/goulash/universe"
)

// NewAPIInstance accepts a URL string and returns a pointer to a new API
// instance and any error encountered.
func NewAPIInstance(url string) (i *apiinstance.APIInstance, err error) {
	i, err = apiinstance.New(url)
	return
}

// NewCookbook accepts a pointer to an APIInstance struct and a name string and
// returns a pointer to a new Cookbook struct and any error encountered.
func NewCookbook(i *apiinstance.APIInstance, name string) (c *cookbook.Cookbook, err error) {
	c, err = cookbook.New(i, name)
	return
}

// NewCookbookVersion accepts a pointer to a Cookbook struct and a version
// string and returns a pointer to a new CookbookVersion struct.
func NewCookbookVersion(c *cookbook.Cookbook, v string) (cv *cookbookversion.CookbookVersion, err error) {
	cv, err = cookbookversion.New(c, v)
	return
}

// NewUniverse accepts a pointer to an APIInstance struct and returns a pointer
// to a new Universe struct and any error encountered.
func NewUniverse(i *apiinstance.APIInstance) (u *universe.Universe, err error) {
	u, err = universe.New(i)
	return
}
