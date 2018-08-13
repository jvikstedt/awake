package runner

import (
	"fmt"
	"log"
	"sync"

	"github.com/jvikstedt/awake"
	"github.com/jvikstedt/awake/internal/domain"
	"github.com/jvikstedt/awake/internal/plugin"
)

type Runner struct {
	jobs             chan domain.Job
	quit             chan struct{}
	log              *log.Logger
	conf             domain.Config
	resultRepository domain.ResultRepository
}

func New(logger *log.Logger, conf domain.Config, resultRepository domain.ResultRepository) *Runner {
	return &Runner{
		jobs:             make(chan domain.Job, 100),
		quit:             make(chan struct{}),
		log:              logger,
		conf:             conf,
		resultRepository: resultRepository,
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
	r.log.Printf("Running job %d\n", job.ID)
	scope := newScope(r.log, r.conf.PerformerConfigs, *job.StepConfigs)

	for i, v := range scope.steps {
		scope.current = i

		performer, ok := plugin.FindPerformer(v.Conf.Tag)
		if !ok {
			scope.addAlert(awake.Alert{Type: awake.AlertError, Value: fmt.Sprintf("argh... performer not found %s", v.Conf.Tag)})
			continue
		}

		if err := performer.Perform(scope); err != nil {
			scope.addAlert(awake.Alert{Type: awake.AlertError, Value: err.Error()})
		}
	}

	result := domain.Result{
		JobID: job.ID,
		Steps: &scope.steps,
	}

	_, err := r.resultRepository.Create(result)
	if err != nil {
		r.log.Println(err)
	}

	r.log.Printf("Job %d execution done\n", job.ID)
}
