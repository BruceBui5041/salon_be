package middleware

import (
	"errors"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
)

func PublicVideoCheck(appCtx component.AppContext) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		// TODO: because get playlist using param but get segment using query

		videoIdFromCtx := c.Param("video_id")
		if videoIdFromCtx == "" {
			videoIdFromCtx = c.Query("video_id")
		}

		uid, err := common.FromBase58(videoIdFromCtx)
		if err != nil {
			panic(err)
		}
		videoId := uid.GetLocalID()

		courseSlug := c.Param("course_slug")
		if courseSlug == "" {
			courseSlug = c.Query("course_slug")
		}

		videoStore := videostore.NewSQLStore(appCtx.GetMainDBConnection())
		courseStore := coursestore.NewSQLStore(appCtx.GetMainDBConnection())
		repo := videorepo.NewGetVideoRepo(videoStore, courseStore)
		biz := videobiz.NewGetVideoBiz(repo)

		video, err := biz.GetVideoById(c.Request.Context(), videoId, courseSlug)
		if err != nil {
			panic(common.ErrInternal(err))
		}

		if !video.AllowPreview && video.Course.IntroVideoID != &videoId && !video.Lesson.AllowPreview {
			panic(common.ErrNoPermission(errors.New("not public video")))
		}

		c.Next()
	}
}
