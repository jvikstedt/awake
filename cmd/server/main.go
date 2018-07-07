package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

	"github.com/jvikstedt/awake/cron"
	"github.com/jvikstedt/awake/internal/domain"
	"github.com/jvikstedt/awake/internal/plugin"
	"github.com/jvikstedt/awake/internal/runner"
)

func main() {
	appPath := getApplicationPath()

	// Setup logger
	f, err := os.OpenFile(filepath.Join(appPath, "awake.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	logger := log.New(f, "", log.LstdFlags)

	registerPerformers(logger, appPath)

	conf, err := loadConfig(logger, appPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	runner := runner.New(logger, conf)

	wg.Add(1)
	go func() {
		defer wg.Done()
		runner.Start()
	}()

	scheduler := cron.New(logger)

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start()
	}()

	for _, j := range conf.Jobs {
		scheduleJob(scheduler, runner, j)
	}

	// Handle signals
	go func() {
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, os.Kill)

		<-sigquit

		log.Println("Stopping everything...")
		scheduler.Stop()
		runner.Stop()
	}()

	wg.Wait()
}

func scheduleJob(scheduler *cron.Scheduler, runner *runner.Runner, job domain.Job) {
	scheduler.AddEntry(cron.EntryID(job.ID), job.Cron, func(id cron.EntryID) {
		runner.AddJob(job)
	})
}

func registerPerformers(logger *log.Logger, appPath string) {
	for _, p := range plugin.BuiltinPerformers() {
		logger.Printf("Registering builtin performer %s\n", p.Info().Name)
		domain.RegisterPerformer(p)
	}

	plugin.PluginPerformers(filepath.Join(appPath, "plugins"), func(performer domain.Performer, err error) {
		logger.Printf("Registering plugin performer %s\n", performer.Info().Name)
		if err != nil {
			logger.Println(err)
			return
		}
		domain.RegisterPerformer(performer)
	})
}

func loadConfig(logger *log.Logger, appPath string) (domain.Config, error) {
	data, err := ioutil.ReadFile(filepath.Join(appPath, "config.json"))
	if err != nil {
		return domain.Config{}, err
	}

	conf := domain.Config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return domain.Config{}, err
	}

	return conf, nil
}
