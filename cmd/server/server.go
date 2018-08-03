package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/jvikstedt/awake/internal/job"
)

type Api struct {
	log        *log.Logger
	jobHandler *job.Handler
}

func (a *Api) handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: a.log}))
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	r.Use(cors.Handler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/jobs", func(r chi.Router) {
			r.Get("/", a.jsonResponseHandler(a.jobHandler.GetAll))
			r.Get("/{id}", a.jsonResponseHandler(a.jobHandler.GetOne))
			r.Put("/{id}", a.jsonResponseHandler(a.jobHandler.Update))
		})
	})

	return r
}

func (a *Api) jsonResponseHandler(handleFunc func(http.ResponseWriter, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := handleFunc(w, r)
		if err != nil {
			a.log.Println(err)
		}
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			a.log.Printf("Could not encode response to output: %v", err)
		}
	}
}
