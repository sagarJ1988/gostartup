package validation

import (
	"regexp"
	"strings"
)

type Message struct {
	Email    string
	Password string
	Errors   map[string]string
}

type SignupForm struct {
	FirstName string
	LastName  string
	UserName  string
	Password  string
	Errors    map[string]string
}

func (msg *Message) ValidateLogin() bool {
	msg.Errors = make(map[string]string)

	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(msg.Email))
	if matched == false {
		msg.Errors["Email"] = "Please enter a valid email address"
	}

	if strings.TrimSpace(msg.Password) == "" {
		msg.Errors["Password"] = "Password required"
	}

	return len(msg.Errors) == 0
}

func (msg *SignupForm) ValidateSignup() bool {
	msg.Errors = make(map[string]string)

	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(msg.UserName))

	if strings.TrimSpace(msg.FirstName) == "" {
		msg.Errors["FirstName"] = "First Name Required"
	}

	if strings.TrimSpace(msg.LastName) == "" {
		msg.Errors["LastName"] = "Last Name Required"
	}

	if matched == false {
		msg.Errors["UserName"] = "Please Enter A Valid Email Address"
	}

	if strings.TrimSpace(msg.Password) == "" {
		msg.Errors["Password"] = "Password Required"
	}
	return len(msg.Errors) == 0
}
