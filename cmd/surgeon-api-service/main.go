package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xnemect1/surgeon-webapi/api"
	"github.com/xnemect1/surgeon-webapi/internal/surgeon_wl"
)

func main() {
    log.Printf("Server started")
    port := os.Getenv("SURGEON_API_PORT")
    if port == "" {
        port = "8080"
    }
    environment := os.Getenv("SURGEON_API_ENVIRONMENT")
    if !strings.EqualFold(environment, "production") { // case insensitive comparison
        gin.SetMode(gin.DebugMode)
    }
    engine := gin.New()
    engine.Use(gin.Recovery())
    // request routings
    surgeon_wl.AddRoutes(engine)
    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}