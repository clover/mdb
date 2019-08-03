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

	"go.etcd.io/etcd/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdOpts struct {
	cfgfile     string
	environment string
	server      string
	debug       int
}

// I hate that I have to make this global, but cobra doesn't seem to have
// a way to carry variables along with it. Sigh.
var kk client.KeysAPI

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mdb",
	Short: "A brief description of your application",
	Long:  `Metadata mangler of doom`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		Loglevel += 2
		// Loglevel = viper.GetInt("debug")
		// FIXME: Should be in config file, when we have a config file
		kk = etcdConnect(viper.GetString("server"), 2379)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// cmdExecute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func cmdExecute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cmdOpts.cfgfile, "config", "c", "/etc/mdb/config.yaml", "config file (default is /etc/mdb/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&cmdOpts.server, "server", "s", "localhost", "etcd server for source of truth")
	rootCmd.PersistentFlags().CountVarP(&Loglevel, "debug", "d", "Turn on debugging")
	// rootCmd.PersistentFlags().Parse()

	viper.BindPFlags(rootCmd.PersistentFlags())

	// viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	// viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// if cmdOpts.cfgfile != "" {
	// 	// Use config file from the flag.
	// 	viper.SetConfigFile(cmdOpts.cfgfile)
	// } else {
	// 	// Find home directory.
	// 	home, err := homedir.Dir()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}

	// 	// Search config in home directory with name ".mdbconfig"
	// 	viper.AddConfigPath(home)
	// 	viper.SetConfigName(".mdbconfig")
	// }

	viper.SetConfigFile(cmdOpts.cfgfile)

	viper.SetEnvPrefix("mdb")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); viper.GetInt("debug") > 0 && err != nil {
		logwarnl("Failed to open config file %s: %s", cmdOpts.cfgfile, err)
	}
}
