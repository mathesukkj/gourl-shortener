package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/mathesukkj/gourl-shortener/app"
)

type Server struct {
	router *chi.Mux

	URLService app.URLService
}

func (s *Server) Serve(port string) {
	if err := http.ListenAndServe(port, s.router); err != nil {
		log.Fatalf("error when listening on port %s\n%v", port, err)
	}
}

func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.router.Use(m)
}

func (s *Server) registerRoutes() {
	s.router.Post("/shorten-url", s.HandleUrlShortener)
	s.router.Get("/{url}", s.handleRedirectToInitialUrl)
}

func (s *Server) HandleUrlShortener(w http.ResponseWriter, r *http.Request) {
	var payload app.URLPayload

	if err := render.Bind(r, &payload); err != nil {
		render.Render(w, r, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	url, err := s.URLService.Create(payload)
	if err != nil {
		render.Render(w, r, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	render.Render(w, r, app.URLResponse{
		Url: url.ShortenedURL,
	})
}

func (s *Server) HandleRedirectToInitialUrl(w http.ResponseWriter, r *http.Request) {
	shortenedUrl := r.PathValue("url")

	url, err := s.URLService.FindByShortened(shortenedUrl)
	if err != nil {
		render.Render(w, r, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	http.Redirect(w, r, url.InitialURL, http.StatusFound)
}

func NewServer() *Server {
	server := &Server{
		router: chi.NewRouter(),
	}

	server.Use(middleware.Logger)
	server.Use(render.SetContentType(render.ContentTypeJSON))
	server.registerRoutes()

	return server
}
