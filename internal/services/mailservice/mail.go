package mailservice

import (
	"bytes"
	"fmt"
	"html/template"
	"project-skbackend/configs"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/producerservice"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/uttemplate"

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
		cfg   *configs.Config
		ruser userrepo.IUserRepository
		sprod producerservice.IProducerService
	}

	IMailService interface {
		SendEmail(req requests.SendEmail) error
		SendResetPasswordEmail(data requests.SendEmailResetPassword) error
		SendVerifyEmail(req requests.SendEmailVerification) error
	}
)

func NewMailService(
	cfg *configs.Config,
	ruser userrepo.IUserRepository,
	sprod producerservice.IProducerService,
) *MailService {
	return &MailService{
		cfg:   cfg,
		ruser: ruser,
		sprod: sprod,
	}
}

func parseTemplate(templateFileName string, data any) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		utlogger.Error(err)
		return "", consttypes.ErrFailedToParseFile
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		utlogger.Error(err)
		return "", consttypes.ErrFailedToWriteFile
	}

	return buf.String(), nil
}
func (s *MailService) SendEmail(
	req requests.SendEmail,
) error {
	var (
		body bytes.Buffer
	)

	templates, err := uttemplate.ParseTemplateDir("templates", req.Template)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrFailedToParseFile
	}

	templates = templates.Lookup(req.Template)

	err = templates.Execute(&body, &req.Data)
	if err != nil {
		utlogger.Error(err)
		return consttypes.ErrFailedToWriteFile
	}

	m := gomail.NewMessage()

	m.SetHeaders(map[string][]string{
		"From":    {s.cfg.Mail.From},
		"To":      {req.Email},
		"Subject": {req.Subject},
	})
	m.SetBody("text/html", body.String())

	// TODO - need to enable this
	// d := gomail.NewDialer(MAIL_HOST, MAIL_PORT, MAIL_FROM, MAIL_PASS)
	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	//
	// if err := d.DialAndSend(m); err != nil {
	// 	return err
	// }

	return nil
}

func (s *MailService) SendResetPasswordEmail(req requests.SendEmailResetPassword) error {
	sereq := requests.SendEmail{
		Template: "reset_password.html",
		Subject:  "Reset Password Request on Meals to Heals",
		Email:    req.Email,
		Data: map[string]any{
			"Name":    req.Name,
			"Email":   req.Email,
			"LinkUrl": template.URL(req.LinkUrl),
		},
	}

	if err := s.sprod.PublishEmail(sereq); err != nil {
		return consttypes.ErrFailedToPublishMessage
	}

	return nil
}

func (s *MailService) SendVerifyEmail(req requests.SendEmailVerification) error {
	sereq := requests.SendEmail{
		Template: "verify_email.html",
		Subject:  fmt.Sprintf("Verify Your Email Address on Meals to Heals"),
		Email:    req.Email,
		Data: map[string]any{
			"Email": req.Email,
			"Token": req.Token,
			"Name":  req.Name,
		},
	}

	if err := s.sprod.PublishEmail(sereq); err != nil {
		return consttypes.ErrFailedToPublishMessage
	}

	return nil
}
