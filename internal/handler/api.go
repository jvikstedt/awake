package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/jvikstedt/awake/internal/job"
	"github.com/jvikstedt/awake/internal/result"
)

type Api struct {
	log           *log.Logger
	jobHandler    *job.Handler
	resultHandler *result.Handler
}

func NewApi(log *log.Logger, jobHandler *job.Handler, resultHandler *result.Handler) *Api {
	return &Api{
		log:           log,
		jobHandler:    jobHandler,
		resultHandler: resultHandler,
	}
}

func (a *Api) Handler() http.Handler {
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
			r.Post("/", a.jsonResponseHandler(a.jobHandler.Create))
			r.Get("/{id}", a.jsonResponseHandler(a.withID(a.jobHandler.GetOne)))
			r.Put("/{id}", a.jsonResponseHandler(a.withID(a.jobHandler.Update)))
			r.Delete("/{id}", a.jsonResponseHandler(a.withID(a.jobHandler.Delete)))
		})
		r.Route("/results", func(r chi.Router) {
			r.Get("/", a.jsonResponseHandler(a.resultHandler.GetAll))
			r.Post("/", a.jsonResponseHandler(a.resultHandler.Create))
			r.Get("/{id}", a.jsonResponseHandler(a.withID(a.resultHandler.GetOne)))
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

func (a *Api) withID(handleFunc func(int, http.ResponseWriter, *http.Request) (interface{}, int, error)) func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		idStr := chi.URLParam(r, "id")
		if idStr == "" {
			return struct{}{}, http.StatusNotFound, fmt.Errorf("URL param %s was empty", "id")
		}

		asInt, err := strconv.Atoi(idStr)
		if err != nil {
			return struct{}{}, http.StatusUnprocessableEntity, err
		}

		return handleFunc(asInt, w, r)
	}
}
