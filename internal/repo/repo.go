package repo

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"player/internal/config"
	"player/internal/music"
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
    CREATE TABLE IF NOT EXISTS playlist(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        song TEXT NOT NULL,
        duration INT NOT NULL
    );
    `
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(song string, duration int) error {
	_, err := r.db.Exec("INSERT INTO playlist(song, duration) values(?,?)", song, duration)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLiteRepository) ShowAllSongs() (music.Playlist, error) {
	err := r.Migrate()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query("SELECT (song, duration) FROM playlist")
	if err != nil {
		return music.Playlist{}, err
	}
	defer rows.Close()

	var all music.Playlist
	for rows.Next() {
		var song music.Song
		err = rows.Scan(&song.Name, &song.Duration)
		if err != nil {
			return all, err
		}
		all = append(all, song)
	}
	return all, nil
}

func (r *SQLiteRepository) CheckFileName(filename string) error {
	if len(filename) == 0 {
		return errors.New("invalid updated filename")
	}
	if err := r.Migrate(); err != nil {
		return err
	}
	all, err := r.ShowAllSongs()
	if err != nil {
		return err
	}
	for _, v := range all {
		if v.Name == filename {
			return nil
		}
	}
	return errors.New("file already have")
}

func (r *SQLiteRepository) Update(filename string, duration int) error {
	_, err := r.db.Exec("UPDATE playlist SET duration = ? WHERE song = ?", duration, filename)
	if err != nil {
		return err
	}
	return nil
}
