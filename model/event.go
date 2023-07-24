package model

import "time"

type Event struct {
	Actor      string    `json:"actor"`
	Name       string    `json:"name"`
	ObjectType string    `json:"object_type"`
	ObjectID   string    `json:"object_id"`
	Time       time.Time `json:"time"`
	Data       string    `json:"data"`
}
