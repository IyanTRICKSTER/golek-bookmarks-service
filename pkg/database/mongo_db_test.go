package database

import (
	"github.com/go-playground/assert/v2"
	"golek_bookmark_service/pkg/config"
	"testing"
)

func TestMongoDB(t *testing.T) {

	//Load .Env
	cfg := config.New("../../.env")

	//Connecting Databases
	db := New(cfg)
	db.Prepare()

	assert.NotEqual(t, db.GetConnection(), nil)

}
