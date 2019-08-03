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
	"strings"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var varsSetCmd = &cobra.Command{
	Use:     "set <environment> <variable>=<value>",
	Aliases: []string{"create"},
	Short:   "set/create a global variable",
	Args:    cobra.ExactArgs(2),
	Run:     varsSetCmdRun,
}

// FIXME: Validate environment name
func varsSetCmdRun(cmd *cobra.Command, args []string) {
	environment := args[0]
	s := strings.SplitN(args[1], "=", 2)

	if len(s) > 2 {
		logerr("invalid variab/evalue specification, use: <variable>=<value>")
		os.Exit(MDB_EXIT_USERFAIL)
	}

	varname := s[0]

	if valid, err := valid_varname(varname); !valid {
		logerr("invalid variable name: %s\n", err)
		os.Exit(MDB_EXIT_USERFAIL)
	}
	varname, _ = standardize_varname(varname)

	val := ""

	if len(s) == 2 {
		val = s[1]
	}

	if valid, err := valid_varvalue(val); !valid {
		logerr("invalid variable value: %s\n", err)
		os.Exit(MDB_EXIT_USERFAIL)
	}

	err := setvar(kk, environment, varname, val)
	if err != nil {
		logerr("etcd communication failure: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}
}

func init() {
	varsCmd.AddCommand(varsSetCmd)
}
