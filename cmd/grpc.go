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
	"os"
	"time"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/storage/sql"

	grpcAPI "github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	DefaultGRPCPort = "50051"
	DefaultGRPCHost = "127.0.0.1"
)

// grpcCmd represents the grpcAPI command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run grpc service",
	Long:  `Run grpc service`,
	Run: func(cmd *cobra.Command, args []string) {
		runGRPCServer()
	},
}

// grpc client for client commands add, delete, clear, auth
type GRPCClientConfig struct {
	addr      string
	timeout   time.Duration
	apiClient grpcAPI.ApiClient
}

var grpcClientConfig *GRPCClientConfig

func init() {
	rootCmd.AddCommand(grpcCmd)
}

func runGRPCServer() {
	port := viper.GetString("GRPC_PORT")
	if port == "" {
		port = DefaultGRPCPort
	}

	l := logger.GetLogger()
	l.Debugf("Run grpcAPI service on port %s", port)

	dbCfg := sql.NewConfigByEnv()

	db, err := sql.Connect(dbCfg)
	if err != nil {
		l.Fatal(err)
	}

	err = grpcAPI.NewAPIByViper(viper.GetViper(), db).Run(port)
	if err != nil {
		l.Fatal(err)
	}
}

func getGRPCClientConfig() *GRPCClientConfig {
	if grpcClientConfig != nil {
		return grpcClientConfig
	}

	// first init

	l := logger.GetLogger()

	host := os.Getenv("GRPC_SERVER_HOST")
	if host == "" {
		host = DefaultGRPCHost
	}

	port := os.Getenv("GRPC_SERVER_PORT")
	if port == "" {
		port = DefaultGRPCPort
	}

	addr := host + ":" + port

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		l.Fatalf("establish connection with `%s` failed: %s", addr, err)
	}

	grpcClientConfig = &GRPCClientConfig{
		addr:      addr,
		timeout:   3 * time.Second,
		apiClient: grpcAPI.NewApiClient(conn),
	}

	return grpcClientConfig
}
