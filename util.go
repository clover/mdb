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

func map2arr(strmap map[string]bool) []string {
	out := []string{}

	for k, _ := range strmap {
		out = append(out, k)
	}

	return out
}

const (
	MDB_EXIT_SUCCESS = 0

	MDB_EXIT_FAILURE  = 1 // Generic failure
	MDB_EXIT_USERFAIL = 2 // Bad user input

	MDB_EXIT_EXISTS    = 10 // Node exists already (where we weren't expecting)
	MDB_EXIT_NOTEXISTS = 11 // Node doesn't exist (where we were expecting)

	MDB_EXIT_ETCD_COMM = 20 // Talking to etcd failed in unspecified way
)
