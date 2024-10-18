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

		// 	serviceSlug := c.Param("service_slug")
		// 	if serviceSlug == "" {
		// 		panic(common.ErrInvalidRequest(errors.New("missing service slug")))
		// 	}

		// 	videoStore := videostore.NewSQLStore(appCtx.GetMainDBConnection())
		// 	serviceStore := servicestore.NewSQLStore(appCtx.GetMainDBConnection())
		// 	repo := videorepo.NewGetVideoRepo(videoStore, serviceStore)
		// 	biz := videobiz.NewGetVideoBiz(repo)

		// 	video, err := biz.GetVideoById(c.Request.Context(), uint32(id), serviceSlug)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	video.Mask(false)

		// 	c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
