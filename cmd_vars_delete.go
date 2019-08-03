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
var varsDeleteCmd = &cobra.Command{
	Use:     "delete <environment> <variable>",
	Aliases: []string{"del", "rm"},
	Short:   "delete a global variable for an environment",
	Args:    cobra.ExactArgs(2),
	Run:     varsDeleteCmdRun,
}

// FIXME: verify environment name
func varsDeleteCmdRun(cmd *cobra.Command, args []string) {
	env := args[0]
	varname, _ := standardize_varname(args[1])

	err := deletevar(kk, env, varname)
	if isKeyNotFound(err) {
		loginfo("Var %s not present in environment %s\n", varname, env)
		os.Exit(MDB_EXIT_SUCCESS) // Not an error if the var doesn't exist
	}

	if err != nil {
		logerr("var delete failed: %s\n", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}
}

func init() {
	varsCmd.AddCommand(varsDeleteCmd)
}
