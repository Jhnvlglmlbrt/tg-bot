package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Jhnvlglmlbrt/tg-bot/lib/e"
	"github.com/Jhnvlglmlbrt/tg-bot/storage"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName); err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

// PickRandom picks random page from storage.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove removes page from storage.
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return e.Wrap("can't remove page", err)
	}

	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return e.Wrap("can't create table", err)
	}

	return nil
}

func (s *Storage) List(ctx context.Context) (string, error) {
	q := `SELECT url FROM pages`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return "", e.Wrap("can't open database", err)
	}
	defer rows.Close()

	var urls []string

	i := 1
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return "", e.Wrap("can't scan row", err)
		}

		urls = append(urls, fmt.Sprintf("%d. %s", i, url))
		i++
	}

	if err := rows.Err(); err != nil {
		return "", e.Wrap("error iterating over rows", err)
	}

	return strings.Join(urls, "\n\n"), nil
}
