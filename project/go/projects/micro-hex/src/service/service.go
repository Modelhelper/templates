package service

import (
	"context"
)

type HespChannelService interface {
	StartChannel(ctx context.Context, workoutId int)
	StopChannel(ctx context.Context, workoutId int)
}
