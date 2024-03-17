package requests

type (
	SendEmail struct {
		Template string         `validate:"required"`
		Subject  string         `validate:"required"`
		Email    string         `validate:"required"`
		Data     map[string]any `validate:"required"`
	}

	SendEmailResetPassword struct {
		Name    string `validate:"required"`
		Email   string `validate:"required,email"`
		LinkUrl string `validate:"required"`
	}

	SendEmailInvitation struct {
		EmailInviter     string `validate:"required,email"`
		EmailInvitee     string `validate:"required,email"`
		OrganizationName string `validate:"required"`
		DepartmentName   string `validate:"required"`
		LinkUrl          string `validate:"required"`
	}

	SendEmailSurveyResult struct {
		Name    string `validate:"required"`
		Email   string `validate:"required,email"`
		LinkUrl string `validate:"required"`
	}

	SendEmailVerification struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
		Token int    `validate:"required"`
	}
)
