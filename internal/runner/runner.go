package runner

import (
	"encoding/json"
	"log"
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
	t := domain.NewTask(r.log, r.conf.PerformerConfigs, job.StepConfigs)
	steps := t.Run()

	data, _ := json.MarshalIndent(steps, "", "  ")
	r.log.Printf("%s\n", data)
}
