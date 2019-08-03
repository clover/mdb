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

// createCmd represents the create command
var endpointsCreateCmd = &cobra.Command{
	Use:   "create <endpoint>",
	Short: "create a new endpoint",
	Args:  cobra.ExactArgs(1),
	Run:   endpointsCreateCmdRun,
}

func endpointsCreateCmdRun(cmd *cobra.Command, args []string) {
	endpoint := args[0]
	args = args[1:]

	if valid, err := valid_endpointname(endpoint); !valid {
		logerr("invalid endpoint name: %s\n", err)
		os.Exit(MDB_EXIT_USERFAIL)
	}

	path := fmt.Sprintf("/mdb/endpoints/%s", endpoint)
	n, _ := checknode(kk, path)
	if n {
		logwarn("endpoint already exists: %s\n", endpoint)
		os.Exit(MDB_EXIT_EXISTS)
	}

	err := createdir(kk, path)
	if err != nil {
		logerr("etcd communication failure: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	err = createdir(kk, fmt.Sprintf("%s/tags", path))
	if err != nil {
		logerr("etcd communication failure: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	return
}

func init() {
	endpointsCmd.AddCommand(endpointsCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
