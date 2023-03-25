package worker

import (
	"context"
	"hesp_channel/service"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
)

const (
	startChannelEvent = "start_live_broadcast"
	stopChannelevent  = "stop_live_broadcast"
)

type message struct {
	id        string
	group     string
	eventType string
}

type streamEventConsumer struct {
	config      *service.Config
	service     service.HespChannelService
	redisClient *redis.ClusterClient
	logger      service.LogService
}

func NewStreamEventConsumer(
	config *service.Config,
	redisClient *redis.ClusterClient,
	logger service.LogService,
	service service.HespChannelService,
) service.BackgroundService {

	return &streamEventConsumer{
		config:      config,
		service:     service,
		logger:      logger,
		redisClient: redisClient,
	}
}

// Start implements the BackgroundWorker interface
func (r *streamEventConsumer) Start(ctx context.Context) error {
	streamName := r.config.Redis.EventStreamName
	return r.consumeFromStream(ctx, streamName)
}

func (r *streamEventConsumer) consumeFromStream(ctx context.Context, streamName string) error {
	var (
		checkPending = true
		err          error
		isExiting    bool
		lastPELId    = "0"
		// wg           sync.WaitGroup
	)

	go func() {
		<-ctx.Done()
		isExiting = true
	}()

	consumerGroup := r.config.Redis.ConsumerGroup
	err = r.redisClient.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		r.logger.Errorf("message=\"Could not connect to stream\" error=\"%s\"", err)
	}

	for {
		if isExiting {
			break
		}

		var fromMessageId string
		if checkPending {
			fromMessageId = lastPELId
		} else {
			fromMessageId = ">"
		}

		entries, err := r.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: r.config.Redis.ConsumerName,
			Streams:  []string{streamName, fromMessageId},
			Count:    1,
			Block:    2 * time.Second,
			NoAck:    false,
		}).Result()
		if err != nil && err.Error() != "redis: nil" {
			return err
		}
		if entries == nil || entries[0].Messages == nil || len(entries[0].Messages) == 0 {
			if checkPending {
				checkPending = false
			}
			continue
		}

		for i := 0; i < len(entries[0].Messages); i++ {

			message := entries[0].Messages[i]

			if checkPending {
				lastPELId = message.ID
			}

			msg := toMessage(entries[0].Messages[i].Values)

			if msg != nil {
				id, _ := strconv.Atoi(msg.id)

				switch msg.eventType {
				case startChannelEvent:
					go r.service.StartChannel(ctx, id)
				case stopChannelevent:
					go r.service.StopChannel(ctx, id)
				}
			}

			r.redisClient.XAck(ctx, streamName, consumerGroup, message.ID)
			// wg.Done()
		}
	}

	return nil
}

func toMessage(values map[string]interface{}) *message {
	msg := message{}

	msg.id = values["id"].(string)

	msg.group = values["group"].(string)
	msg.eventType = values["type"].(string)

	return &msg
}
