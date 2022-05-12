package setup

import (
	"fmt"
	"log"

	"github.com/ory/dockertest"
	"github.com/vignesh-innblockchain/blockatlas/db"
	"github.com/vignesh-innblockchain/blockatlas/db/models"
	"gorm.io/gorm"
)

const (
	pgUser = "user"
	pgPass = "pass"
	pgDB   = "blockatlas"
)

var (
	pgResource     *dockertest.Resource
	pgContainerENV = []string{
		"POSTGRES_USER=" + pgUser,
		"POSTGRES_PASSWORD=" + pgPass,
		"POSTGRES_DB=" + pgDB,
	}

	tables = []interface{}{
		&models.Tracker{},
		&models.Asset{},
		&models.Subscription{},
		&models.SubscriptionsAssetAssociation{},
	}

	url string
)

func runPgContainerAndInitConnection() (*db.Instance, error) {
	pool := runPgContainer()
	var (
		dbConn *db.Instance
		err    error
	)
	err = pool.Retry(func() error {
		dbConn, err = db.New(url, false)
		return err
	})
	if err != nil {
		return nil, err
	}
	autoMigrate(dbConn.Gorm)
	return dbConn, nil
}

func CleanupPgContainer(dbConn *gorm.DB) {
	if err := dbConn.Migrator().DropTable(tables...); err != nil {
		log.Fatal(err)
	}
	autoMigrate(dbConn)
}

func autoMigrate(dbConn *gorm.DB) {
	if err := dbConn.AutoMigrate(tables...); err != nil {
		log.Fatal(err)
	}
}

func stopPgContainer() error {
	return pgResource.Close()
}

func runPgContainer() *dockertest.Pool {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	pgResource, err = pool.Run("postgres", "latest", pgContainerENV)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	url = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		pgUser, pgPass, pgResource.GetPort("5432/tcp"), pgDB,
	)
	return pool
}
