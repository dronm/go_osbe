package validEmail

import "net/mail"

func Check(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}
