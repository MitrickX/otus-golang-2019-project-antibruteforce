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

	cfgPath := flag.String("config", "", `--config=<path>`)

	features := flag.String("features", "", `-features="create_event,delete_event"`)
	featuresPath := flag.String("features-path", "", `-features-path="./features/"`)

	flag.Parse()

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

		l.Infof("run on features %s", paths)
	} else if *featuresPath != "" {
		cfg.RunnerPaths = append(cfg.RunnerPaths, *featuresPath)
		l.Infof("run on features %", cfg.RunnerPaths)
	}

	host := os.Getenv("GRPC_SERVER_HOST")
	if host == "" {
		l.Fatalf("env var `GRPC_SERVER_HOST` is required")
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

	cfg.timeout = 3 * time.Second
	cfg.apiClient = grpcAPI.NewApiClient(conn)

	cfg.LimitsConfig = grpcAPI.NewLimitsConfigByViper(v)
}

func GetConfig() *Config {
	return cfg
}
