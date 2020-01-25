package tests

import (
	"context"
	"flag"
	"os"
	"strings"
	"time"

	grpcAPI "github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/grpc"
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const (
	DefaultGRPCPort = "50051"
)

var cfg *Config

type Config struct {
	apiClient   grpcAPI.ApiClient
	timeout     time.Duration // timeout for context when call api methods
	RunnerPaths []string
}

func init() {

	v := viper.GetViper()

	logger.InitLogger(v)
	log := logger.GetLogger()

	features := flag.String("features", "", `-features="create_event,delete_event"`)
	featuresPath := flag.String("features-path", "", `-features-path="./features/"`)

	flag.Parse()

	cfg = &Config{}

	if *features != "" {
		featureList := strings.Split(*features, ",")
		pathPrefix := "../features/"
		if *featuresPath != "" {
			pathPrefix = *featuresPath
		}
		var paths []string
		for _, f := range featureList {
			paths = append(paths, pathPrefix+f+".feature")
		}
		cfg.RunnerPaths = paths

		log.Infof("run on features %s", paths)
	} else if *featuresPath != "" {
		cfg.RunnerPaths = append(cfg.RunnerPaths, *featuresPath)
		log.Infof("run on features %", cfg.RunnerPaths)
	}

	host := os.Getenv("GRPC_SERVER_HOST")
	if host == "" {
		log.Fatalf("env var `GRPC_SERVER_HOST` is required")
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
		log.Fatalf("establish connection with `%s` failed: %s", addr, err)
	}

	cfg.timeout = 3 * time.Second
	cfg.apiClient = grpcAPI.NewApiClient(conn)

}

func GetConfig() *Config {
	return cfg
}
