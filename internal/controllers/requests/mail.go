package requests

import (
	"encoding/json"
	"project-skbackend/packages/utils/utlogger"
)

type (
	SendMail struct {
		Template string   `binding:"required"`
		To       []string `binding:"required"`
		Subject  string   `binding:"required"`
	}

	ResetPasswordEmail struct {
		Token string
	}
)

func (s ResetPasswordEmail) ToString() string {
	json, err := json.Marshal(s)
	if err != nil {
		utlogger.LogError(err)
		return ""
	}
	return string(json)
}
