package utils

import (
	"fmt"
	"net/http"
	"path/filepath"
	"salon_be/appconst"
	pb "salon_be/proto/video_service/video_service"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RenameFile(originalName, newBaseName string) string {
	ext := filepath.Ext(originalName)
	return newBaseName + ext
}

func RemoveFileExtension(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex <= 0 {
		return filename
	}
	return filename[:dotIndex]
}

func WriteServerJWTTokenCookie(ctx *gin.Context, accessToken string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     appconst.AccessTokenName,
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		// SameSite: http.SameSiteStrictMode,
		Path:    "/",
		Domain:  "localhost",
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})
}

func ClearServerJWTTokenCookie(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     appconst.AccessTokenName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		// SameSite: http.SameSiteStrictMode,
		Expires: time.Now().Add(-1 * time.Hour),
	})
}

func StringToProcessResolution(resolution string) (pb.ProcessResolution, error) {
	switch strings.ToLower(strings.TrimSpace(resolution)) {
	case "360p":
		return pb.ProcessResolution_RESOLUTION_360P, nil
	case "480p":
		return pb.ProcessResolution_RESOLUTION_480P, nil
	case "720p":
		return pb.ProcessResolution_RESOLUTION_720P, nil
	case "1080p":
		return pb.ProcessResolution_RESOLUTION_1080P, nil
	default:
		return pb.ProcessResolution_ProcessResolution_NONE, fmt.Errorf("invalid resolution: %s", resolution)
	}
}
