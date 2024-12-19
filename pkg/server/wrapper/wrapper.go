package wrapper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, HttpResponse{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	})
}

func OkWithNil(c *gin.Context) {
	c.JSON(http.StatusOK, HttpResponse{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    nil,
	})
}

func Fail(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, HttpResponse{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	})
	c.Abort()
}

func FailWithError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, HttpResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Data:    nil,
	})
	c.Abort()
}

func FailWithCode(c *gin.Context, code int, message string) {
	c.JSON(code, HttpResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
	c.Abort()
}
