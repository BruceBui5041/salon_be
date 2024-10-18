package coursemodel

import (
	"video_server/common"
)

type VideoLessonResonse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title"`
}

type VideoProcessInfoResponse struct {
	ProcessResolution string `json:"process_resolution"`
	ProcessState      string `json:"process_state"`
}

type CourseVideoResponse struct {
	common.SQLModel `json:",inline"`
	Title           string                     `json:"title"`
	Description     string                     `json:"description"`
	ThumbnailURL    string                     `json:"thumbnail_url"`
	Duration        int                        `json:"duration"`
	Order           int                        `json:"order"`
	AllowPreview    bool                       `json:"allow_preview"`
	Lesson          VideoLessonResonse         `json:"lesson"`
	ProcessInfos    []VideoProcessInfoResponse `json:"process_infos"`
}

func (g *CourseVideoResponse) Mask(isAdmin bool) {
	g.GenUID(common.DbTypeVideo)
}

type CourseVideosResponse struct {
	common.SQLModel `json:",inline"`
	Title           string                `json:"title"`
	Videos          []CourseVideoResponse `json:"videos"`
}

func (c *CourseVideosResponse) Mask(isAdmin bool) {
	c.GenUID(common.DbTypeCourse)
}
