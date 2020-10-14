package alert

import (
	"os"
	"testing"
)

// SENDAPI==`kiya shared get sendgrid/apikey/moneypenny` FROM=... TO=... go test
func TestSendEmail(t *testing.T) {
	apikey := os.Getenv("SENDAPI")
	if len(apikey) == 0 {
		t.Skip()
	}
	if err := SendEmail("TestSendEmail", os.Getenv("FROM"), os.Getenv("TO"), "email.json", "email_template.txt", apikey); err != nil {
		t.Error(err)
	}
}
