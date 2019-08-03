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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var tagsListCmd = &cobra.Command{
	Use:   "list <endpoint>",
	Short: "list endpoint tags",
	Args:  cobra.ExactArgs(1),
	Run:   tagsListCmdRun,
}

func tagsListCmdRun(cmd *cobra.Command, args []string) {
	host := args[0]

	tagmap := map[string]bool{}
	err := gettaglist(kk, host, "", tagmap)
	if err != nil {
		logerr("Failed to retrieve tag list for %s: %s\n", host, err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	taglist := []string{}
	for key, _ := range tagmap {
		if !strings.HasPrefix(key, "#") {
			taglist = append(taglist, key)
		}
	}

	for _, tag := range taglist {
		fmt.Println(tag)
	}
	// log("Tags: %#v\n", taglist)

	os.Exit(MDB_EXIT_SUCCESS)
}

func init() {
	tagsCmd.AddCommand(tagsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
