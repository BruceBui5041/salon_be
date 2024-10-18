package apihandler

import (
	"salon_be/component"

	"github.com/gin-gonic/gin"
)

func SegmentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// uid, err := common.FromBase58(c.Query("video_id"))
		// if err != nil {
		// 	panic(err)
		// }

		// videoId := uid.GetLocalID()
		// serviceSlug := c.Query("service_slug")
		// resolution := c.Query("resolution")
		// segmentNumber := c.Query("number")

		// if resolution == "" || segmentNumber == "" || serviceSlug == "" {
		// 	logger.AppLogger.Error(c, "Missing required parameters")
		// 	panic(errors.New("missing required parameters"))
		// }

		// appCache := appCtx.GetAppCache()
		// cacheInfo, err := appCache.GetVideoCache(c.Request.Context(), serviceSlug, c.Query("video_id"))
		// if err != nil {
		// 	logger.AppLogger.Error(c, "Error getting cached URL from DynamoDB", zap.Error(err))
		// }

		// var videoURL string
		// if cacheInfo != nil && cacheInfo.VideoURL != "" {
		// 	videoURL = cacheInfo.VideoURL
		// } else {
		// 	db := appCtx.GetMainDBConnection()
		// 	videoStore := videostore.NewSQLStore(db)
		// 	serviceStore := servicestore.NewSQLStore(db)
		// 	repo := videorepo.NewGetVideoRepo(videoStore, serviceStore)
		// 	biz := videobiz.NewGetVideoBiz(repo)

		// 	video, err := biz.GetVideoById(c.Request.Context(), uint32(videoId), serviceSlug)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	videoURL = video.VideoURL

		// 	err = appCache.SetVideoCache(c.Request.Context(), serviceSlug, *video)
		// 	if err != nil {
		// 		logger.AppLogger.Error(c.Request.Context(), "Error caching URL in DynamoDB", zap.Error(err))
		// 	}
		// }

		// key := filepath.Join(
		// 	videoURL,
		// 	resolution,
		// 	fmt.Sprintf("segment_%s.ts", segmentNumber),
		// )

		// svc := appCtx.GetS3Client()
		// vidSegment, err := storagehandler.GetFileFromCloudFrontOrS3(c, svc, appconst.AWSVideoS3BuckerName, key)
		// if err != nil {
		// 	logger.AppLogger.Error(c.Request.Context(), "Error getting segment file", zap.Error(err), zap.String("key", key))
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting segment file: %v", err)})
		// 	return
		// }
		// defer vidSegment.Close()

		// c.Header("Content-Type", "video/MP2T")

		// c.Stream(func(w io.Writer) bool {
		// 	_, err := io.Copy(w, vidSegment)
		// 	if err != nil {
		// 		logger.AppLogger.Error(c, "Error streaming segment file", zap.Error(err))
		// 		return false
		// 	}
		// 	return false
		// })
	}
}
