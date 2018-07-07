package domain

type MailConfig struct {
	Identity string `json:"identity"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	To       string `json:"to"`
	From     string `json:"from"`
}

type Config struct {
	Jobs       []Job `json:"jobs"`
	MailConfig `json:"mailConfig"`
}
