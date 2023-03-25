package service

import "testing"

func TestLogFormat(t *testing.T) {
	actual := logMessageBuilder("message", "run", "asd", "GET", 1234, nil)
	expected := "sessionId=1234 workoutId=1234 action=run channel=asd httpMethod=GET message=\"message\""

	if actual != expected {
		t.Fatalf("Expected %s but got %s", expected, actual)
	}
}

func TestCanDoAction(t *testing.T) {
	m := buildStatusMap()

	actual := m["stopped"].canStart
	expected := true

	if actual != expected {
		t.Fatalf("Expected %v but got %v", expected, actual)
	}
}
