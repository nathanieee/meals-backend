package requests

import (
	"encoding/json"
	"project-skbackend/packages/utils/utlogger"
)

type (
	SendEmailRequest struct {
		Template string `binding:"required,oneof=email_veification.html"`
		Subject  string `binding:"required"`
		Name     string
		Email    string `binding:"required,email"`
		Token    int
	}
)

func (s SendEmailRequest) ToString() string {
	json, err := json.Marshal(s)
	if err != nil {
		utlogger.LogError(err)
		return ""
	}
	return string(json)
}
