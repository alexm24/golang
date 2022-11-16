package handler

import (
	"net/http"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"

	"github.com/alexm24/golang/internal/handler/api"
	"github.com/alexm24/golang/internal/handler/route"
	"github.com/alexm24/golang/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services}
}

func (h *Handler) InitRoutes(basePath string) http.Handler {
	logHttp := httplog.NewLogger("backend", httplog.Options{LogLevel: "warn"})
	routes := route.NewRoute(h.services)
	router := chi.NewRouter()
	router.Use(httplog.RequestLogger(logHttp))
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}).Handler)

	router.Use(middleware.SetHeader("Content-Type", "application/json"))
	router.Use(middleware.SetHeader("Accept", "application/json"))

	router.Route(path.Join("/", basePath), func(router chi.Router) {
		router.Use(middleware.NoCache)
		router.Mount("/", api.Handler(routes))
	})

	return router
}
