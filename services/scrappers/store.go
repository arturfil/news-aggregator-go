package scrappers

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// save article
func (s *Store) SaveArticle() error {
    return nil
}

