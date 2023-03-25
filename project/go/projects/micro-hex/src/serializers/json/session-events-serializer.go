package json

import (
	"encoding/json"
	"hesp_channel/service"

	"github.com/pkg/errors"
)

type SessionEvents struct{}

func (c *SessionEvents) Decode(input []byte) ([]*service.SessionEvent, error) {
	item := make([]*service.SessionEvent, 0)
	if err := json.Unmarshal(input, &item); err != nil {
		return nil, errors.Wrap(err, "serializer.SessionEvents.Decode")
	}
	return item, nil
}

func (c *SessionEvents) Encode(input []*service.SessionEvent) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.SessionEvents.Encode")
	}
	return rawMsg, nil
}
