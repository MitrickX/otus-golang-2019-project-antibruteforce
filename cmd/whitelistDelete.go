/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// whitelistDeleteCmd represents the whitelistDelete command
var whitelistDeleteCmd = &cobra.Command{
	Use:   "delete <ip> [flags]",
	Short: "Delete IP from white list",
	Long: `Delete IP from white list. 
IP could be as of host or as of subnet

See examples below:
	delete 193.70.18.0/24	- delete IP of subnet, bit mask with length 24 bits 
	delete 193.70.18.123	- without slash, delete host ip (not ip of subnet)

Note that we delete IP as it in list - i.e.
if we has IP of subnet and IP of host from this subnet and we call delete subnet IP 
then only IP on subnet will deleted from list`,
	DisableFlagsInUseLine: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("ip is required. Run with --help for more information")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("whitelistDelete called")
	},
}

func init() {
	whitelistCmd.AddCommand(whitelistDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// whitelistDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// whitelistDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
