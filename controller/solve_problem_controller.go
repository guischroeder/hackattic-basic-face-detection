package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"hackattic-basic-face-detection/libs/basic_face_detection/services"
)

func SolveProblem(context *gin.Context) {
    result, err := services.SolveProblem()

    if err != nil {
        context.Error(err)
    }

    if len(context.Errors) > 0 {
        context.IndentedJSON(http.StatusBadRequest, gin.H{"errors": context.Errors})
        return
    }

    context.IndentedJSON(http.StatusOK, gin.H{"data": result})
}
