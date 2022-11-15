package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WakeHandle(c *gin.Context) {
	c.String(http.StatusOK, "health")
}
