// Copyright 2019 Clover Network, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"go.etcd.io/etcd/client"
)

var reTagName = regexp.MustCompile(`(?i)^[0-9a-z._-]+$`)

func valid_tagname(n string) (bool, error) {
	if !reTagName.MatchString(n) {
		// FIXME: Make this more specific in the future
		return false, errors.New("valid characters must match: [0-9a-z._-]")
	}

	return true, nil
}

var reTagValue = regexp.MustCompile(`(?i)^[\w.,=/-]*$`)

func valid_tagvalue(v string) (bool, error) {
	if !reTagValue.MatchString(v) {
		// FIXME: Make this more specific in the future
		return false, errors.New("valid characters must match: [\\w.,=/-]")
	}

	return true, nil
}

// convert to a standard form
//
// FIXME: Should we have this call valid_tagname, or do that separately?
//
// For right now, it's just "force to lowercase"
func standardize_tagname(n string) (string, error) {
	return strings.ToLower(n), nil
}

// FIXME: Doesn't handle @tags
// FIXME: How should this report errors? Or differentiate between tag
// missing and tag empty?
func gettag(k client.KeysAPI, endpoint string, tag string) (string, error) {
	var path string
	path = fmt.Sprintf("/mdb/endpoints/%s/tags/%s", endpoint, tag)

	val, err := getnode(k, path)
	return val, err
}

// FIXME: Instead of having this iterate over a bunch of tags, just have etcd
// recursively give us everything under the tag path.
func getalltags(k client.KeysAPI, endpoint string) (map[string]string, error) {
	tagdump := map[string]string{}

	tagmap := map[string]bool{}
	err := gettaglist(kk, endpoint, "", tagmap)
	if err != nil {
		return map[string]string{}, err
	}

	for tag, _ := range tagmap {
		val, err := gettag(k, endpoint, tag)

		// An error is bad, because we might have created a partial list
		// which we don't want interpreted as a full list (which might
		// unconfigure things), so return an error
		if err != nil {
			return map[string]string{}, err
		}

		tagdump[tag] = val
	}

	// If we made it this far, we're done and we've not had errors
	return tagdump, nil
}

// FIXME: The semantics here suuuuuuuck
// FIXME: Make sure this actually errors if there's an etcd error at any point
func gettaglist(k client.KeysAPI, endpoint string, tag string, tagmap map[string]bool) error {
	var path string
	if endpoint != "" {
		path = fmt.Sprintf("/mdb/endpoints/%s/tags/", endpoint)
	} else {
		path = fmt.Sprintf("/mdb/templates/tags/%s/", tag)
	}

	nl, err := getfilelist(k, path) // fix this to use tag not path
	if err != nil {
		return err
	}
	check(err)

	for _, t := range nl {
		// Have we seen it already?
		if _, ok := tagmap[t]; !ok {
			tagmap[t] = true
			if strings.HasPrefix(t, "@") {
				// Not doing recursive tags right now, commented out
				// t = strings.TrimPrefix(t, "@")
				// gettaglist(k, "", t, tagmap)
			}
		}
	}

	// no error
	return nil
}

func settag(k client.KeysAPI, endpoint string, tag string, val string) error {
	var path string
	path = fmt.Sprintf("/mdb/endpoints/%s/tags/%s", endpoint, tag)

	err := setnode(k, path, val, "")

	return err
}

func deletetag(k client.KeysAPI, endpoint string, tag string) error {
	var path string
	path = fmt.Sprintf("/mdb/endpoints/%s/tags/%s", endpoint, tag)

	err := deletenode(k, path, false)

	return err
}
