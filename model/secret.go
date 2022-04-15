package model

import "time"

type Secret struct {
	ID     string        `json:"id"`
	Data   string        `json:"data"`
	TTL    time.Duration `json:"ttl"`
	IsSeen bool          `json:"isSeen"`
}

type SecretResponse struct {
	Data string `json:"data"`
}
