package task

type Job struct {
	ID          int
	Cron        string
	StepConfigs []StepConfig
}
