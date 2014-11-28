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
Package universe implements the building blocks that make up the top-level
Universe struct.

This file defines a universe-style CookbookVersion struct.
*/
package universe

import (
	"github.com/RoboticCheese/goulash/common"
)

// CookbookVersion implements a struct for each cookbook version underneath a
// Universe.
type CookbookVersion struct {
	Version      string
	LocationType string            `json:"location_type"`
	LocationPath string            `json:"location_path"`
	DownloadURL  string            `json:"download_url"`
	Dependencies map[string]string `json:"dependencies"`
}

// NewCookbookVersion generates an empty CookbookVersion struct.
func NewCookbookVersion() (cv *CookbookVersion) {
	cv = new(CookbookVersion)
	cv.Dependencies = map[string]string{}
	return
}

// Empty checks whether a CookbookVersion struct has been populated with
// anything or still holds all the base defaults.
func (cv *CookbookVersion) Empty() (empty bool) {
	empty = common.Empty(cv)
	return
}

// Equals implements an equality test for a CookbookVersion struct
func (cv *CookbookVersion) Equals(cv2 *CookbookVersion) (res bool) {
	res = common.Equals(cv, cv2)
	return
}

// Diff returns any attributes that have been changed from one CookbookVersion
// struct to another.
func (cv *CookbookVersion) Diff(cv2 *CookbookVersion) (pos, neg *CookbookVersion) {
	ipos, ineg := common.Diff(cv, cv2, &CookbookVersion{}, &CookbookVersion{})
	if ipos != nil {
		cpos := ipos.(*CookbookVersion)
		pos = cpos
	} else {
		pos = nil
	}
	if ineg != nil {
		cneg := ineg.(*CookbookVersion)
		neg = cneg
	} else {
		neg = nil
	}
	return
}
