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

// listCmd represents the list command
var varsListCmd = &cobra.Command{
	Use:   "list <environment>",
	Short: "list variables in an environment",
	Args:  cobra.ExactArgs(1),
	Run:   varsListCmdRun,
}

func varsListCmdRun(cmd *cobra.Command, args []string) {
	env := args[0]

	varnames := getvarlist(kk, env)

	for _, varname := range varnames {
		fmt.Println(varname)
	}
	// log("Tags: %#v\n", taglist)

	os.Exit(MDB_EXIT_SUCCESS)
}

func init() {
	varsCmd.AddCommand(varsListCmd)
}
