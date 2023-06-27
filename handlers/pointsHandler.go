package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetPointsResponse struct {
	Points int `json:"points"`
}

func GetPoints(ctx *gin.Context) {
	id, paramExists := ctx.Params.Get("id")
	if !paramExists {
		ctx.String(http.StatusBadRequest, "Id parameter is missing")
		return
	}
	value, exists := PointsMap[id]
	if !exists {
		ctx.String(http.StatusBadRequest, "Id does not exist.")
		return
	}
	ctx.JSON(http.StatusOK, GetPointsResponse{
		Points: value,
	})
}
