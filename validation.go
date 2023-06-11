package main

import (
	"errors"
	"fmt"
	"strings"

	proto "github.com/nutreet/common/gen/user"
)

func ValidateRegisterRequest(data *proto.RegisterRequest) error {
	if err := validateEmail(data.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	return nil
}

func validateEmail(email string) error {
	if len(email) == 0 {
		return errors.New("email is required")
	} else if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("invalid email format")
	}

	return nil
}
