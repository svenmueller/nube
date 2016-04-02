// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/svenmueller/nube/common"

	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/viper"
)

var cfgFile string
var profile string
var Cfg *viper.Viper

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nube",
	Short: "Manage Rackspace and AWS resources",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file (default is $HOME/.nube.yaml)")
	RootCmd.PersistentFlags().StringVar(&profile, "profile", "default", "Profile name. ")

	RootCmd.PersistentFlags().StringP("output", "o", "yaml", "Output format [yaml|json]")

	// rackspace
	RootCmd.PersistentFlags().StringP("rackspace-username", "", "", "Rackspace API username")
	RootCmd.PersistentFlags().StringP("rackspace-api-key", "", "", "Rackspace API key")
	RootCmd.PersistentFlags().StringP("rackspace-region", "", "LON", "Rackspace region name")

	// aws
	RootCmd.PersistentFlags().StringP("aws-access-key-id", "", "", "AWS Access Key ID")
	RootCmd.PersistentFlags().StringP("aws-secret-access-key", "", "", "AWS Secret Access Key")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(".nube") // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.SetEnvPrefix("nube")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if viper.IsSet(profile) {
		Cfg = viper.Sub(profile)
	} else {
		common.HandleError(fmt.Errorf("Profile %q not found in configuration file %q\n", profile, viper.ConfigFileUsed()), RootCmd)
		os.Exit(1)
	}

	Cfg.BindPFlag("output", RootCmd.PersistentFlags().Lookup("output"))
	Cfg.BindPFlag("rackspace-username", RootCmd.PersistentFlags().Lookup("rackspace-username"))
	Cfg.BindPFlag("rackspace-api-key", RootCmd.PersistentFlags().Lookup("rackspace-api-key"))
	Cfg.BindPFlag("rackspace-region", RootCmd.PersistentFlags().Lookup("rackspace-region"))
	Cfg.BindPFlag("aws-access-key-id", RootCmd.PersistentFlags().Lookup("aws-access-key-id"))
	Cfg.BindPFlag("aws-secret-access-key", RootCmd.PersistentFlags().Lookup("aws-secret-access-key"))

}
