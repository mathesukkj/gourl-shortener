package api

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

func (u ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
