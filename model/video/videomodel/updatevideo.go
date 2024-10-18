package videomodel

type UpdateVideo struct {
	Title        *string `json:"title" form:"title"`
	Description  *string `json:"description" form:"description"`
	VideoURL     *string `json:"video_url" form:"video_url"`
	ThumbnailURL *string `json:"thumbnail_url" form:"thumbnail_url"`
	Duration     *int    `json:"duration" form:"duration"`
	Order        *int    `json:"order" form:"order"`
	AllowPreview *bool   `json:"allow_preview" form:"allow_preview"`
	LessonID     *uint32 `json:"lesson_id"`
	RawVideoURL  string  `json:"raw_video_url"`
}

func (UpdateVideo) TableName() string {
	return "video"
}
