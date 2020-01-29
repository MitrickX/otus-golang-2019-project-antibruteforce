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
	"fmt"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <kind> <ip> [flags]",
	Short: "Delete IP from black or white list",
	Long: `Delete IP from black or white list. 
IP could be as of host or as of subnet

<kind> = black | white
<ip> = host IP or subnet IP in CIDR notation

See examples below:
	delete black 193.70.18.0/24	- delete subnet IP from black list 
	delete black 193.70.18.123	- delete host IP from black list
	delete white 192.70.18.0/24	- delete subnet IP from white list
	delete white 192.70.18.123	- delete host IP from white list
`,
	DisableFlagsInUseLine: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return validateListCmdArgs(args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		runDeleteCommand(args[0], entities.IP(args[1]))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func runDeleteCommand(kind string, ip entities.IP) {
	cfg := getGRPCClientConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	var err error

	if kind == "black" {
		_, err = cfg.apiClient.DeleteFromBlackList(ctx, &grpc.IPRequest{Ip: string(ip)})
	} else {
		_, err = cfg.apiClient.DeleteFromBlackList(ctx, &grpc.IPRequest{Ip: string(ip)})
	}

	if err != nil {
		fmt.Printf("FAIL: %s\n", err)
	} else {
		fmt.Printf("OK: %s delete in %s list\n", ip, kind)
	}
}
