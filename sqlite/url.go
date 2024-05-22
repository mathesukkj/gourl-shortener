package sqlite

import (
	"database/sql"

	"github.com/mathesukkj/gourl-shortener/app"
)

const BASE_URL = "http://localhost:4090/"

type URLService struct {
	db *sql.DB
}

func NewURLService(db *sql.DB) *URLService {
	return &URLService{db: db}
}

func (s URLService) Create(payload app.URLPayload) (*app.URL, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	shortenedUrl := app.ShortenUrl(payload.Url)

	id, err := createUrl(tx, payload.Url, shortenedUrl)
	if err != nil {
		return nil, err
	}

	url := &app.URL{
		Id:           id,
		InitialURL:   payload.Url,
		ShortenedURL: BASE_URL + shortenedUrl,
	}

	return url, tx.Commit()
}

func (s URLService) FindByShortened(shortenedUrl string) (*app.URL, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(`SELECT * FROM urls WHERE shortened_url = ?;`, shortenedUrl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var URL app.URL
	for rows.Next() {
		if err := rows.Scan(&URL.Id, &URL.InitialURL, &URL.ShortenedURL); err != nil {
			return nil, err
		}
	}

	return &URL, nil
}

func createUrl(tx *sql.Tx, initialUrl, shortenedUrl string) (int64, error) {
	res, err := tx.Exec(`
    INSERT INTO urls (initial_url, shortened_url) 
    VALUES (?, ?)
  `,
		initialUrl,
		shortenedUrl,
	)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}
