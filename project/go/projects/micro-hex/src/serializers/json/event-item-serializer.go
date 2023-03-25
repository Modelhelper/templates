package json

import (
	"encoding/json"
	"hesp_channel/service"

	"github.com/pkg/errors"
)

type EventItem struct{}

func (c *EventItem) Decode(input []byte) (*service.EventItem, error) {
	item := &service.EventItem{}
	if err := json.Unmarshal(input, item); err != nil {
		return nil, errors.Wrap(err, "serializer.EventItem.Decode")
	}
	return item, nil
}

func (c *EventItem) Encode(input *service.EventItem) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.EventItem.Encode")
	}
	return rawMsg, nil
}
