package mailservice

import (
	"bytes"
	"log"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/packages/utils/utlogger"
	"text/template"

	"gopkg.in/gomail.v2"
)

var (
	cfg = configs.GetInstance()

	MAIL_NAME   = cfg.Mail.Name
	MAIL_FROM   = cfg.Mail.From
	MAIL_PASS   = cfg.Mail.Password
	MAIL_HOST   = cfg.Mail.SMTPHost
	MAIL_PORT   = cfg.Mail.SMTPPort
	MAIL_TEMDIR = cfg.Mail.TemplateDir
)

type (
	MailService struct {
		cfg *configs.Config
	}

	IMailService interface {
		SendMail(req requests.SendMail, data any) error
		SendResetPassword(reqdef requests.SendMail, reqdata requests.ResetPasswordEmail) error
	}
)

func NewMailService(
	cfg *configs.Config,
) *MailService {
	return &MailService{
		cfg: cfg,
	}
}

func parseTemplate(templateFileName string, data any) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		utlogger.LogError(err)
		return "", err
	}

	return buf.String(), nil
}

func (s *MailService) SendMail(req requests.SendMail, data any) error {
	// * parsing the template from template dir and its data
	result, err := parseTemplate(MAIL_TEMDIR+req.Template, data)
	if err != nil {
		return err
	}

	// * create a new mailer config with existing parameter
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", MAIL_NAME)
	mailer.SetHeader("To", req.To...)
	mailer.SetHeader("Subject", req.Subject)
	mailer.SetBody("text/html", result)

	// * create new dialer from existing config
	dialer := gomail.NewDialer(
		MAIL_HOST,
		MAIL_PORT,
		MAIL_FROM,
		MAIL_PASS,
	)

	// * send the mail using the new dialer
	err = dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func (s *MailService) SendResetPassword(reqdef requests.SendMail, reqdata requests.ResetPasswordEmail) error {
	err := s.SendMail(reqdef, reqdata)
	if err != nil {
		return err
	}

	log.Println("Email has been sent!")

	return nil
}
