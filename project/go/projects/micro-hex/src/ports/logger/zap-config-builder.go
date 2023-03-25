package zap_logger

import (
	"encoding/json"
	"fmt"
	"{{ .Name | kebab }}/service"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BuildConfig(config *service.Config) (*zap.Config, error) {
	fileContent, err := ioutil.ReadFile("zap.json")
	if err != nil {
		return nil, errors.Wrap(err, "Port.ZapConfigBuilder.BuildConfig")
	}

	zapConfig := &zap.Config{}
	err = json.Unmarshal(fileContent, zapConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Port.ZapConfigBuilder.BuildConfig")
	}

	zapConfig.InitialFields["version"] = config.Version.Version
	zapConfig.InitialFields["build"] = config.Version.Build

	envLogLevel := os.Getenv("ZAP_LOG_LEVEL")
	if envLogLevel != "" {
		level := zap.NewAtomicLevel()
		switch envLogLevel {
		case "debug":
			level.SetLevel(zapcore.DebugLevel)
			break
		case "info":
			break
		case "warn":
			level.SetLevel(zapcore.WarnLevel)
			break
		case "error":
			level.SetLevel(zapcore.ErrorLevel)
			break
		default:
			return nil, fmt.Errorf("unsupported log level '%s", envLogLevel)
		}
		zapConfig.Level = level
	} else if config.Environment == "local" {
		zapConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	return zapConfig, nil
}
