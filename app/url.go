package app

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const LETTERS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const BASE_URL = "http://localhost:4090"

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
	FindByShortened(shortenedUrl string) (*URL, error)
	Create(payload URLPayload) (*URL, error)
}

func ShortenUrl(urlStr string) string {
	randomInt := rand.Intn(6) + 5

	shortenedUrl := make([]byte, randomInt)
	for i := range randomInt {
		shortenedUrl[i] = LETTERS[rand.Intn(len(LETTERS))]
	}

	return string(shortenedUrl)
}
