// +build integration

package db_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vignesh-innblockchain/blockatlas/db"
	"github.com/vignesh-innblockchain/blockatlas/tests/integration/setup"
)

var database *db.Instance

func TestMain(m *testing.M) {
	database = setup.RunPgContainer()
	code := m.Run()
	setup.StopPgContainer()
	os.Exit(code)
}

func TestPgSetup(t *testing.T) {
	assert.NotNil(t, database)
	assert.NotNil(t, database.Gorm)
}
