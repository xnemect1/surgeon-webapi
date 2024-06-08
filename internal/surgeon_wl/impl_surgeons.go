package surgeon_wl

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xnemect1/surgeon-webapi/internal/db_service"
)

// Kópia zakomentovanej časti z api_ambulances.go
// CreateSurgeon - Saves new ambulance definition
func (api *implSurgeonsAPI) CreateSurgeon(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}
  
	db, ok := value.(db_service.DbService[Surgeon])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}
  
	surgeon := Surgeon{}
	err := ctx.BindJSON(&surgeon)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}
  
	if surgeon.Id == "" {
		surgeon.Id = uuid.New().String()
	}
  
	err = db.CreateDocument(ctx, surgeon.Id, &surgeon)
  
	switch err {
	case nil:
		ctx.JSON(
			http.StatusCreated,
			surgeon,
		)
	case db_service.ErrConflict:
		ctx.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Surgeon already exists",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create surgeon in database",
				"error":   err.Error(),
			},
		)
	}
}

// DeleteAmbulance - Deletes specific ambulance
func (api *implSurgeonsAPI) DeleteSurgeon(ctx *gin.Context) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}
  
	db, ok := value.(db_service.DbService[Surgeon])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}
  
	surgeonId := ctx.Param("surgeonId")
	err := db.DeleteDocument(ctx, surgeonId)
  
	switch err {
	case nil:
		ctx.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Surgeon not found",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete surgeon from database",
				"error":   err.Error(),
			})
	}
}

// GetAllSurgeons handles GET requests and retrieves all surgeons
func (this *implSurgeonsAPI) GetAllSurgeons(ctx *gin.Context) {
	fmt.Println("Function get all surgeons called")
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
			"error":   "db_service not found",
		})
		return
	}

	db, ok := value.(db_service.DbService[Surgeon])
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of type db_service.DbService",
			"error":   "cannot cast db_service context to db_service.DbService",
		})
		return
	}

	surgeons, err := db.GetAllDocuments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "Failed to retrieve surgeons",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, surgeons,)
}