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
	"context"
	"errors"
	"fmt"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <kind> <ip> [flags]",
	Short: "Add IP into black or white list",
	Long: `Add IP into black or white list. 
IP could be as of host or as of subnet

<kind> = black | white
<ip> = host IP or subnet IP in CIDR notation

See examples below:
	add black 193.70.18.0/24	- add subnet IP into black list 
	add black 193.70.18.123		- add host IP into black list
	add white 192.70.18.0/24	- add subnet IP into white list
	add white 192.70.18.123		- add host IP into white list
`,
	DisableFlagsInUseLine: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return validateAddArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runAddCommand(args[0], entities.IP(args[1]))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func validateAddArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("<kind> and <ip> is required. Run with --help for more information")
	}

	if args[0] != "white" && args[0] != "black" {
		return fmt.Errorf("unpexpected <kind> value `%s`. Supports: black | white. Run with --help for more information", args[0])
	}

	if len(args) < 2 {
		return errors.New("<ip> is required. Run with --help for more information")
	}

	_, err := entities.New(args[1])
	if err != nil {
		return fmt.Errorf("%s. Run with --help for more information", err)
	}

	return nil
}

func runAddCommand(kind string, ip entities.IP) {
	cfg := getGRPCClientConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	var err error

	if kind == "black" {
		_, err = cfg.apiClient.AddInBlackList(ctx, &grpc.IPRequest{Ip: string(ip)})
	} else {
		_, err = cfg.apiClient.AddInBlackList(ctx, &grpc.IPRequest{Ip: string(ip)})
	}

	if err != nil {
		fmt.Printf("FAIL: %s\n", err)
	} else {
		fmt.Printf("OK: %s added in %s list\n", ip, kind)
	}
}
