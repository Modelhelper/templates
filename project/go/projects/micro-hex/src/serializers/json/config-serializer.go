package json

import (
	"encoding/json"
	"hesp_channel/service"

	"github.com/pkg/errors"
)

type Config struct{}

func (c *Config) Decode(input []byte) (*service.Config, error) {
	config := &service.Config{}
	if err := json.Unmarshal(input, config); err != nil {
		return nil, errors.Wrap(err, "serializer.Config.Decode")
	}
	return config, nil
}

func (c *Config) Encode(input *service.Config) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Config.Encode")
	}
	return rawMsg, nil
}
