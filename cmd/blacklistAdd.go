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

// blacklistAddCmd represents the blacklistAdd command
var blacklistAddCmd = &cobra.Command{
	Use:   "add <ip> [flags]",
	Short: "Add IP into black list",
	Long: `Add IP into black list. 
IP could be as of host or as of subnet

See examples below:
	add 193.70.18.0/24	- add IP of subnet, bit mask with length 24 bits 
	add 193.70.18.123	- without slash, add host ip (not ip of subnet)
`,
	DisableFlagsInUseLine: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("ip is required. Run with --help for more information")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("blacklistAdd called")
	},
}

func init() {
	blacklistCmd.AddCommand(blacklistAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blacklistAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// blacklistAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
