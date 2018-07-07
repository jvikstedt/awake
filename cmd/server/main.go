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

	runner := runner.New(logger)

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

//	// Set up authentication information for mailer
//	auth := smtp.PlainAuth(
//		"",
//		conf.MailConfig.Username,
//		conf.MailConfig.Password,
//		conf.MailConfig.Host,
//	)

// 		t := domain.New(logger, job.StepConfigs)
// 		steps := t.Run()
//
// 		data, _ := json.MarshalIndent(steps, "", "  ")
// 		logger.Printf("%s\n", data)
//
// 	Loop:
// 		for _, s := range steps {
// 			if s.Err != nil {
// 				if job.MailerEnabled {
// 					mail(logger, auth, conf.MailConfig, fmt.Sprintf("Something went wrong with job %d", job.ID), data)
// 				}
// 				break Loop
// 			}
// 		}
// 	})

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

type mailConfig struct {
	Identity string `json:"identity"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	To       string `json:"to"`
	From     string `json:"from"`
}

type config struct {
	Jobs       []domain.Job `json:"jobs"`
	MailConfig mailConfig   `json:"mailConfig"`
}

func loadConfig(logger *log.Logger, appPath string) (config, error) {
	data, err := ioutil.ReadFile(filepath.Join(appPath, "config.json"))
	if err != nil {
		return config{}, err
	}

	conf := config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return config{}, err
	}

	return conf, nil
}

// func mail(logger *log.Logger, auth smtp.Auth, conf mailConfig, subject string, body []byte) {
// 	msg := "From: " + conf.From + "\n" +
// 		"To: " + conf.To + "\n" +
// 		"Subject: " + subject + "\n\n" +
// 		string(body)
//
// 	err := smtp.SendMail(
// 		fmt.Sprintf("%s:%s", conf.Host, conf.Port),
// 		auth,
// 		conf.From,
// 		[]string{conf.To},
// 		[]byte(msg),
// 	)
// 	if err != nil {
// 		logger.Println(err)
// 	}
// }
