package handler

import (
	"back/internal/middleware"
	"back/internal/service"
	"back/internal/util"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	services  *service.Service
	validator *util.Validator
}

func NewHandler(services *service.Service, validator *util.Validator) *Handler {
	return &Handler{services: services, validator: validator}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", h.SignUp)
		r.Post("/login", h.SignIn)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.UserIdentity)

		r.Route("/profile", func(r chi.Router) {
			r.Get("/", h.getProfile)
			r.Put("/username", h.updateUsername)
			r.Put("/password", h.updatePassword)
		})

		r.Route("/collections", func(r chi.Router) {
			r.Get("/", h.getAllCollections)
			r.Post("/", h.createCollection)

			r.Route("/{collectionID}", func(r chi.Router) {
				r.Put("/", h.editCollection)
				r.Delete("/", h.removeCollection)
				r.Get("/train", h.startPractise)

				r.Route("/cards", func(r chi.Router) {
					r.Post("/", h.createCard)
					r.Get("/", h.getCardsByCollectionID)
				})
			})
		})

		r.Route("/cards/{cardID}", func(r chi.Router) {
			r.Put("/", h.editCard)
			r.Delete("/", h.removeCard)
		})
	})

	return r
}
