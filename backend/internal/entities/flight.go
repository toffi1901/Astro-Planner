package entities

import "time"

type Flight struct {
	NoradID      int64     `json:"norad_id"`
	SatName      string    `json:"satellite_name"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	MaxAltitude  float64   `json:"max_altitude"`  // максимальная высота над горизонтом, градусы
	StartAzimuth float64   `json:"start_azimuth"` // с какой стороны появляется
	EndAzimuth   float64   `json:"end_azimuth"`   // в какую сторону уходит
}
