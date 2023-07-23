package handlers

import (
	"gomarket/internal/errors"
	"regexp"
	"strconv"
)

const (
	MinNumLogin = 4
	MaxNumLogin = 20
	MinNumPass  = 6
	MaxNumPass  = 40
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

func ValidateLoadOrder(orderID string) *errors.ErrorApp {
	if isValidOrderID(orderID) {
		return errors.NewErrFailedValidation("invalid order_id")
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

func isValidOrderID(orderID string) bool {
	if _, err := strconv.Atoi(orderID); err != nil {
		return true
	}

	sum := 0
	alternate := false
	for i := len(orderID) - 1; i >= 0; i-- {
		digit := int(orderID[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}

	return !(sum%10 == 0)
}
