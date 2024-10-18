package videotransport

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"video_server/common"
	"video_server/component"
	"video_server/component/logger"
	"video_server/model/course/coursestore"
	"video_server/model/video/videobiz"
	"video_server/model/video/videomodel"
	"video_server/model/video/videorepo"
	"video_server/model/video/videostore"
	"video_server/model/videoprocessinfo/videoprocessinfostore"
	pb "video_server/proto/video_service/video_service"
	"video_server/watermill"
	"video_server/watermill/messagemodel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateVideoHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input videomodel.CreateVideo

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		videoFile, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No video file uploaded"})
			return
		}

		thumbnailFile, err := c.FormFile("thumbnail")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No thumbnail file uploaded"})
			return
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}

		db := appCtx.GetMainDBConnection()
		svc := appCtx.GetS3Client()

		if err = db.Transaction(func(tx *gorm.DB) error {
			courseStore := coursestore.NewSQLStore(tx)
			videoStore := videostore.NewSQLStore(tx)
			videoProcessStore := videoprocessinfostore.NewSQLStore(tx)
			repo := videorepo.NewCreateVideoRepo(videoStore, courseStore, videoProcessStore, svc)
			biz := videobiz.NewCreateVideoBiz(repo)

			video, err := biz.CreateNewVideo(c.Request.Context(), &input, videoFile, thumbnailFile)
			if err != nil {
				panic(common.ErrInternal(err))
			}

			requester.Mask(false)

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(video))
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
				CourseId:          "CourseId",
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
