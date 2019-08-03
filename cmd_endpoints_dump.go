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
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var doEnc = false

// listCmd represents the list command
var endpointsDumpCmd = &cobra.Command{
	Use:   "dump <endpoint>",
	Short: "dump all info for an existing endpoint",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run:   endpointsDumpCmdRun,
}

func endpointsDumpCmdRun(cmd *cobra.Command, args []string) {
	endpoint := args[0]
	args = args[1:]

	n, err := checkendpoint(kk, endpoint)
	if err != nil {
		logerr("Couldn't find endpoint, etcd error: %s", err)
		os.Exit(MDB_EXIT_ETCD_COMM)
	}

	if n == false {
		logwarnl("No such endpoint: %s", endpoint)
		os.Exit(MDB_EXIT_NOTEXISTS)
	}

	alltags, err := getalltags(kk, endpoint)

	// For sanity checking -- if we're asking to dump and enpoint that has
	// no tags, that seems pretty weird, exit with a generic error
	if len(alltags) < 1 {
		logwarnl("endpoint %s has zero tags!?", endpoint)
		os.Exit(MDB_EXIT_FAILURE)
	}

	if doEnc == false {
		outstring, err := json.Marshal(alltags)
		if err != nil {
			logerrl("Couldn't marshal json output (this shouldn't happen): %s", err)
			os.Exit(MDB_EXIT_FAILURE)
		}
		// It's good! Ship it!
		stdout("%s", outstring)
	} else {
		var role string
		if _, ok := alltags["puppet.role"]; ok {
			role = alltags["puppet.role"]
		}

		if role == "" {
			role = "none"
		}
		role_path := strings.Replace(role, "::", "/", -1)

		// Yes, this was as painful to type as it is to read
		outmap := map[string]map[string]interface{}{
			"parameters": map[string]interface{}{
				"role":      role,
				"role_path": role_path,
				"mdb":       alltags,
			},
		}
		outstring, err := yaml.Marshal(outmap)
		if err != nil {
			logerr("Couldn't marshal yaml output (this shouldn't happen): %s\n", err)
			os.Exit(MDB_EXIT_FAILURE)
		}
		// It's good! Ship it!
		stdout("%s", outstring)
	}
}

func init() {
	endpointsCmd.AddCommand(endpointsDumpCmd)

	endpointsDumpCmd.PersistentFlags().BoolVar(&doEnc, "enc", false, "output in puppet ENC-compatable format?")
}
