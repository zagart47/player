package repo

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"player/internal/config"
)

var DB, _ = sql.Open(config.SqlConnect.DriverName, config.SqlConnect.DataSourceName)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) SQLiteRepository {
	return SQLiteRepository{
		db: db,
	}
}

// Migrate prepares the table of database for work.
func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS songs(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        song_name TEXT NOT NULL,
        duration INT NOT NULL,
    );
    `
	_, err := r.db.Exec(query)
	return err
}
