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

// tagsCmd represents the tags command
var initCmd = &cobra.Command{
	Use:    "init",
	Short:  "Initialize a mdb database structure in etcd",
	Hidden: true,
	Args:   cobra.NoArgs,
	Run:    initCmdRun,
}

func initCmdRun(cmd *cobra.Command, args []string) {
	// Just try to create the directories, ignore errors.
	err := createdir(kk, "/mdb")
	if err != nil && !isKeyExists(err) && !isKeyNotFile(err) {
		logerr("Couldn't init /mdb: %s\n", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	err = createdir(kk, "/mdb/endpoints")
	if err != nil && !isKeyExists(err) && !isKeyNotFile(err) {
		logerr("Couldn't init /mdb/endpoints: %s\n", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	err = createdir(kk, "/mdb/vars")
	if err != nil && !isKeyExists(err) && !isKeyNotFile(err) {
		logerr("Couldn't init /mdb/vars: %s\n", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
