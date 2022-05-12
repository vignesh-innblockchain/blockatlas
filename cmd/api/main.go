package main

import (
	"context"

	"github.com/vignesh-innblockchain/blockatlas/internal/metrics"

	golibsGin "github.com/trustwallet/golibs/network/gin"

	"github.com/trustwallet/golibs/network/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/vignesh-innblockchain/blockatlas/api"
	"github.com/vignesh-innblockchain/blockatlas/config"
	"github.com/vignesh-innblockchain/blockatlas/db"
	_ "github.com/vignesh-innblockchain/blockatlas/docs"
	"github.com/vignesh-innblockchain/blockatlas/internal"
	"github.com/vignesh-innblockchain/blockatlas/platform"
	"github.com/vignesh-innblockchain/blockatlas/services/tokenindexer"
)

const (
	defaultPort       = "8420"
	defaultConfigPath = "../../config.yml"
)

var (
	ctx            context.Context
	cancel         context.CancelFunc
	port, confPath string
	engine         *gin.Engine
	database       *db.Instance
	tokenIndexer   tokenindexer.Instance
)

func init() {
	port, confPath = internal.ParseArgs(defaultPort, defaultConfigPath)
	ctx, cancel = context.WithCancel(context.Background())
	var err error

	internal.InitConfig(confPath)

	if err := middleware.SetupSentry(config.Default.Sentry.DSN); err != nil {
		log.Error(err)
	}

	engine = internal.InitEngine(config.Default.Gin.Mode)
	platform.Init(config.Default.Platform)

	database, err = db.New(config.Default.Postgres.URL, config.Default.Postgres.Log)
	if err != nil {
		log.Fatal(err)
	}

	metrics.Setup(database)

	tokenIndexer = tokenindexer.Init(database)
}

func main() {
	api.SetupTokensIndexAPI(engine, tokenIndexer)
	api.SetupSwaggerAPI(engine)
	api.SetupPlatformAPI(engine)
	api.SetupMetrics(engine)

	golibsGin.SetupGracefulShutdown(ctx, port, engine)
	cancel()
}
