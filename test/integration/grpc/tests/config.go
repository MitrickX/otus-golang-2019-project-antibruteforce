// +build !unit

package tests

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"
	"time"

	grpcAPI "github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	DefaultGRPCPort    = "50051"
	DefaultCfgFilePath = "../../../../configs/config.yml"
)

var cfg *Config

type Config struct {
	grpcAPI.LimitsConfig
	apiClient   grpcAPI.ApiClient
	timeout     time.Duration // timeout for context when call api methods
	RunnerPaths []string
}

func init() {
	cfg = &Config{}

	fs := flag.NewFlagSet("integraion tests", flag.ContinueOnError)

	// where is config
	cfgPath := fs.String("config", "", `--config=<path>`)

	// what features to test
	features := fs.String("features", "", `--features="create_event,delete_event"`)

	// where is features
	featuresPath := fs.String("features-path", "", `--features-path="./features/"`)

	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Printf("Parse arguments error %s", err)
	}

	if *cfgPath == "" {
		*cfgPath = DefaultCfgFilePath
	}

	viper.SetConfigFile(*cfgPath)

	// If a logger file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}

	v := viper.GetViper()

	logger.InitLogger(v)
	l := logger.GetLogger()

	host := os.Getenv("GRPC_SERVER_HOST")
	if host == "" {
		l.Fatalf("env var `GRPC_SERVER_HOST` is required")
	}

	cfg.RunnerPaths = getRunnerPaths(*features, *featuresPath)
	l.Infof("run on features %", cfg.RunnerPaths)

	cfg.timeout = 3 * time.Second

	port := os.Getenv("GRPC_SERVER_PORT")
	if port == "" {
		port = DefaultGRPCPort
	}

	addr := host + ":" + port

	cfg.apiClient, err = newAPIClient(addr)
	if err != nil {
		l.Fatalf("establish connection with `%s` failed: %s", addr, err)
	}

	cfg.LimitsConfig = grpcAPI.NewLimitsConfigByViper(v)
}

func getRunnerPaths(features, featuresPath string) []string {
	var paths []string

	if features != "" {
		featureList := strings.Split(features, ",")

		pathPrefix := "../features/"
		if featuresPath != "" {
			pathPrefix = featuresPath
		}

		for _, f := range featureList {
			paths = append(paths, pathPrefix+f+".feature")
		}

		return paths
	}

	if featuresPath != "" {
		paths = append(paths, featuresPath)
	}

	return paths
}

func newAPIClient(addr string) (grpcAPI.ApiClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return grpcAPI.NewApiClient(conn), nil
}

func GetConfig() *Config {
	return cfg
}
