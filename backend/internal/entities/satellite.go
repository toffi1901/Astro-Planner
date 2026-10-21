package entities

import "time"

type Satellite struct {
	NoradID        int64     `json:"norad_id"`
	Name           string    `json: "name"`
	Type           string    `json: "type"`
	TLELine1       string    `json: "tle_line1"`
	TLELine2       string    `json: "tle_line2"`
	LastUpdatedTLE time.Time `json: "last_updated_tle"`
	IsActive       bool      `json: "is_active"`
}

type PredictRequest struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
