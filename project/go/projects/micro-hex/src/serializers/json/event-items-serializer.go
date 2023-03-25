package json

import (
	"encoding/json"
	"hesp_channel/service"

	"github.com/pkg/errors"
)

type EventItems struct{}

func (c *EventItems) Decode(input []byte) ([]*service.EventItem, error) {
	item := make([]*service.EventItem, 0)
	if err := json.Unmarshal(input, &item); err != nil {
		return nil, errors.Wrap(err, "serializer.EventItems.Decode")
	}
	return item, nil
}

func (c *EventItems) Encode(input []*service.EventItem) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.EventItems.Encode")
	}
	return rawMsg, nil
}
