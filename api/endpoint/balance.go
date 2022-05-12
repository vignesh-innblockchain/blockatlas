package endpoint

import (
	// "errors"
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/blockatlas"
)

// @Summary Get Balance
// @ID Balance_v2
// @Description Get Balance information
// @Accept json
// @Produce json
// @Tags Transactions
// @Param coin path string true "the coin name" default(zilliqa)
// @Param address path string true "the query address" default(850321)
// @Failure 500 {object} ErrorResponse
// @Router /v2/{coin}/balances/{address} [get]
func GetBalanceByAddress(c *gin.Context, balanceAPI blockatlas.BalanceAPI) {
	address := c.Param("address")
	if address == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(blockatlas.ErrInvalidAddr))
		return
	}

	result, err := balanceAPI.GetBalanceByAddress(address)

	if err != nil {
		switch err {
		case blockatlas.ErrInvalidAddr:
			c.AbortWithStatusJSON(
				http.StatusBadRequest,
				errorResponse(blockatlas.ErrInvalidAddr),
			)
			return
		case blockatlas.ErrNotFound:
			c.AbortWithStatusJSON(
				http.StatusNotFound,
				errorResponse(blockatlas.ErrNotFound),
			)
			return
		case blockatlas.ErrSourceConn:
			c.AbortWithStatusJSON(
				http.StatusServiceUnavailable,
				errorResponse(blockatlas.ErrSourceConn),
			)
			return
		default:
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				errorResponse(err),
			)
			return
		}
	}

	c.JSON(http.StatusOK, &result)
}
