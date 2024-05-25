package surgeon_wl

import (
	"net/http"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_surgeon_waiting_list.go

// CreateSurgeryEntry - Saves new entry into waiting list
func (this *implSurgeriesListAPI) CreateSurgeryEntry(ctx *gin.Context) {
	updateSurgeonFunc(ctx, func(c *gin.Context, surgeon *Surgeon) (*Surgeon,  interface{},  int){
        var entry SurgeryEntry

        if err := c.ShouldBindJSON(&entry); err != nil {
            return nil, gin.H{
                "status": http.StatusBadRequest,
                "message": "Invalid request body",
                "error": err.Error(),
            }, http.StatusBadRequest
        }

        if entry.PatientId == "" {
            return nil, gin.H{
                "status": http.StatusBadRequest,
                "message": "Patient ID is required",
            }, http.StatusBadRequest
        }

        if entry.Id == "" || entry.Id == "@new" {
            entry.Id = uuid.NewString()
        }

        conflictIndx := slices.IndexFunc( surgeon.SurgeriesList, func(waiting SurgeryEntry) bool {
            return entry.Id == waiting.Id || entry.PatientId == waiting.PatientId
        })

        if conflictIndx >= 0 {
            return nil, gin.H{
                "status": http.StatusConflict,
                "message": "Entry already exists",
            }, http.StatusConflict
        }

        surgeon.SurgeriesList = append(surgeon.SurgeriesList, entry)
        
        // entry was copied by value return reconciled value from the list
        entryIndx := slices.IndexFunc( surgeon.SurgeriesList, func(waiting SurgeryEntry) bool {
            return entry.Id == waiting.Id
        })
        if entryIndx < 0 {
            return nil, gin.H{
                "status": http.StatusInternalServerError,
                "message": "Failed to save entry",
            }, http.StatusInternalServerError
        }
        return surgeon, surgeon.SurgeriesList[entryIndx], http.StatusOK
    })
}

// DeleteSurgeryEntry - Deletes specific entry
func (this *implSurgeriesListAPI) DeleteSurgeryEntry(ctx *gin.Context) {
	updateSurgeonFunc(ctx, func(c *gin.Context, surgeon *Surgeon) (*Surgeon, interface{}, int) {
        entryId := ctx.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(surgeon.SurgeriesList, func(waiting SurgeryEntry) bool {
            return entryId == waiting.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        surgeon.SurgeriesList = append(surgeon.SurgeriesList[:entryIndx], surgeon.SurgeriesList[entryIndx+1:]...)
        return surgeon, nil, http.StatusNoContent
    })
}

// GetSurgeryEntries - Provides the surgeries list
func (this *implSurgeriesListAPI) GetSurgeryEntries(ctx *gin.Context) {
	updateSurgeonFunc(ctx, func(c *gin.Context, surgeon *Surgeon) (*Surgeon, interface{}, int) {
        result := surgeon.SurgeriesList
        if result == nil {
            result = []SurgeryEntry{}
        }
        // return nil surgeon - no need to update it in db
        return nil, result, http.StatusOK
    })
}

// GetSurgeryEntry - Provides details about surgery entry
func (this *implSurgeriesListAPI) GetSurgeryEntry(ctx *gin.Context) {
	updateSurgeonFunc(ctx, func(c *gin.Context, surgeon *Surgeon) (*Surgeon, interface{}, int) {
        entryId := ctx.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(surgeon.SurgeriesList, func(waiting SurgeryEntry) bool {
            return entryId == waiting.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        // return nil surgeon - no need to update it in db
        return nil, surgeon.SurgeriesList[entryIndx], http.StatusOK
    })
}

// UpdateSurgeryEntry - Updates specific entry
func (this *implSurgeriesListAPI) UpdateSurgeryEntry(ctx *gin.Context) {
	updateSurgeonFunc(ctx, func(c *gin.Context, surgeon *Surgeon) (*Surgeon, interface{}, int) {
        var entry SurgeryEntry

        if err := c.ShouldBindJSON(&entry); err != nil {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Invalid request body",
                "error":   err.Error(),
            }, http.StatusBadRequest
        }

        entryId := ctx.Param("entryId")

        if entryId == "" {
            return nil, gin.H{
                "status":  http.StatusBadRequest,
                "message": "Entry ID is required",
            }, http.StatusBadRequest
        }

        entryIndx := slices.IndexFunc(surgeon.SurgeriesList, func(waiting SurgeryEntry) bool {
            return entryId == waiting.Id
        })

        if entryIndx < 0 {
            return nil, gin.H{
                "status":  http.StatusNotFound,
                "message": "Entry not found",
            }, http.StatusNotFound
        }

        if entry.PatientId != "" {
            surgeon.SurgeriesList[entryIndx].PatientId = entry.PatientId
        }

		if entry.SurgeonId != "" {
            surgeon.SurgeriesList[entryIndx].SurgeonId = entry.SurgeonId
        }

        if entry.Id != "" {
            surgeon.SurgeriesList[entryIndx].Id = entry.Id
        }

		if entry.Id != "" {
            surgeon.SurgeriesList[entryIndx].Id = entry.Id
        }

		if entry.Date != "" {
            surgeon.SurgeriesList[entryIndx].Date = entry.Date
        }

		if entry.SurgeryNote != "" {
            surgeon.SurgeriesList[entryIndx].SurgeryNote = entry.SurgeryNote
        }

        

        
        return surgeon, surgeon.SurgeriesList[entryIndx], http.StatusOK
    })
}