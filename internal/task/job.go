package task

type Job struct {
	ID            int          `json:"id"`
	Cron          string       `json:"cron"`
	StepConfigs   []StepConfig `json:"stepConfigs"`
	MailerEnabled bool         `json:"mailerEnabled"`
}
