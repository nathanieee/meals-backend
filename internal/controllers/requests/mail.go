package requests

import (
	"encoding/json"
	"fmt"
)

type (
	SendEmailRequest struct {
		Template string `validate:"required,oneof=email_veification.html"`
		Subject  string `validate:"required"`
		Name     string
		Email    string `validate:"required,email"`
		Token    int
	}
)

func (s SendEmailRequest) ToString() string {
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ""
	}
	return string(b)
}
