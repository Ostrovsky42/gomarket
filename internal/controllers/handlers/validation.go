package handlers

import (
	"gomarket/internal/errors"
	"regexp"
	"strconv"
)

const (
	minNumLogin = 4
	maxNumLogin = 20
	minNumPass  = 6
	maxNumPass  = 40
)

var loginMask = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func ValidationAuth(login string, pass string) *errors.ErrorApp {
	if err := validateLogin(login); err != nil {
		return err
	}

	return validatePassword(pass)
}

func ValidateLoadOrder(orderID string) *errors.ErrorApp {
	if isValidOrderID(orderID) {
		return errors.NewErrFailedValidation("invalid order_id")
	}

	return nil
}

func ValidateUsePoints(orderID string, points float64) *errors.ErrorApp {
	if isValidOrderID(orderID) {
		return errors.NewErrFailedValidation("invalid order_id")
	}

	return validatePoints(points)
}

func validatePoints(points float64) *errors.ErrorApp {
	if points < 0 {
		return errors.NewErrFailedValidation("less zero")
	}

	return nil
}

func validateLogin(login string) *errors.ErrorApp {
	if !loginMask.MatchString(login) {
		return errors.NewErrFailedValidation("invalid login")
	}

	if len(login) < minNumLogin || len(login) > maxNumLogin {
		return errors.NewErrFailedValidation("invalid login")
	}

	return nil
}

func validatePassword(password string) *errors.ErrorApp {
	if len(password) < minNumPass || len(password) > maxNumPass {
		return errors.NewErrFailedValidation("invalid password")
	}

	return nil
}

// nolint:gomnd
func isValidOrderID(orderID string) bool {
	if _, err := strconv.Atoi(orderID); err != nil {
		return true
	}

	sum := 0
	alternate := false

	i := len(orderID) - 1
	for ; i >= 0; i-- {
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
