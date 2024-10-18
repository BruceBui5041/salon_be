package videotransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
)

func ListCourseVideos(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseSlug := c.Param("course_slug")

		db := appCtx.GetMainDBConnection()
		videoStore := videostore.NewSQLStore(db)
		courseStore := coursestore.NewSQLStore(db)
		repo := videorepo.NewListVideoRepo(videoStore, courseStore)

		biz := videobiz.NewListVideoBiz(repo)

		conditions := map[string]interface{}{"course_slug": courseSlug}
		videos, err := biz.ListCourseVideos(c.Request.Context(), conditions)
		if err != nil {
			panic(err)
		}

		for i := range videos {
			videos[i].Mask(false)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(videos))
	}
}
