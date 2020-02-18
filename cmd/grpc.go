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
	// DefaultGRPCPort is default port on which GRPC server runs
	DefaultGRPCPort = "50051"
	// DefaultGRPCHost is default IP of host on which GRPC server runs
	DefaultGRPCHost = "127.0.0.1"
	// DefaultDialTimeout is timeout for client dial connection to GRPC server
	DefaultDialTimeout = 5 * time.Second
	// DefaultClientContextTimeout is context timeout for calling GRPC methods by client
	DefaultClientContextTimeout = 3 * time.Second
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

// GRPCClientConfig is config for client commands add, delete, clear, auth
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

	err = grpcAPI.NewAPIByViper(viper.GetViper(), db, l).Run(port)
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

	ctx, cancel := context.WithTimeout(context.Background(), DefaultDialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		l.Fatalf("establish connection with `%s` failed: %s", addr, err)
	}

	grpcClientConfig = &GRPCClientConfig{
		addr:      addr,
		timeout:   DefaultClientContextTimeout,
		apiClient: grpcAPI.NewApiClient(conn),
	}

	return grpcClientConfig
}
