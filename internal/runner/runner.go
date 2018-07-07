package runner

import (
	"fmt"
	"log"
	"sync"

	"github.com/jvikstedt/awake/internal/domain"
)

type Runner struct {
	jobs chan domain.Job
	quit chan struct{}
	log  *log.Logger
}

func New(logger *log.Logger) *Runner {
	return &Runner{
		jobs: make(chan domain.Job, 100),
		quit: make(chan struct{}),
		log:  logger,
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
				// DO stuff
				fmt.Println(j)
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
