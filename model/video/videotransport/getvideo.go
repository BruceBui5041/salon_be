package videotransport

import (
	"errors"
	"net/http"
	"strconv"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
)

func GetVideoById(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		courseSlug := c.Param("course_slug")
		if courseSlug == "" {
			panic(common.ErrInvalidRequest(errors.New("missing course slug")))
		}

		videoStore := videostore.NewSQLStore(appCtx.GetMainDBConnection())
		courseStore := coursestore.NewSQLStore(appCtx.GetMainDBConnection())
		repo := videorepo.NewGetVideoRepo(videoStore, courseStore)
		biz := videobiz.NewGetVideoBiz(repo)

		video, err := biz.GetVideoById(c.Request.Context(), uint32(id), courseSlug)
		if err != nil {
			panic(err)
		}

		video.Mask(false)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
	}
}
