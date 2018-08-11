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
	scope := newScope(r.log, r.conf.PerformerConfigs, *job.StepConfigs)

	for i, v := range scope.steps {
		scope.current = i

		performer, ok := domain.FindPerformer(v.Conf.Tag)
		if !ok {
			v.Err = fmt.Errorf("argh... performer not found %s", v.Conf.Tag)
			v.ErrMsg = v.Err.Error()
			continue
		}

		if v.Err = performer.Perform(scope); v.Err != nil {
			v.ErrMsg = v.Err.Error()
		}
	}

	r.log.Printf("Job %d execution done\n", job.ID)
}
