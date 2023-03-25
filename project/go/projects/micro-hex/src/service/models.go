package service

type AppVersion struct {
	Build   string
	Version string
}

type Config struct {
	Environment string        `json:"environment"`
	Redis       *RedisOptions `json:"redis"`
	Sql         *SqlConfig    `json:"sql"`
	Hesp        *HespConfig   `json:"hesp"`
	Version     AppVersion    `json:"version"`
}

type EventItem struct {
	EventType string `json:"eventType"`
	Id        int    `json:"id"`
	RoomId    int    `json:"roomId"`
	Timestamp int64  `json:"timestamp"`
}

type SessionEvent struct {
	Items    []*EventItem `json:"items"`
	UnixTime int64        `json:"unixTime"`
}

type RedisOptions struct {
	Address         string   `json:"address"`
	Addresses       []string `json:"addresses"`
	Password        string   `json:"password"`
	EventStreamName string   `json:"eventStreamName"`
	ConsumerGroup   string   `json:"consumerGroup"`
	ConsumerPrefix  string   `json:"consumerPrefix"`
	ConsumerName    string   `json:"consumerName,omitempty"`
}

type HespConfigData struct {
	Id        int    `json:"id"`
	StreamKey string `json:"streamKey"`
	Player    string `json:"player"`
}

type HespConfig struct {
	Api       string `json:"api"`
	ClientKey string `json:"clientKey"`
}
type SqlConfig struct {
	ConnectionString string `json:"connectionString"`
}
