// Package handlers contains general http handlers utils
package handlers

import (
	"todoapiservice/internal/http/httpdto"

	"github.com/gin-gonic/gin"
)

func SendErrorResponse(c *gin.Context, code int) {
	c.IndentedJSON(
		code,
		httpdto.GeneralResponse{
			Status: httpdto.StatusError,
		})
}
