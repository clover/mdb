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
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// findCmd represents the find command
var tagsUniqueCmd = &cobra.Command{
	Use:   "unique",
	Long:  ``,
	Short: "list all tag keys from all endpoints",
	Args:  cobra.ExactArgs(0),
	Run:   tagsUniqueCmdRun,
}

func tagsUniqueCmdRun(cmd *cobra.Command, args []string) {
	endpoints := getendpointlist(kk)

	tagmap := map[string]bool{}
	for _, endpoint := range endpoints {
		err := gettaglist(kk, endpoint, "", tagmap)
		if err != nil {
			logerrl("Error retrieving tag lists: %s", err)
			os.Exit(MDB_EXIT_ETCD_COMM)
		}
	}

	tagarr := map2arr(tagmap)
	sort.Strings(tagarr)

	for _, tag := range tagarr {
		stdoutl("%s", tag)
	}

	os.Exit(MDB_EXIT_SUCCESS)
}

func init() {
	tagsCmd.AddCommand(tagsUniqueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// findCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// findCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
