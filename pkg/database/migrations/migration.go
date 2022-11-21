package migrations

import "golek_bookmark_service/pkg/database"

type Migration struct {
	DB *database.Database
}

func New(db *database.Database) *Migration {
	return &Migration{DB: db}
}
