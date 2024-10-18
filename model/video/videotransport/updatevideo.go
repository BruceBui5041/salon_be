package videotransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/video/videobiz"
	"video_server/model/video/videomodel"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateVideoHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			apperr := common.ErrInvalidRequest(err)
			c.JSON(http.StatusBadRequest, gin.H{"message": apperr.Message})
			panic(apperr)
		}

		id := uid.GetLocalID()

		var input videomodel.UpdateVideo

		if err := c.ShouldBind(&input); err != nil {
			panic(err)
		}

		// videoFile, _ := c.FormFile("video")
		// thumbnailFile, _ := c.FormFile("thumbnail")

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		userId := requester.GetUserId()

		svc := appCtx.GetS3Client()
		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			store := videostore.NewSQLStore(tx)
			repo := videorepo.NewUpdateVideoRepo(store, svc)
			biz := videobiz.NewUpdateVideoBiz(repo)

			// video, err := biz.UpdateVideo(c.Request.Context(), id, &input, videoFile, thumbnailFile, userId)
			video, err := biz.UpdateVideo(c.Request.Context(), id, &input, nil, nil, userId)
			if err != nil {
				return err
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
			return nil
		}); err != nil {
			panic(err)
		}

	}
}
