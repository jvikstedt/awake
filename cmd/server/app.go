package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/jvikstedt/awake/cron"
	"github.com/jvikstedt/awake/internal/database"
	"github.com/jvikstedt/awake/internal/domain"
	"github.com/jvikstedt/awake/internal/job"
	"github.com/jvikstedt/awake/internal/plugin"
	"github.com/jvikstedt/awake/internal/runner"
)

type App struct {
	log           *log.Logger
	port          string
	srv           *http.Server
	wg            sync.WaitGroup
	config        domain.Config
	appPath       string
	scheduler     *cron.Scheduler
	runner        *runner.Runner
	db            *sqlx.DB
	jobRepository domain.JobRepository
	jobHandler    *job.Handler
}

func newApp(logger *log.Logger, port string, config domain.Config, appPath string) (*App, error) {
	db, err := database.NewDB("sqlite3", filepath.Join(appPath, "awake.db"))
	if err != nil {
		return nil, err
	}
	if err := database.EnsureTables(db); err != nil {
		return nil, err
	}

	api := &Api{}

	runner := runner.New(logger, config)
	scheduler := cron.New(logger)

	jobRepository := job.NewRepository(db)
	jobHandler := job.NewHandler(api, jobRepository, runner, scheduler)

	api.log = logger
	api.jobHandler = jobHandler

	srv := &http.Server{Addr: ":" + port, Handler: api.handler()}

	return &App{
		log:           logger,
		port:          port,
		srv:           srv,
		config:        config,
		appPath:       appPath,
		scheduler:     scheduler,
		runner:        runner,
		db:            db,
		jobRepository: jobRepository,
		jobHandler:    jobHandler,
	}, nil
}

func (a *App) startServices() error {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.runner.Start()
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.scheduler.Start()
	}()

	jobs, err := a.jobRepository.GetAll()
	if err != nil {
		return err
	}

	for _, j := range jobs {
		if j.Active {
			a.scheduleJob(j)
		}
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.srv.ListenAndServe(); err != http.ErrServerClosed {
			a.log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	return nil
}

func (a *App) stopServices() {
	a.scheduler.Stop()
	a.runner.Stop()

	if err := a.srv.Shutdown(context.Background()); err != nil {
		a.log.Printf("HTTP server Shutdown: %v", err)
	}
}

func (a *App) wait() {
	a.wg.Wait()

	a.db.Close()
}

func (a *App) scheduleJob(job domain.Job) {
	a.scheduler.AddEntry(cron.EntryID(job.ID), job.Cron, func(id cron.EntryID) {
		a.runner.AddJob(job)
	})
}

func (a *App) registerPerformers() {
	for _, p := range plugin.BuiltinPerformers() {
		a.log.Printf("Registering builtin performer %s\n", p.Info().Name)
		domain.RegisterPerformer(p)
	}

	plugin.PluginPerformers(filepath.Join(a.appPath, "plugins"), func(performer domain.Performer, err error) {
		a.log.Printf("Registering plugin performer %s\n", performer.Info().Name)
		if err != nil {
			a.log.Println(err)
			return
		}
		domain.RegisterPerformer(performer)
	})
}
