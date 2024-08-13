package mailservice

import (
	"bytes"
	"crypto/tls"
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

type (
	MailService struct {
		cfg   *configs.Config
		ruser userrepo.IUserRepository
		sprod producerservice.IProducerService

		mailname   string
		mailfrom   string
		mailpass   string
		mailhost   string
		mailport   int
		mailtemdir string

		logourl string
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

		mailname:   cfg.Mail.Name,
		mailfrom:   cfg.Mail.From,
		mailpass:   cfg.Mail.Password,
		mailhost:   cfg.Mail.SMTPHost,
		mailport:   cfg.Mail.SMTPPort,
		mailtemdir: cfg.Mail.TemplateDir,

		logourl: cfg.LogoURL,
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

	templates, err := uttemplate.ParseTemplateDir(s.mailtemdir, req.Template)
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
		"To":      {"jonathanvnc@gmail.com", "shavyerachristine@gmail.com"}, // TODO: change this into the correct email
		"Subject": {req.Subject},
	})
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer(s.mailhost, s.mailport, s.mailfrom, s.mailpass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *MailService) SendResetPasswordEmail(req requests.SendEmailResetPassword) error {
	sereq := requests.SendEmail{
		Template: "reset_password.html",
		Subject:  "Reset Password Request on Meals to Heals",
		Email:    req.Email,
		Data: map[string]any{
			"LogoUrl": s.logourl,
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
			"LogoUrl": s.logourl,
			"Email":   req.Email,
			"Token":   req.Token,
			"Name":    req.Name,
		},
	}

	if err := s.sprod.PublishEmail(sereq); err != nil {
		return consttypes.ErrFailedToPublishMessage
	}

	return nil
}
