package handlers

import (
	"regexp"

	"gomarket/internal/errors"
)

const (
	MinNumLogin = 4
	MaxNumLogin = 20
	MinNumPass  = 6
	MaxNumPass  = 20
)

var loginMask = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func ValidationAuth(login string, pass string) *errors.ErrorApp {
	if err := validateLogin(login); err != nil {
		return err
	}

	if err := validatePassword(pass); err != nil {
		return err
	}

	return nil
}

func validateLogin(login string) *errors.ErrorApp {
	if !loginMask.MatchString(login) {
		return errors.NewErrFailedValidation("invalid login")
	}

	if len(login) < MinNumLogin || len(login) > MaxNumLogin {
		return errors.NewErrFailedValidation("invalid login")
	}

	return nil
}

func validatePassword(password string) *errors.ErrorApp {
	if len(password) < MinNumPass || len(password) > MaxNumPass {
		return errors.NewErrFailedValidation("invalid password")
	}

	return nil
}
