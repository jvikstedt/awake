package main

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/jvikstedt/awake/cron"
	"github.com/jvikstedt/awake/internal/domain"
	"github.com/jvikstedt/awake/internal/plugin"
	"github.com/jvikstedt/awake/internal/runner"
)

type App struct {
	log       *log.Logger
	port      string
	wg        sync.WaitGroup
	config    domain.Config
	appPath   string
	scheduler *cron.Scheduler
	runner    *runner.Runner
}

func newApp(logger *log.Logger, port string, config domain.Config, appPath string) *App {
	return &App{
		log:       logger,
		port:      port,
		config:    config,
		appPath:   appPath,
		scheduler: cron.New(logger),
		runner:    runner.New(logger, config),
	}
}

func (a *App) startServices() {
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

	for _, j := range a.config.Jobs {
		a.scheduleJob(j)
	}
}

func (a *App) stopServices() {
	a.scheduler.Stop()
	a.runner.Stop()
}

func (a *App) wait() {
	a.wg.Wait()
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
