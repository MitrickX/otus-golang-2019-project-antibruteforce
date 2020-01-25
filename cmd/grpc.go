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
	"strconv"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/logger"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/bucket"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/memory/ip"
	"github.com/spf13/viper"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"

	"github.com/spf13/cobra"
)

const (
	DefaultLoginBucketLimit    = 10
	DefaultPasswordBucketLimit = 100
	DefaultIPBucketLimit       = 1000
	DefaultPort                = "50051"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run grpc service",
	Long:  `Run grpc service`,
	Run: func(cmd *cobra.Command, args []string) {
		runGRPC()
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
}

func newAPI() *grpc.API {
	limits := viper.GetStringMapString("limits")

	loginBucketLimit := getIntFromStringMap(limits, "login", DefaultLoginBucketLimit)
	passwordBucketLimit := getIntFromStringMap(limits, "password", DefaultPasswordBucketLimit)
	ipBucketLimit := getIntFromStringMap(limits, "ip", DefaultIPBucketLimit)

	api := &grpc.API{
		BlackList:              ip.NewList(),
		WhiteList:              ip.NewList(),
		LoginBucketsStorage:    bucket.NewStorage(),
		PasswordBucketsStorage: bucket.NewStorage(),
		IPBucketsStorage:       bucket.NewStorage(),
		LoginBucketLimit:       uint(loginBucketLimit),
		PasswordBucketLimit:    uint(passwordBucketLimit),
		IPBucketLimit:          uint(ipBucketLimit),
	}

	return api
}

func runGRPC() {
	port := viper.GetString("GRPC_PORT")
	if port == "" {
		port = DefaultPort
	}

	l := logger.GetLogger()
	l.Debugf("Run grpc service on port %s", port)

	err := newAPI().Run(port)
	if err != nil {
		l.Error(err)
	}
}

func getIntFromStringMap(m map[string]string, key string, defaultVal int) int {
	val, ok := m[key]
	if !ok {
		return defaultVal
	}
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return valInt
}
