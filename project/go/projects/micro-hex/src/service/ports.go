package service

import (
	"context"
)

type BackgroundService interface {
	Start(ctx context.Context) error
}

// type Cache interface {
// 	//StringGet gets a string from the cache by the given key. Returns ErrNotFound if the key does not exist.
// 	StringGet(ctx context.Context, key string) (string, error)
// 	//StringSet sets a string value for the given key. Item will be removed from the cache after the given duration.
// 	StringSet(ctx context.Context, key string, value string, duration time.Duration) error
// }

type ConfigService interface {
	LoadConfig() (*Config, error)
}

// type EventItemRepository interface {
// 	//GetEvents returns all the events for the given Unix timestamp. Returns ErrNotFound if there are no events.
// 	GetEvents(ctx context.Context, ts int64) ([]*EventItem, error)
// 	//SetEvents sets all the events on the given Unix timestamp.
// 	SetEvents(ctx context.Context, ts int64, events []*EventItem) error
// }

type HespDataRepository interface {
	GetFromWorkout(ctx context.Context, id int) (*HespConfigData, error)
}

type LogService interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warning(args ...interface{})
	Warningf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Panic(args ...interface{})
	Panicf(template string, args ...interface{})
}
