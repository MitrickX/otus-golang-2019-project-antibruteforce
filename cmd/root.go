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
	"os"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/logger"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	BlackListKind = "black"
	WhiteListKind = "white"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "otus-golang-2019-project-antibruteforce",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initLogger,
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.otus-golang-2019-project-antibruteforce.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".otus-golang-2019-project-antibruteforce" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".otus-golang-2019-project-antibruteforce")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initLogger() {
	logger.InitLogger(viper.GetViper())
}

// Validate cmd args for list commands: add, delete
func validateListCmdArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("<kind> and <ip> is required. Run with --help for more information")
	}

	if args[0] != WhiteListKind && args[0] != BlackListKind {
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

func runDeleteCommand(kind, ip string) {
	cfg := getGRPCClientConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	var err error

	if kind == BlackListKind {
		_, err = cfg.apiClient.DeleteFromBlackList(ctx, &grpc.IPRequest{Ip: ip})
	} else {
		_, err = cfg.apiClient.DeleteFromWhiteList(ctx, &grpc.IPRequest{Ip: ip})
	}

	if err != nil {
		fmt.Printf("FAIL: %s\n", err)
	} else {
		fmt.Printf("OK: %s delete in %s list\n", ip, kind)
	}
}

func runAddCommand(kind, ip string) {
	cfg := getGRPCClientConfig()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	var err error

	if kind == BlackListKind {
		_, err = cfg.apiClient.AddInBlackList(ctx, &grpc.IPRequest{Ip: ip})
	} else {
		_, err = cfg.apiClient.AddInWhiteList(ctx, &grpc.IPRequest{Ip: ip})
	}

	if err != nil {
		fmt.Printf("FAIL: %s\n", err)
	} else {
		fmt.Printf("OK: %s added in %s list\n", ip, kind)
	}
}
