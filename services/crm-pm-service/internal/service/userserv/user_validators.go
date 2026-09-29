package userserv

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"

	"github.com/erickgreco/startup-backend/services/crm-pm-service/internal/domain"
	"github.com/nyaruka/phonenumbers/v2"
)

var supportedLanguages = map[string]bool{
	"es": true,
	"en": true,
}

// ! Uso unicamente en user_validators.go
func validateNameOrLastName(s string) (string, bool) {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return "", false
	}

	for _, f := range fields {
		for _, r := range f {
			if !unicode.IsLetter(r) && r != '-' && r != '\'' {
				return "", false
			}
		}
	}

	return strings.Join(fields, " "), true
}

func validateFullName(name, lastName string) (string, string, error) {
	n, nameOK := validateNameOrLastName(name)
	ln, lastNameOk := validateNameOrLastName(lastName)

	var nameErr, lastNameErr error
	if !nameOK {
		nameErr = domain.ErrInvalidName
	}
	if !lastNameOk {
		lastNameErr = domain.ErrInvalidLastName
	}

	if err := errors.Join(nameErr, lastNameErr); err != nil {
		return "", "", err
	}

	return n, ln, nil
}

func validateEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email {
		return "", domain.ErrInvalidEmail
	}

	at := strings.LastIndex(parsed.Address, "@")
	if !hasValidDomainShape(parsed.Address[at+1:]) {
		return "", domain.ErrInvalidEmail
	}

	return parsed.Address, nil
}

// ! Uso unicamente en user_validators.go
func hasValidDomainShape(d string) bool {
	if strings.HasPrefix(d, "[") {
		return false
	}
	if strings.Contains(d, "..") {
		return false
	}

	dot := strings.LastIndex(d, ".")

	return dot > 0 && dot < len(d)-1
}

func validatePhoneNumber(p string) (string, error) {
	num, err := phonenumbers.Parse(p, "")
	if err != nil || !phonenumbers.IsValidNumber(num) {
		return "", domain.ErrInvalidPhoneNumber
	}

	return phonenumbers.Format(num, phonenumbers.E164), nil
}

func isValidLanguage(lang string) bool {
	return supportedLanguages[lang]
}
