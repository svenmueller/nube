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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "Manage AWS Route53 DNS resources",
}

func init() {

	RootCmd.AddCommand(dnsCmd)

	dnsCmd.PersistentFlags().StringP("aws-access-key-id", "", "", "AWS Access Key ID")
	dnsCmd.PersistentFlags().StringP("aws-secret-access-key", "", "", "AWS Secret Access Key")

	viper.BindPFlag("aws-access-key-id", dnsCmd.PersistentFlags().Lookup("aws-access-key-id"))
	viper.BindPFlag("aws-secret-access-key", dnsCmd.PersistentFlags().Lookup("aws-secret-access-key"))

}