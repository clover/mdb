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

var reVarName = regexp.MustCompile(`(?i)^[0-9a-z._-]+$`)

func valid_varname(n string) (bool, error) {
	if !reVarName.MatchString(n) {
		// FIXME: Make this more specific in the future
		return false, errors.New("valid characters must match: [0-9a-z._-]")
	}

	return true, nil
}

var reVarValue = regexp.MustCompile(`(?i)^[\w.,=/-]*$`)

func valid_varvalue(v string) (bool, error) {
	if !reVarValue.MatchString(v) {
		// FIXME: Make this more specific in the future
		return false, errors.New("valid characters must match: [\\w.,=/-]")
	}

	return true, nil
}

// convert to a standard form
//
// FIXME: Should we have this call valid_varname, or do that separately?
//
// For right now, it's just "force to lowercase"
func standardize_varname(n string) (string, error) {
	return strings.ToLower(n), nil
}

// FIXME: validate environment name
// FIXME: How should this report errors? Or differentiate between var
// missing and var empty?
func getvar(k client.KeysAPI, env string, varname string) (string, error) {
	var path string
	path = fmt.Sprintf("/mdb/vars/%s/%s", env, varname)

	val, err := getnode(k, path)
	return val, err
}

func getvarlist(k client.KeysAPI, env string) []string {
	var path string
	path = fmt.Sprintf("/mdb/vars/%s/", env)

	nl, err := getfilelist(k, path) // fix this to use var not path?
	check(err)

	varnames := []string{}

	for _, t := range nl {
		varnames = append(varnames, t)
	}

	return varnames
}

func setvar(k client.KeysAPI, env string, varname string, varval string) error {
	var path string
	path = fmt.Sprintf("/mdb/vars/%s/%s", env, varname)

	err := setnode(k, path, varval, "")

	return err
}

func deletevar(k client.KeysAPI, env string, varname string) error {
	var path string
	path = fmt.Sprintf("/mdb/vars/%s/%s", env, varname)

	err := deletenode(k, path, false)

	return err
}

func watchvar(k client.KeysAPI, env string, varname string) (string, error) {
	path := fmt.Sprintf("/mdb/vars/%s/%s", env, varname)

	return watch(k, path)
}
