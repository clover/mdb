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

	"go.etcd.io/etcd/client"
)

// func cmdEndpointsGet(k client.KeysAPI, args []string) int {
// 	host := args[0]
// 	args = args[1:]

// 	val, err := gettag(k, host, args[0])

// 	if isKeyNotFound(err) {
// 		fmt.Println("No such tag.")
// 		return 1
// 	}

// 	if err != nil {
// 		fmt.Printf("ERROR: %s\n", err)
// 		return 10
// 	}

// 	fmt.Println(val)
// 	return 0
// }

// func cmdEndpointsSet(k client.KeysAPI, args []string) int {
// 	host := args[0]
// 	tag := args[1]
// 	val := args[2]

// 	args = args[3:]

// 	_ = settag(k, host, tag, val)

// 	return 0
// }

// // FIXME: Make sure we have enough args!
// func cmdEndpointsFind(k client.KeysAPI, args []string) int {
// 	s := strings.SplitN(args[0], "=", 2)
// 	endpoints := getendpointlist(k)

// 	tag := s[0]
// 	var searchval string
// 	if len(s) > 1 {
// 		searchval = s[1]
// 	}

// 	for _, endpoint := range endpoints {
// 		t, err := gettag(k, endpoint, tag)
// 		if isKeyNotFound(err) {
// 			fmt.Printf("Tag not present on %s\n", endpoint)
// 			continue
// 		}

// 		if len(searchval) > 0 && t != searchval {
// 			fmt.Printf("Tag present but wrong value (%s != %s) on %s\n",
// 				t, searchval, endpoint)
// 			continue
// 		}
// 		fmt.Printf("%s\n", endpoint)
// 	}

// 	return 0
// }

// // Delete a tag.
// func cmdEndpointsDelete(k client.KeysAPI, args []string) int {
// 	host := args[0]
// 	tag := args[1]
// 	args = args[2:]

// 	err := deletetag(k, host, tag)
// 	if isKeyNotFound(err) {
// 		log("Tag %s not present on %s\n", tag, host)
// 		return 0 // Not an error if the tag doesn't exist
// 	}

// 	if err != nil {
// 		fmt.Printf("ERROR: %s\n", err)
// 		return 10
// 	}

// 	return 0
// }

// Right now, just check basic character set. In the future, we probably
// want to make this more detailed -- check for things like double dots,
// the correct domain suffix, etc.
var reEndpointName = regexp.MustCompile(`(?i)^@?[0-9a-z._-]+$`)

func valid_endpointname(n string) (bool, error) {
	if !reEndpointName.MatchString(n) {
		// FIXME: Make this more specific in the future
		return false, errors.New("endpoint names must match: [0-9a-z._-]")
	}

	return true, nil
}

func getendpointlist(k client.KeysAPI) []string {
	path := "/mdb/endpoints/"
	endpoints, err := getdirlist(k, path)
	check(err)

	return endpoints
}

// FIXME: Does this even make sense?
func getendpoint(k client.KeysAPI, host string, tag string) (string, error) {
	var path string
	path = fmt.Sprintf("/mdb/endpoints/%s/tags/%s", host, tag)

	val, err := getnode(k, path)
	return val, err
}

// FIXME: Needs porting
// func getendpointlist(k client.KeysAPI, host string, tag string, tagmap map[string]bool) {
// 	var path string
// 	if host != "" {
// 		path = fmt.Sprintf("/mdb/endpoints/%s/tags/", host)
// 	} else {
// 		path = fmt.Sprintf("/mdb/templates/tags/%s/", tag)
// 	}

// 	nl := getfilelist(k, path) // fix this to use tag not path
// 	for _, t := range nl {
// 		// Have we seen it already?
// 		if _, ok := tagmap[t]; !ok {
// 			tagmap[t] = true
// 			if strings.HasPrefix(t, "@") {
// 				// Not doing recursive tags right now, commented out
// 				// t = strings.TrimPrefix(t, "@")
// 				// gettaglist(k, "", t, tagmap)
// 			}
// 		}
// 	}
// }

// FIXME: Does this even make sense?
func setendpoint(k client.KeysAPI, host string, tag string, val string) error {
	var path string
	path = fmt.Sprintf("/mdb/endpoints/%s/tags/%s", host, tag)

	err := setnode(k, path, val, "")

	return err
}

func checkendpoint(k client.KeysAPI, endpoint string) (bool, error) {
	path := fmt.Sprintf("/mdb/endpoints/%s", endpoint)
	return checknode(k, path)
}

func deleteendpoint(k client.KeysAPI, endpoint string) error {
	// Sanity check: make sure we're not just going to delete the
	// whole endpoint tree because someone passed an empty string
	if len(endpoint) < 1 {
		return errors.New("deleteendpoint called with empty endpoint")
	}

	// Not gonna blow _everything_ away, so blow something away
	var path string
	path = fmt.Sprintf("/mdb/endpoints/%s", endpoint)

	err := deletenode(k, path, true)

	return err
}
