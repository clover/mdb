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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mdb",
	Run: func(cmd *cobra.Command, args []string) {
		stdoutl("mdb version %s", version)
		stdoutl("Using config file: %s", viper.ConfigFileUsed())
		stdoutl("Using server: %s", viper.GetString("server"))

		stdoutl("debug: %d", viper.GetInt("debug"))
		if viper.GetInt("debug") >= LogLvlDebug {
			stdoutl("")
			stdoutl("Viper state:")
			viper.Debug()
		}
	},
}
