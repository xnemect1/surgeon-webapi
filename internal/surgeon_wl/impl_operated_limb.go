package surgeon_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_surgeon_conditions.go
func (api *implOperatedLimbAPI) GetOperatedLimbList(ctx *gin.Context) {
	limbs := []OperatedLimb{
		{Value: "Lava ruka", Code: "Left hand"},
		{Value: "Prava ruka", Code: "Right hand"},
		{Value: "Lava noha", Code: "Left leg"},
		{Value: "Prava noha", Code: "Right leg"},
		{Value: "Hlava", Code: "Head"},
		{Value: "Brucho", Code: "Body"},
    }

    ctx.JSON(http.StatusOK, limbs)
}