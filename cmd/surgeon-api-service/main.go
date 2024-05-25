package main

import (
	"log"
	"os"
	"strings"

	"context"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xnemect1/surgeon-webapi/api"
	"github.com/xnemect1/surgeon-webapi/internal/db_service"
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
    corsMiddleware := cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
        AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
        ExposeHeaders:    []string{""},
        AllowCredentials: false,
        MaxAge: 12 * time.Hour,
    })
    engine.Use(corsMiddleware)

    // setup context update  middleware
    dbService := db_service.NewMongoService[surgeon_wl.Surgeon](db_service.MongoServiceConfig{})
    defer dbService.Disconnect(context.Background())
    engine.Use(func(ctx *gin.Context) {
        ctx.Set("db_service", dbService)
        ctx.Next()
    })

    // request routings
    surgeon_wl.AddRoutes(engine)
    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}