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

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var tagsDeleteCmd = &cobra.Command{
	Use:     "delete <endpoint> <tag>",
	Aliases: []string{"del", "rm"},
	Short:   "delete a tag from an endpoint",
	Args:    cobra.ExactArgs(2),
	Run:     tagsDeleteCmdRun,
}

func tagsDeleteCmdRun(cmd *cobra.Command, args []string) {
	host := args[0]
	tag, _ := standardize_tagname(args[1])

	err := deletetag(kk, host, tag)
	if isKeyNotFound(err) {
		loginfo("Tag %s not present on %s\n", tag, host)
		os.Exit(MDB_EXIT_SUCCESS) // Not an error if the tag doesn't exist
	}

	if err != nil {
		logerr("tag delete failed: %s\n", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}
}

func init() {
	tagsCmd.AddCommand(tagsDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
