package entities

type CelestialObject struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
