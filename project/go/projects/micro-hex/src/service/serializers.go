package service

type ConfigSerializer interface {
	Decode(input []byte) (*Config, error)
	Encode(input *Config) ([]byte, error)
}

type EventItemSerializer interface {
	Decode(input []byte) (*EventItem, error)
	Encode(input *EventItem) ([]byte, error)
}

type EventItemsSerializer interface {
	Decode(input []byte) ([]*EventItem, error)
	Encode(input []*EventItem) ([]byte, error)
}

type SessionEventSerializer interface {
	Decode(input []byte) (*HespConfigData, error)
	Encode(input *HespConfigData) ([]byte, error)
}

type SessionEventsSerializer interface {
	Decode(input []byte) ([]*SessionEvent, error)
	Encode(input []*SessionEvent) ([]byte, error)
}
