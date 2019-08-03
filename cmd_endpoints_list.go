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
	"sort"

	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var endpointsListCmd = &cobra.Command{
	Use:   "list",
	Short: "list existing endpoints",
	Long:  ``,
	Args:  cobra.RangeArgs(0, 1),
	Run:   endpointsListCmdRun,
}

func endpointsListCmdRun(cmd *cobra.Command, args []string) {
	endpoints := getendpointlist(kk)

	// I hate having to do this, but I need the scope. :/
	var g glob.Glob
	var err error

	if len(args) > 0 {
		g, err = glob.Compile(args[0])
		if err != nil {
			logerr("Invalid glob: %s\n", err)
			os.Exit(MDB_EXIT_USERFAIL)
		}
	} else {
		g = glob.MustCompile("*")
	}

	// This sort would be far more efficient if we generated the whole
	// list of matches first, and _then_ sorted. But it probably makes
	// zero difference without thousands of endpoints
	sort.Strings(endpoints)
	for _, endpoint := range endpoints {
		if g.Match(endpoint) {
			fmt.Println(endpoint)
		}
	}
}

func init() {
	endpointsCmd.AddCommand(endpointsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
