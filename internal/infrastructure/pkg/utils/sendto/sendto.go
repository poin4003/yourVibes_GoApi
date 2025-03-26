package sendto

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func GetMailServiceSettings() (string, string, string, string) {
	return global.Config.MailService.SMTPHost,
		global.Config.MailService.SMTPPort,
		global.Config.MailService.SMTPUsername,
		global.Config.MailService.SMTPPassword
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From.Name)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ","))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func SendTemplateEmail(
	to []string,
	from string,
	nameTemplate string,
	dataTemplate map[string]interface{},
	subject string,
) error {
	htmlBody, err := getMailTemplate(nameTemplate, dataTemplate)
	if err != nil {
		return err
	}

	return send(to, from, htmlBody, subject)
}

func getMailTemplate(nameTemplate string, dataTemplate map[string]interface{}) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(nameTemplate).ParseFiles("templates/" + nameTemplate))
	err := t.Execute(htmlTemplate, dataTemplate)
	if err != nil {
		return "", err
	}
	return htmlTemplate.String(), nil
}

func send(to []string, from string, htmlTemplate string, subject string) error {
	contentEmail := Mail{
		From:    EmailAddress{Address: from, Name: "YourVibes"},
		To:      to,
		Subject: subject,
		Body:    htmlTemplate,
	}

	messageMail := BuildMessage(contentEmail)

	SMTPHost, SMTPPort, SMTPUsername, SMTPPassword := GetMailServiceSettings()

	// Send smtp
	authentication := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)

	err := smtp.SendMail(SMTPHost+":"+SMTPPort, authentication, from, to, []byte(messageMail))

	if err != nil {
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}

	return nil
}
