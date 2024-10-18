package apihandler

import (
	"salon_be/component"

	"github.com/gin-gonic/gin"
)

func GetPlaylistHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		// videoUID, err := common.FromBase58(c.Param("video_id"))
		// if err != nil {
		// 	panic(err)
		// }

		// videoId := videoUID.GetLocalID()
		// serviceSlug := c.Param("service_slug")
		// resolution := c.Param("resolution")
		// playlistName := c.Param("playlistName")

		// if serviceSlug == "" {
		// 	logger.AppLogger.Error(c, "Missing service slug")
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing service slug"})
		// 	return
		// }

		// appCache := appCtx.GetAppCache()
		// cacheInfo, err := appCache.GetVideoCache(c.Request.Context(), serviceSlug, c.Param("video_id"))
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

		// key := filepath.Join(videoURL, "master.m3u8")

		// if playlistName != "" {
		// 	key = filepath.Join(videoURL, resolution, playlistName)
		// }

		// svc := appCtx.GetS3Client()
		// playlist, err := storagehandler.GetFileFromCloudFrontOrS3(c.Request.Context(), svc, appconst.AWSVideoS3BuckerName, key)
		// if err != nil {
		// 	logger.AppLogger.Error(c.Request.Context(), "Error getting playlist file", zap.Error(err), zap.String("key", key))
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error getting playlist file: %v", err)})
		// 	return
		// }
		// defer playlist.Close()

		// c.Header("Content-Type", "application/vnd.apple.mpegurl")

		// c.Stream(func(w io.Writer) bool {
		// 	_, err := io.Copy(w, playlist)
		// 	if err != nil {
		// 		logger.AppLogger.Error(c, "Error streaming playlist file", zap.Error(err))
		// 		return false
		// 	}
		// 	return false
		// })
	}
}
