package config

import (
	"fmt"
	js "{{ .Name | kebab }}/serializers/json"
	"{{ .Name | kebab }}/service"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type configService struct{}

func NewConfigService() service.ConfigService {
	return &configService{}
}

func (c *configService) serializer() service.ConfigSerializer {
	return &js.Config{}
}

func (c *configService) LoadConfig() (*service.Config, error) {
	var cfg *service.Config

	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	configFile := "config.json"
	if env == "local" {
		configFile = "config.local.json"
	}
	cfgFileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "Ports.Config.LoadConfig")
	}

	cfg, err = c.serializer().Decode(cfgFileContent)
	if err != nil {
		return nil, errors.Wrap(err, "port.Config.LoadConfig")
	}

	cfg.Environment = env

	ver, err := loadVersionFile()
	if err == nil {
		cfg.Version = ver
	}

	consumerNameSuffix := os.Getenv("CONSUMER_NAME_SUFFIX")
	if consumerNameSuffix == "" {
		consumerNameSuffix = "service"
	}
	cfg.Redis.ConsumerName = fmt.Sprint(cfg.Redis.ConsumerPrefix, consumerNameSuffix)

	secrets, _ := loadSecrets()
	if cs, ok := secrets["Sql_ConnectionString"]; ok {
		cfg.Sql.ConnectionString = cs
	}
	if hs, ok := secrets["HESP_ClientKey"]; ok {
		cfg.Hesp.ClientKey = hs
	}

	return cfg, nil
}

func loadSecrets() (map[string]string, error) {
	secretDir := "/var/run/secrets"
	files, err := ioutil.ReadDir(secretDir)
	if err != nil {
		return nil, err
	}

	secrets := make(map[string]string)
	errs := make([]error, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		secretPath := filepath.Join(secretDir, name)
		content, err := ioutil.ReadFile(secretPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		secrets[name] = string(content)
	}

	numErrs := len(errs)
	if numErrs > 0 {
		return secrets, fmt.Errorf("%d errors occured while reading secret files", numErrs)
	}

	return secrets, nil
}

func loadVersionFile() (service.AppVersion, error) {
	var ver service.AppVersion

	versionFile, err := ioutil.ReadFile("version")
	if err != nil {
		return ver, err
	}

	lines := strings.Split(string(versionFile), "\n")
	for _, v := range lines {
		if !strings.Contains(v, "=") {
			continue
		}
		kvp := strings.Split(v, "=")
		key := kvp[0]
		value := kvp[1]

		switch key {
		case "version":
			ver.Version = value
			break
		case "build":
			ver.Build = value
		}
	}

	return ver, nil
}
