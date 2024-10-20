package server

import (
	"fmt"
	"log"

	"salon_be/component"
	"salon_be/component/routes"
	"salon_be/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func StartHTTPServer(appCtx component.AppContext) {
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	config.AllowOrigins = []string{"http://localhost:8200"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.MaxAge = 300

	r.Use(cors.New(config))
	r.Use(middleware.OtelTracing())
	r.Use(middleware.Recover(appCtx))

	routes.SetupRoutes(r, appCtx)

	restPort := viper.GetString("REST_PORT")
	log.Printf("Starting HTTP server on :%s", restPort)
	if err := r.Run(fmt.Sprintf(":%s", restPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
