package mail

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"text/template"
)

type MailService struct {
	cfg *configs.Config
}

func NewMailService(cfg *configs.Config) *MailService {
	return &MailService{cfg: cfg}
}

func (m *MailService) SendVerificationEmail(req requests.SendEmailRequest) error {

	/* ------------------------- setup from and password ------------------------ */

	from := m.cfg.Mail.From
	pass := m.cfg.Mail.Password
	to := []string{req.Email}

	/* ------------------- setup smtp host, port, and address ------------------- */

	smtphost := m.cfg.Mail.SMTPHost
	smtpport := m.cfg.Mail.SMTPPort
	address := smtphost + ":" + smtpport

	/* --------------- setup mail template directory and file name -------------- */

	temdir := m.cfg.Mail.TemplateDirectory
	temfile := req.Template

	/* ------------------------ setup plain auth setting ------------------------ */

	auth := smtp.PlainAuth("", from, pass, smtphost)

	/* -------------------------- setup parse template -------------------------- */

	t, err := template.ParseFiles(temdir + temfile)
	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	/* -------------------- setup the content of the template ------------------- */

	// TODO CHANGE THE CONTENT ACCORDINGLY
	t.Execute(&body, struct {
		Name    string
		Message string
		Token   int
	}{
		Name:    req.Name,
		Message: "This is a message",
		Token:   req.Token,
	})

	/* ------------------------------- send email ------------------------------- */

	err = smtp.SendMail(address, auth, from, to, body.Bytes())
	if err != nil {
		return err
	}

	log.Println("Email has been sent!")

	return nil
}
