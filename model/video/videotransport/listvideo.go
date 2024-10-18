package videotransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"

	"github.com/gin-gonic/gin"
)

func ListServiceVideos(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// serviceSlug := c.Param("service_slug")

		// db := appCtx.GetMainDBConnection()
		// videoStore := videostore.NewSQLStore(db)
		// serviceStore := servicestore.NewSQLStore(db)
		// repo := videorepo.NewListVideoRepo(videoStore, serviceStore)

		// biz := videobiz.NewListVideoBiz(repo)

		// conditions := map[string]interface{}{"service_slug": serviceSlug}
		// videos, err := biz.ListServiceVideos(c.Request.Context(), conditions)
		// if err != nil {
		// 	panic(err)
		// }

		// for i := range videos {
		// 	videos[i].Mask(false)
		// }

		// c.JSON(http.StatusOK, common.SimpleSuccessResponse(videos))
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
