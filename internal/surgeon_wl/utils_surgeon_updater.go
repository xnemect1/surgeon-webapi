package surgeon_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xnemect1/surgeon-webapi/internal/db_service"
)

type surgeonUpdater = func(
    ctx *gin.Context,
    surgeon *Surgeon,
) (updatedSurgeon *Surgeon, responseContent interface{}, status int)

func updateSurgeonFunc(ctx *gin.Context, updater surgeonUpdater) {
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

    surgeon, err := db.FindDocument(ctx, surgeonId)

    switch err {
    case nil:
        // continue
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Surgeon not found",
                "error":   err.Error(),
            },
        )
        return
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to load surgeon from database",
                "error":   err.Error(),
            })
        return
    }

    if !ok {
        ctx.JSON(
            http.StatusInternalServerError,
            gin.H{
                "status":  "Internal Server Error",
                "message": "Failed to cast surgeon from database",
                "error":   "Failed to cast surgeon from database",
            })
        return
    }

    updatedSurgeon, responseObject, status := updater(ctx, surgeon)

    if updatedSurgeon != nil {
        err = db.UpdateDocument(ctx, surgeonId, updatedSurgeon)
    } else {
        err = nil // redundant but for clarity
    }

    switch err {
    case nil:
        if responseObject != nil {
            ctx.JSON(status, responseObject)
        } else {
            ctx.AbortWithStatus(status)
        }
    case db_service.ErrNotFound:
        ctx.JSON(
            http.StatusNotFound,
            gin.H{
                "status":  "Not Found",
                "message": "Surgeon was deleted while processing the request",
                "error":   err.Error(),
            },
        )
    default:
        ctx.JSON(
            http.StatusBadGateway,
            gin.H{
                "status":  "Bad Gateway",
                "message": "Failed to update surgeon in database",
                "error":   err.Error(),
            })
    }

}