package videotransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"

	"github.com/gin-gonic/gin"
)

func GetVideoById(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 	id, err := strconv.Atoi(c.Param("id"))
		// 	if err != nil {
		// 		panic(common.ErrInvalidRequest(err))
		// 	}

		// 	courseSlug := c.Param("course_slug")
		// 	if courseSlug == "" {
		// 		panic(common.ErrInvalidRequest(errors.New("missing course slug")))
		// 	}

		// 	videoStore := videostore.NewSQLStore(appCtx.GetMainDBConnection())
		// 	courseStore := coursestore.NewSQLStore(appCtx.GetMainDBConnection())
		// 	repo := videorepo.NewGetVideoRepo(videoStore, courseStore)
		// 	biz := videobiz.NewGetVideoBiz(repo)

		// 	video, err := biz.GetVideoById(c.Request.Context(), uint32(id), courseSlug)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	video.Mask(false)

		// 	c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
