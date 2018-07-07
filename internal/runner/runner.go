package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"sync"

	"github.com/jvikstedt/awake/internal/domain"
)

type Runner struct {
	jobs chan domain.Job
	quit chan struct{}
	log  *log.Logger
	conf domain.Config
}

func New(logger *log.Logger, conf domain.Config) *Runner {
	return &Runner{
		jobs: make(chan domain.Job, 100),
		quit: make(chan struct{}),
		log:  logger,
		conf: conf,
	}
}

func (r *Runner) AddJob(j domain.Job) {
	r.jobs <- j
}

func (r *Runner) Start() {
	r.log.Println("Started runner")
	var wg sync.WaitGroup
Loop:
	for {
		select {
		case j := <-r.jobs:
			wg.Add(1)
			go func() {
				defer wg.Done()
				r.handleJob(j)
			}()
		case <-r.quit:
			break Loop
		}
	}
	r.log.Println("Runner waiting for jobs to finish...")
	wg.Wait()
	r.log.Println("Stopped runner")
}

func (r *Runner) Stop() {
	r.log.Println("Stopping runner...")
	r.quit <- struct{}{}
}

func (r *Runner) handleJob(job domain.Job) {
	t := domain.New(r.log, job.StepConfigs)
	steps := t.Run()

	data, _ := json.MarshalIndent(steps, "", "  ")
	r.log.Printf("%s\n", data)

	auth := smtp.PlainAuth(
		"",
		r.conf.MailConfig.Username,
		r.conf.MailConfig.Password,
		r.conf.MailConfig.Host,
	)

Loop:
	for _, s := range steps {
		if s.Err != nil {
			if job.MailerEnabled {
				mail(r.log, auth, r.conf.MailConfig, fmt.Sprintf("Something went wrong with job %d", job.ID), data)
			}
			break Loop
		}
	}
}

func mail(logger *log.Logger, auth smtp.Auth, conf domain.MailConfig, subject string, body []byte) {
	msg := "From: " + conf.From + "\n" +
		"To: " + conf.To + "\n" +
		"Subject: " + subject + "\n\n" +
		string(body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		auth,
		conf.From,
		[]string{conf.To},
		[]byte(msg),
	)
	if err != nil {
		logger.Println(err)
	}
}
