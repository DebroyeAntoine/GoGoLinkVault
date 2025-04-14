package handler

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{Success: true, Data: data})
}

func ErrorResponse(c *gin.Context, status int, err string) {
	c.JSON(status, Response{Success: false, Error: err})
}
