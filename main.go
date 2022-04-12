package main

import (
	"strings"
	"time"

	"encoding/json"
	"fmt"
	"net/http"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "okay"}`))
	})
	r.Post("/v1/verification", handelGetEmailVerification)

	http.ListenAndServe(":3000", r)
}

type EmailRequest struct {
	Email string `json:"email"`
}

type EmailResponse struct {
	Valid       bool          `json:"valid"`
	Suggestions string        `json:"suggestion,omitempty"`
	Reason      string        `json:"reason,omitempty"`
	Elapsed     time.Duration `json:"elapsed"`
}

func handelGetEmailVerification(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	decoder := json.NewDecoder(r.Body)

	var er EmailRequest

	err := decoder.Decode(&er)
	if err != nil {
		panic(err)
	}

	emailResponse := verifyEmail(er)

	time_now := time.Now()
	emailResponse.Elapsed = time.Duration(time_now.Sub(start).Milliseconds())
	bytes, err := json.Marshal(emailResponse)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, _ = fmt.Fprint(w, string(bytes))
}

func verifyEmail(er EmailRequest) EmailResponse {
	verifier := emailverifier.NewVerifier()
	var emailResponse EmailResponse
	emailResponse.Valid = true

	ret, err := verifier.Verify(er.Email)
	if err != nil {
		emailResponse.Reason = "Failed to verify"
		emailResponse.Valid = false
	}
	if !ret.Syntax.Valid {
		emailResponse.Reason = "Email syntax is incorrect"
		emailResponse.Valid = false
	}

	if emailResponse.Valid == true {
		domain := strings.Split(er.Email, "@")[1]
		suggestion := verifier.SuggestDomain(domain)

		if suggestion != "" {
			emailResponse.Valid = false
			emailResponse.Reason = "Did you mean " + suggestion + " instead of " + domain
			emailResponse.Suggestions = suggestion
		}
	}
	return emailResponse
}
