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
	"log"
	"strings"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear bucket",
	Long: `Clear bucket for login, password or ip. 
IP should be of host

See examples below:
	clear --login=test        - clear bucket for login=test
	clear --password=1234     - clear bucket for password=1234
	clear --ip=193.70.18.123  - clear bucket for ip=193.70.18.123

	clear --login=test --password=1234     - clear bucket for login=test and password=1234
	clear --login=test --ip=193.70.18.123  - clear bucket for login=test and password=1234

	clear --login=test --password=1234 --ip=193.70.18.123 - clear buckets for login, password and ip
`,
	Run: func(cmd *cobra.Command, args []string) {
		login, password, ip, err := validateClearBucketArgs(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		runClearBucketCommand(login, password, ip)
	},
}

func init() {
	clearCmd.Flags().String("login", "", "--login=test")
	clearCmd.Flags().String("password", "", "--password=1234")
	clearCmd.Flags().String("ip", "", "--ip=193.70.18.123")

	rootCmd.AddCommand(clearCmd)
}

func validateClearBucketArgs(cmd *cobra.Command, args []string) (login string, password string, ip string, err error) {
	err = cmd.Flags().Parse(args)
	if err != nil {
		return "", "", "", err
	}

	login, err = cmd.Flags().GetString("login")
	if err != nil {
		return login, "", "", err
	}

	password, err = cmd.Flags().GetString("password")
	if err != nil {
		return login, password, "", err
	}

	ip, err = cmd.Flags().GetString("ip")
	if err != nil {
		return login, password, "", err
	}

	if ip == "" {
		return login, password, "", nil
	}

	_, err = entities.NewWithoutMaskPart(ip)
	if err != nil {
		return login, password, "", err
	}

	return login, password, ip, nil
}

func runClearBucketCommand(login string, password string, ip string) {
	cfg := getGRPCClientConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	_, err := cfg.apiClient.ClearBucket(ctx, &grpc.BucketRequest{
		Login:    login,
		Password: password,
		Ip:       ip,
	})

	if err != nil {
		fmt.Printf("FAIL: %s\n", err)
		return
	}

	var parts []string

	if login != "" {
		parts = append(parts, "login="+login)
	}

	if password != "" {
		parts = append(parts, "password="+password)
	}

	if ip != "" {
		parts = append(parts, "ip="+ip)
	}

	if len(parts) > 1 {
		fmt.Printf("OK: clear buckets for %s\n", strings.Join(parts, ", "))
	} else {
		fmt.Printf("OK: clear bucket for %s\n", strings.Join(parts, ", "))
	}
}
