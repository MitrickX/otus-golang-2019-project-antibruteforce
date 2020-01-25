package logger

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func InitLogger(v *viper.Viper) {

	zapConfig := zap.NewProductionConfig()

	loggerCfg := v.GetStringMap("logger")
	if val, ok := loggerCfg["level"]; ok {
		if level, okCast := val.(string); okCast {
			// if level unrecognized, just will use NewProductionConfig level, so ignore error
			_ = zapConfig.Level.UnmarshalText([]byte(level))
		}
	}

	if val, ok := loggerCfg["output_paths"]; ok {
		if paths, ok := convertEmptyInterfaceStringSlice(val); ok {
			paths = dropOffNotExistedDirPaths(paths, []string{"stderr", "stdout"})
			zapConfig.OutputPaths = paths
		}
	}

	if val, ok := loggerCfg["error_output_paths"]; ok {
		if paths, ok := convertEmptyInterfaceStringSlice(val); ok {
			paths = dropOffNotExistedDirPaths(paths, []string{"stderr", "stdout"})
			zapConfig.ErrorOutputPaths = paths
		}
	}

	var err error
	zapLogger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("Can't init zap logger: %s\n", err)
	}

	logger = zapLogger.Sugar()
}

func convertEmptyInterfaceStringSlice(data interface{}) ([]string, bool) {
	switch res := data.(type) {
	case []string:
		return res, true
	case []interface{}:
		var values []string
		for _, v := range res {
			str, ok := v.(string)
			if !ok {
				return nil, false
			}
			values = append(values, str)

		}
		return values, true
	}
	return nil, false
}

func dropOffNotExistedDirPaths(paths []string, exclude []string) []string {
	var resultPaths []string
OUTER:
	for _, path := range paths {
		for _, excl := range exclude {
			if path == excl {
				resultPaths = append(resultPaths, path)
				continue OUTER
			}
		}
		dirPath := filepath.Dir(path)
		_, err := os.Stat(dirPath)
		if err == nil {
			resultPaths = append(resultPaths, path)
		}
	}
	return resultPaths
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		log.Fatal("Logger is not inited")
	}
	return logger
}
