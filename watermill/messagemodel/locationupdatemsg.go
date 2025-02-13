package messagemodel

type UpdateLocationBody struct {
	Event     string  `json:"event"`
	Latitude  float64 `json:"lat" binding:"required"`
	Longitude float64 `json:"long" binding:"required"`
	Accuracy  float32 `json:"accuracy,omitempty"`
	UserId    string  `json:"user_id"`
}
