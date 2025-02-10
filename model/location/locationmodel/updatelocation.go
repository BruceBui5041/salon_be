package locationmodel

type UpdateLocationInput struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Accuracy  float32 `json:"accuracy,omitempty"`
	UserId    uint32  `json:"-"`
}
