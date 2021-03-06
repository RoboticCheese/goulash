Goulash
=======

[![Build Status](http://img.shields.io/travis/RoboticCheese/goulash.svg)][travis]

[travis]: http://travis-ci.org/RoboticCheese/goulash

A Go client library for the Chef Supermarket API.

Requirements
============

This project currently depends only on packages in the Go standard library.

Usage
=====

To do anything, you first need to import the goulash package and set up a
connection to your API instance:

    import (
        "goulash"
    )

    i, err := goulash.NewAPIInstance("https://supermarket.chef.io") // Or your API server

That instance can then be used to examine cookbook data:

    cb, err := goulash.NewCookbook(i, "nginx") // Or your API instance and cookbook name
    fmt.Print(cb.Name)
    fmt.Print(cb.Maintainer)
    fmt.Print(cb.Description)
    fmt.Print(cb.Category)
    fmt.Print(cb.LatestVersion)
    fmt.Print(cb.ExternalURL)
    fmt.Print(cb.AverageRating)
    fmt.Print(cb.CreatedAt)
    fmt.Print(cb.UpdatedAt)
    fmt.Print(cb.Deprecated)
    fmt.Print(cb.FoodcriticFailure)
    fmt.Print(cb.Versions)
    fmt.Print(cb.Versions[0])
    fmt.Print(cb.Metrics)
    fmt.Print(cb.Metrics.Downloads)
    fmt.Print(cb.Metrics.Downloads.Total)
    fmt.Print(cb.Metrics.Downloads.Versions)
    fmt.Print(cb.Metrics.Downloads.Versions["0.1.0"]) // Or your version number
    fmt.Print(cb.Metrics.Followers)

And that cookbook can be used to examine cookbook version data:

    cv, err := goulash.NewCookbookVersion(cb, "0.1.0") // Or your cookbook and version string
    fmt.Print(cv.License)
    fmt.Print(cv.TarballFileSize)
    fmt.Print(cv.Version)
    fmt.Print(cv.AverageRating)
    fmt.Print(cv.Cookbook)
    fmt.Print(cv.File)
    fmt.Print(cv.Dependencies)
    fmt.Print(cv.Dependencies["chef"]) // Or your dependency cookbook name

The instance can also be used to examine the Berkshelf-style `universe`
endpoint:

    u, err := goulash.NewUniverse(i)
    fmt.Print(u["nginx"]["2.7.4"].LocationType)
    fmt.Print(u["nginx"]["2.7.4"].LocationPath)
    fmt.Print(u["nginx"]["2.7.4"].DownloadURL)
    fmt.Print(u["nginx"]["2.7.4"].Dependencies["apt"])

Each data structure has tests for...

***Emptiness***

Return a boolean representing whether all fields in a struct are their zero
values.

c, err := goulash.NewCookbook(api, "nginx")
empty := c.Empty()

***Equality***

Return a boolean representing whether all fields in two structs are equal.

c1, err := goulash.NewCookbook(api, "nginx")
c2, err := goulash.NewCookbook(api, "othernginx")
equal := c1.Equals(c2)

***Diffs***

Return two new instances of the type being compared that represent a positive
(any new or modified) and negative (any removed or modified) diff of each field.

c1, err := goulash.NewCookbook(api, "nginx")
c2, err := goulash.NewCookbook(api, "othernginx")
positiveDiff, negativeDiff := c1.Diff(c2)

Contributing
============

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Include unit tests with any changes and make sure they pass (`go test ./...`)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request

License & Authors
=================

- Author: Jonathan Hartman <j@p4nt5.com>

Copyright 2014, Jonathan Hartman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
