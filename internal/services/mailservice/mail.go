package mailservice

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"text/template"
)

type (
	MailService struct {
		cfg     *configs.Config
		mailset *MailSetup
	}

	MailSetup struct {
		from    string
		pass    string
		host    string
		port    string
		address string
		temdir  string
	}

	IMailService interface {
		SendVerificationEmail(req requests.SendEmail) error
	}
)

func NewMailService(
	cfg *configs.Config,
) *MailService {
	mailset := MailSetup{
		from:    cfg.Mail.From,
		pass:    cfg.Mail.Password,
		host:    cfg.Mail.SMTPHost,
		port:    cfg.Mail.SMTPPort,
		address: cfg.Mail.SMTPHost + ":" + cfg.Mail.SMTPPort,
		temdir:  cfg.Mail.TemplateDirectory,
	}

	return &MailService{
		cfg:     cfg,
		mailset: &mailset,
	}
}

func (m *MailService) SendVerificationEmail(req requests.SendEmail) error {
	to := []string{req.Email}
	temfile := req.Template

	/* ------------------------ setup plain auth setting ------------------------ */
	auth := smtp.PlainAuth("", m.mailset.from, m.mailset.pass, m.mailset.host)

	/* -------------------------- setup parse template -------------------------- */
	t, err := template.ParseFiles(m.mailset.temdir + temfile)
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
		Token   string
	}{
		Name:    req.Name,
		Message: "This is a message",
		Token:   req.Token,
	})

	/* ------------------------------- send email ------------------------------- */
	err = smtp.SendMail(m.mailset.address, auth, m.mailset.from, to, body.Bytes())
	if err != nil {
		return err
	}

	log.Println("Email has been sent!")

	return nil
}
