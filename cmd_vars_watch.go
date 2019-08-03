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
var varsWatchCmd = &cobra.Command{
	Use:   "watch <environment> <variable>",
	Short: "block until value of a variable is updated",
	Args:  cobra.ExactArgs(2),
	Run:   varsWatchCmdRun,
}

func varsWatchCmdRun(cmd *cobra.Command, args []string) {
	env := args[0]
	varname, err := standardize_varname(args[1])

	val, err := watchvar(kk, env, varname)

	if isKeyNotFound(err) {
		logwarnl("No such variable %s in environment %s", varname, env)
		os.Exit(MDB_EXIT_NOTEXISTS)
	}

	if err != nil {
		logerrl("etcd communication failure: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	fmt.Println(val)
}

func init() {
	varsCmd.AddCommand(varsWatchCmd)
}
