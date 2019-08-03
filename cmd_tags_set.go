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
var tagsSetCmd = &cobra.Command{
	Use:     "set <endpoint> <tag>=<value>",
	Aliases: []string{"create"},
	Short:   "set/create an endpoint tag",
	Args:    cobra.ExactArgs(2),
	Run:     tagSetCmdRun,
}

func tagSetCmdRun(cmd *cobra.Command, args []string) {
	endpoint := args[0]
	s := strings.SplitN(args[1], "=", 2)

	if len(s) > 2 {
		logerr("invalid tag/value specification, use: <tag>=<value>")
		os.Exit(MDB_EXIT_USERFAIL)
	}

	tag := s[0]

	if valid, err := valid_tagname(tag); !valid {
		logerr("invalid tag name: %s\n", err)
		os.Exit(MDB_EXIT_USERFAIL)
	}
	tag, _ = standardize_tagname(tag)

	val := ""

	if len(s) == 2 {
		val = s[1]
	}

	if valid, err := valid_tagvalue(tag); !valid {
		logerr("invalid tag value: %s\n", err)
		os.Exit(MDB_EXIT_USERFAIL)
	}

	// FIXME: Should this be at a lower layer?
	n, err := checkendpoint(kk, endpoint)
	if err != nil {
		logerrl("etcd communication failure: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	if n == false {
		logerrl("endpoint %s does not exist", endpoint)
		os.Exit(MDB_EXIT_NOTEXISTS)
	}

	err = settag(kk, endpoint, tag, val)
	if err != nil {
		logerrl("etcd communication failure: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}
}

func init() {
	tagsCmd.AddCommand(tagsSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
