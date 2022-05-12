package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vignesh-innblockchain/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/coin"
	"github.com/trustwallet/golibs/types"
)

// @Summary Get Collection
// @ID collection_v4
// @Description Get a collection from the address
// @Accept json
// @Produce json
// @Tags Collections
// @Param coin path string true "the coin name" default(ethereum)
// @Param owner path string true "the query address" default(0x0875BCab22dE3d02402bc38aEe4104e1239374a7)
// @Param collection_id path string true "the query collection" default(0x06012c8cf97bead5deae237070f9587f8e7a266d)
// @Success 200 {object} types.CollectionPage
// @Failure 500 {object} ErrorResponse
// @Router /v4/{coin}/collections/{owner}/collection/{collection_id} [get]
func GetCollectiblesForSpecificCollectionAndOwner(c *gin.Context, api blockatlas.CollectionsAPI) {
	collectibles, err := api.GetCollectibles(c.Param("owner"), c.Param("collection_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, &collectibles)
}

// @Description Get collection categories
// @ID collection_categories_v4
// @Summary Get list of collections from a specific coin and addresses
// @Accept json
// @Produce json
// @Tags Collections
// @Param data body string true "Payload" default({"60": ["0xb3624367b1ab37daef42e1a3a2ced012359659b0"]})
// @Success 200 {object} blockatlas.ResultsResponse
// @Router /v4/collectibles/categories [post]
func GetCollectionCategoriesFromList(c *gin.Context, apis blockatlas.CollectionsAPIs) {
	var reqs map[string][]string
	if err := c.BindJSON(&reqs); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	reqIds := []int{}
	coinIds := []int{}
	for k := range reqs {
		coinId, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		reqIds = append(reqIds, coinId)
	}

	// old iOS client requests all accounts
	if len(reqIds) > 2 {
		coinIds = append(coinIds, coin.ETHEREUM)
	} else {
		coinIds = reqIds
	}

	batch := make(types.CollectionPage, 0)
	for _, coinId := range coinIds {
		p, ok := apis[uint(coinId)]
		if !ok {
			continue
		}
		addresses := reqs[strconv.Itoa(coinId)]
		for _, address := range addresses {
			collections, err := p.GetCollections(address)
			if err != nil {
				continue
			}
			batch = append(batch, collections...)
		}
	}
	c.JSON(http.StatusOK, &batch)
}
