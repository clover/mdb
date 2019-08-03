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

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var tagsGetCmd = &cobra.Command{
	Use:   "get <endpoint> <tag>",
	Short: "get the value of an endpoint tag",
	Args:  cobra.ExactArgs(2),
	Run:   tagsGetCmdRun,
}

func tagsGetCmdRun(cmd *cobra.Command, args []string) {
	host := args[0]
	tag, err := standardize_tagname(args[1])

	val, err := gettag(kk, host, tag)

	if isKeyNotFound(err) {
		logwarn("No such tag: %s", tag)
		os.Exit(MDB_EXIT_NOTEXISTS)
	}

	if err != nil {
		logerr("etcd communication failure: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	fmt.Println(val)
}

func init() {
	tagsCmd.AddCommand(tagsGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
