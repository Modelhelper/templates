package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	startAction  = "start"
	stopAction   = "stop"
	statusAction = "status"
)

type hespChannelService struct {
	hespConfig      *HespConfig
	eventRepository HespDataRepository
	logger          LogService
}

type statusMap struct {
	canStart bool
	canStop  bool
}

type channelStatus struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewHespChannelService(
	hespConfig *HespConfig,
	eventRepository HespDataRepository,
	logger LogService) HespChannelService {

	return &hespChannelService{
		hespConfig:      hespConfig,
		eventRepository: eventRepository,
		logger:          logger,
	}
}

func (hesp *hespChannelService) StartChannel(ctx context.Context, workoutId int) {

	hesp.performAction(ctx, workoutId, startAction)

}

func (hesp *hespChannelService) StopChannel(ctx context.Context, workoutId int) {
	hesp.performAction(ctx, workoutId, stopAction)

}

func (hesp *hespChannelService) performAction(ctx context.Context, workoutId int, action string) error {

	hesp.logger.Info(logMessageBuilder("Event message received", action, "", "", workoutId, nil))
	hespData, err := hesp.eventRepository.GetFromWorkout(ctx, workoutId)

	if err != nil {
		return err
	}

	status, err := hesp.getChannelStatus(ctx, hespData.StreamKey)

	if err != nil {
		return err
	}

	hesp.logger.Infof("status=%s", logMessageBuilder("", action, hespData.StreamKey, http.MethodPost, workoutId, nil))

	do := canDoAction(status, action)
	if do {

		err = hesp.postHespRequest(ctx, workoutId, hespData.StreamKey, action)

		if err != nil {
			return err
		}
	} else {

		hesp.logger.Infof("status=%s %s", status, logMessageBuilder("Current status of channel indicates that the action can not be performed", action, hespData.StreamKey, "", workoutId, nil))

	}

	return nil
}

func (hesp *hespChannelService) postHespRequest(ctx context.Context, workoutId int, channel string, action string) error {
	hesp.logger.Debug(logMessageBuilder("Start HESP request", action, channel, "", workoutId, nil))

	req, err := hesp.createRequestWithContext(ctx, channel, action, http.MethodPost)
	if err != nil {

		hesp.logger.Error(logMessageBuilder("Failed to create the request", action, channel, http.MethodPost, workoutId, err))
		return nil
	}

	client := http.DefaultClient

	resp, err := client.Do(req)

	if err != nil {

		hesp.logger.Error(logMessageBuilder("Failed to execute the request", action, channel, http.MethodPost, workoutId, err))
		return nil
	}

	// defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {

		hesp.logger.Infof("success=true %s", logMessageBuilder("Request to HESP was successful", action, channel, http.MethodPost, workoutId, nil))
	} else {
		hesp.logger.Infof("httpStatusCode=%d %s", resp.StatusCode, logMessageBuilder("Request to HESP was successful but with an unexpected status code", statusAction, channel, http.MethodGet, -1, nil))
	}

	return nil
}

func (hesp *hespChannelService) getChannelStatus(ctx context.Context, channel string) (string, error) {
	hesp.logger.Info(logMessageBuilder("Get channel status", statusAction, channel, "", -1, nil))

	req, err := hesp.createRequestWithContext(ctx, channel, statusAction, http.MethodGet)
	if err != nil {

		hesp.logger.Error(logMessageBuilder("Failed to create the request", "", channel, http.MethodGet, -1, err))
		return "", nil
	}

	client := http.DefaultClient
	resp, err := client.Do(req)

	if err != nil {
		hesp.logger.Error(logMessageBuilder("Failed to execute the request", statusAction, channel, http.MethodGet, -1, err))
		return "", nil
	}

	if resp.StatusCode == http.StatusOK {
		hesp.logger.Info(logMessageBuilder("Request to HESP was successful", statusAction, channel, http.MethodGet, -1, nil))

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			hesp.logger.Error(logMessageBuilder("Failed to read the response body", statusAction, channel, http.MethodGet, -1, err))

		}

		var status channelStatus

		if err := json.Unmarshal(b, &status); err != nil {
			hesp.logger.Error(logMessageBuilder("Failed to unmarshal the json response", statusAction, channel, "", -1, err))
			return "", err
		}

		hesp.logger.Infof("status=%s %s", status.Code, logMessageBuilder(status.Message, statusAction, channel, "", -1, nil))

		return status.Code, nil

	} else {
		hesp.logger.Infof("httpStatusCode=%d %s", resp.StatusCode, logMessageBuilder("Request to HESP was successful but with an unexpected status code", statusAction, channel, http.MethodGet, -1, nil))
		return "unknown", nil
	}

}

func (hesp *hespChannelService) createRequestWithContext(ctx context.Context, channel, action, httpMethod string) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s/%s", hesp.hespConfig.Api, channel, action)

	var authEncoded = make([]byte, base64.StdEncoding.EncodedLen(len(hesp.hespConfig.ClientKey)))
	base64.StdEncoding.Encode(authEncoded, []byte(hesp.hespConfig.ClientKey))
	auth := fmt.Sprintf("%s %s", "Basic", string(authEncoded))

	req, err := http.NewRequestWithContext(ctx, httpMethod, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", auth)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func canDoAction(status, action string) bool {

	statuses := buildStatusMap()

	s, found := statuses[status]

	if !found {
		return false
	}

	if action == startAction {
		return s.canStart
	} else if action == stopAction {
		return s.canStop
	}

	return false
}

func buildStatusMap() map[string]statusMap {
	m := make(map[string]statusMap)

	m["creating"] = statusMap{canStart: false, canStop: false}
	m["stopping"] = statusMap{canStart: true, canStop: false}
	m["stopped"] = statusMap{canStart: true, canStop: false}
	m["deploying"] = statusMap{canStart: false, canStop: true}
	m["starting"] = statusMap{canStart: false, canStop: true}
	m["error"] = statusMap{canStart: true, canStop: true}
	m["ingesting"] = statusMap{canStart: false, canStop: true}
	m["waiting"] = statusMap{canStart: false, canStop: true}
	m["playing"] = statusMap{canStart: false, canStop: true}
	m["deleting"] = statusMap{canStart: false, canStop: false}
	m["deleted"] = statusMap{canStart: false, canStop: false}

	return m
}

func logMessageBuilder(message, action, channel, httpMethod string, workoutId int, err error) string {
	var sb strings.Builder

	if workoutId > 0 {
		sb.WriteString(fmt.Sprintf("sessionId=%d workoutId=%[1]d ", workoutId))
	}

	if action != "" {
		sb.WriteString(fmt.Sprintf("action=%s ", action))
	}

	if channel != "" {
		sb.WriteString(fmt.Sprintf("channel=%s ", channel))
	}

	if httpMethod != "" {
		sb.WriteString(fmt.Sprintf("httpMethod=%s ", httpMethod))
	}

	if message != "" {
		sb.WriteString(fmt.Sprintf("message=\"%s\" ", message))
	}

	if err != nil {
		sb.WriteString(fmt.Sprintf("error=\"%s\" ", err))
	}

	return strings.TrimSpace(sb.String())
}
