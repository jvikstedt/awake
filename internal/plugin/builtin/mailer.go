package builtin

import (
	"fmt"
	"net/smtp"

	"github.com/jvikstedt/awake"
)

type Mailer struct{}

func (Mailer) Info() awake.PerformerInfo {
	return awake.PerformerInfo{
		Name:        "builtin_mailer",
		DisplayName: "Mailer",
	}
}

func (Mailer) Perform(scope awake.Scope) error {
	// identity, _ := scope.ValueAsString("identity")
	username, _ := scope.ValueAsString("username")
	password, _ := scope.ValueAsString("password")
	host, _ := scope.ValueAsString("host")
	port, _ := scope.ValueAsString("port")
	to, _ := scope.ValueAsString("to")
	from, _ := scope.ValueAsString("from")
	subject, _ := scope.ValueAsString("subject")
	body, _ := scope.ValueAsString("body")

	shouldRun, ok := scope.ValueAsBool("shouldRun")
	if !ok {
		shouldRun = true
	}

	if shouldRun {
		auth := smtp.PlainAuth(
			"",
			username,
			password,
			host,
		)

		return mail(auth, from, to, host, port, subject, body)
	}
	return nil
}

func mail(auth smtp.Auth, from string, to string, host string, port string, subject string, body string) error {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		string(body)

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", host, port),
		auth,
		from,
		[]string{to},
		[]byte(msg),
	)
}
