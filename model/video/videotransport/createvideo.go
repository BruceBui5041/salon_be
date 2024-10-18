package videotransport

import (
	"fmt"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	pb "salon_be/proto/video_service/video_service"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateVideoHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var input videomodel.CreateVideo

		// if err := c.ShouldBind(&input); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		// videoFile, err := c.FormFile("video")
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "No video file uploaded"})
		// 	return
		// }

		// requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		// if !ok {
		// 	panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		// }

		db := appCtx.GetMainDBConnection()
		// svc := appCtx.GetS3Client()

		if err := db.Transaction(func(tx *gorm.DB) error {

			return nil
		}); err != nil {
			panic(common.ErrInternal(err))
		}

	}
}

func CreateVideoHandlerTest(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		for i := range 100 {
			utcTime := time.Now().UTC()
			timestamp := fmt.Sprintf("%d", utcTime.UnixNano())

			videoUploadedInfo := &messagemodel.RequestProcessVideoInfo{
				RawVidS3Key:       "RawVidS3Key",
				UploadedBy:        "UploadedBy",
				ServiceId:         "ServiceId",
				VideoId:           "VideoId",
				Timestamp:         timestamp,
				RequestResolution: pb.ProcessResolution_RESOLUTION_360P.Enum(),
			}

			err := watermill.PublishVideoUploadedEvent(c.Request.Context(), appCtx, videoUploadedInfo)
			if err != nil {
				logger.AppLogger.Error(c.Request.Context(), "publish video uploaded event", zap.Error(err), zap.Int("index", i))
				panic(err)
			}
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
