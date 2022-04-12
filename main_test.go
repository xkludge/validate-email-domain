package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type suiteVerifyEmail struct {
	suite.Suite
}

func (s *suiteVerifyEmail) TestGetStatusHandler() {
	_, err := http.NewRequest("GET", "/status", nil)

	assert.NoError(s.T(), err, "error")
}

func (s *suiteVerifyEmail) TestVerifyEmail() {
	var er EmailRequest
	er.Email = "wasd@gmail.com"

	emailResponse := verifyEmail(er)

	assert.EqualValues(s.T(), emailResponse.Valid, true)
}

func (s *suiteVerifyEmail) TestVerifyEmailFailsWithSuggestion() {
	var er EmailRequest
	er.Email = "wasd@gail.com"

	emailResponse := verifyEmail(er)

	assert.EqualValues(s.T(), emailResponse.Valid, false)
	assert.EqualValues(s.T(), emailResponse.Reason, "Did you mean gmail.com instead of gail.com")
}

func (s *suiteVerifyEmail) TestVerifyEmailInvalidSyntax() {
	var er EmailRequest
	er.Email = "wasd@@gail.com"

	emailResponse := verifyEmail(er)

	assert.EqualValues(s.T(), emailResponse.Valid, false)
	assert.EqualValues(s.T(), emailResponse.Reason, "Email syntax is incorrect")
}

func (s *suiteVerifyEmail) TestVerifyEmailNoEmailSupplied() {
	var er EmailRequest
	er.Email = ""

	emailResponse := verifyEmail(er)

	assert.EqualValues(s.T(), emailResponse.Valid, false)
	assert.EqualValues(s.T(), emailResponse.Reason, "Email syntax is incorrect")
}
func TestSuiteRest(t *testing.T) {
	suite.Run(t, new(suiteVerifyEmail))
}
