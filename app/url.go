package app

import (
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type URL struct {
	Id           int64  `json:"-"`
	InitialURL   string `json:"initial_url"`
	ShortenedURL string `json:"shortened_url"`
}

type URLPayload struct {
	Url string `json:"url" validate:"required,http_url"`
}

func (u URLPayload) Bind(r *http.Request) error {
	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		return errors.New("body failed validation: " + err.Error())
	}

	return nil
}

type URLResponse struct {
	Url string `json:"url"`
}

func (u URLResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type URLService interface {
	FindById(id int64) (*URL, error)
	FindAll() ([]*URL, error)
	Create(payload URLPayload) (*URL, error)
	Update(payload URLPayload, id int64) (*URL, error)
	Delete(id int64)
}

func ShortenUrl(urlStr string) string {
	var shortenedUrl string

	urlObj, _ := url.Parse(urlStr)
	searchableUrl := getSearchableUrlStr(urlObj)

	randomInt := rand.Intn(5) + 5
	for range randomInt {
		letterIndex := rand.Intn(len(searchableUrl))
		shortenedUrl += string(searchableUrl[letterIndex])
	}

	return shortenedUrl
}

func getSearchableUrlStr(urlObj *url.URL) string {
	url := urlObj.Hostname() + urlObj.EscapedPath() + urlObj.RawQuery

	searchableUrl := ""
	for _, v := range url {
		if unicode.IsLetter(v) {
			searchableUrl += string(v)
		}
	}

	return searchableUrl
}
