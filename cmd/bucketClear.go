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

// bucketClearCmd represents the bucketDelete command
var bucketClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear bucket(s)",
	Long: `Clear bucket(s) for login(s) or/and IP(s)

IP also support subnet notation
Login and IP could be list (separated by ',')
See examples below

Example:
	clear --login=login
	clear --login=login1,login2...
	clear --ip=127.0.0.1
	clear --ip=193.198.170.0/24
	clear --ip=127.0.0.1,193.198.170.0/24
	clear --login=login1,login2 --ip=127.0.0.1,193.198.170.0/24
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("missed flags")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bucketDelete called")
	},
}

func init() {
	bucketCmd.AddCommand(bucketClearCmd)

	bucketClearCmd.PersistentFlags().String("login", "", `for which login(s) delete bucket(s) 
use ',' as separator to pass list: --login=login1,login2`)

	bucketClearCmd.PersistentFlags().String("ip", "", `for which IP(s) delete bucket(s) 
use ',' as separator to pass list: --ip=127.0.0.1,193.198.170.0/24`)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bucketClearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bucketClearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
