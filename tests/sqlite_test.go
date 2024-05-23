package tests

import (
	"fmt"
	"log"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/mathesukkj/gourl-shortener/app"
	"github.com/mathesukkj/gourl-shortener/sqlite"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := sqlite.NewURLService(db)

	tests := []struct {
		name       string
		initialUrl string
		wantErr    bool
	}{
		{
			name:       "Valid URL",
			initialUrl: "https://google.com/shshs",
			wantErr:    false,
		},
		{
			name:       "Empty URL",
			initialUrl: "",
			wantErr:    true,
		},
		{
			name:       "Invalid URL format",
			initialUrl: "htp://invalid-url",
			wantErr:    true,
		},
		{
			name:       "URL with query parameters",
			initialUrl: "https://example.com/path?query=value",
			wantErr:    false,
		},
		{
			name:       "URL with fragment",
			initialUrl: "https://example.com/path#section",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO urls").
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			createdUrl, err := service.Create(app.URLPayload{Url: tt.initialUrl})

			assert.NoError(t, err)
			assert.NotNil(t, createdUrl)
			assert.Equal(t, tt.initialUrl, createdUrl.InitialURL)

			urlObj, err := url.Parse(createdUrl.ShortenedURL)

			assert.NoError(t, err)
			assert.Equal(t, app.BASE_URL, fmt.Sprintf("%s://%s:%s",
				urlObj.Scheme,
				urlObj.Hostname(),
				urlObj.Port(),
			))
			assert.False(t, len(urlObj.Path) < 5 && len(urlObj.Path) > 10)
		})
	}
}

func TestFindByShortened(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := sqlite.NewURLService(db)

	initialUrl := "https://google.com/abab"
	shortenedUrl := app.BASE_URL + "/abXCkeIU"

	res, err := db.Exec(`
    INSERT INTO urls (initial_url, shortened_url) 
    VALUES (?, ?)
  `,
		initialUrl,
		shortenedUrl,
	)
	assert.NoError(t, err)

	log.Println(res.LastInsertId())

	tests := []struct {
		name     string
		url      string
		foundURL *app.URL
	}{
		{
			name: "URL found",
			url:  shortenedUrl,
			foundURL: &app.URL{
				Id:           1,
				InitialURL:   initialUrl,
				ShortenedURL: shortenedUrl,
			},
		},
		{
			name:     "URL not found",
			url:      "http://sdokfe/cd",
			foundURL: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT * FROM urls")

			foundUrl, err := service.FindByShortened(tt.foundURL.ShortenedURL)
			assert.NoError(t, err)
			assert.Equal(t, tt.foundURL, foundUrl)
		})
	}
}
