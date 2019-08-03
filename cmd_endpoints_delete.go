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

// createCmd represents the create command
var endpointsDeleteCmd = &cobra.Command{
	Use:   "delete <endpoint>",
	Short: "delete an existing endpoint",
	Args:  cobra.ExactArgs(1),
	Run:   endpointsDeleteCmdRun,
}

func endpointsDeleteCmdRun(cmd *cobra.Command, args []string) {
	endpoint := args[0]
	args = args[1:]

	if valid, err := valid_endpointname(endpoint); !valid {
		logerr("invalid endpoint name: %s\n", err)
		os.Exit(MDB_EXIT_USERFAIL)
	}

	// when we try to delete an endpoint and it doesn't exist, actually
	// exit with an error. Deleting endpoints is rare and controlled enough
	// that it seems like trying to delete one that's not there is
	// significant enough to error. We may want to revisit this decision
	// at some point.
	//
	// This also differs from tags, where deleting a missing tag still
	// exits with success
	n, _ := checkendpoint(kk, endpoint)
	if !n {
		logwarn("Endpoint %s does not exist\n", endpoint)
		os.Exit(MDB_EXIT_NOTEXISTS)
	}

	err := deleteendpoint(kk, endpoint)

	if err != nil {
		logerr("error deleting endpoint %s: %s", endpoint, err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	return
}

func init() {
	endpointsCmd.AddCommand(endpointsDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
