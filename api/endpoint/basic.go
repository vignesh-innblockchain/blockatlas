package endpoint

import (
	"github.com/gin-gonic/gin"
	"github.com/vignesh-innblockchain/blockatlas/internal"
	"net/http"
)

func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"build":  internal.Build,
		"date":   internal.Date,
	})
}
