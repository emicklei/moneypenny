package alert

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/emicklei/moneypenny/util"
	"github.com/emicklei/tre"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends HTML content to an email address(es) unless the data (jsonFilename) is missing.
func SendEmail(subject, fromAddress, toAddress string, jsonFilename, htmlTemplateFilename, apikey string) error {
	util.CheckNonEmpty("from email", fromAddress)
	util.CheckNonEmpty("to email(s)", toAddress)
	util.CheckNonEmpty("subject email", subject)
	util.CheckNonEmpty("api-key", apikey)

	dataJSON, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		log.Println("Warning: no JSON data file found, skip sending email", jsonFilename)
		return nil
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return tre.New(err, "parsing JSON")
	}
	// fallback for non-html will have the JSON formatted data
	plainTextContent := string(dataJSON)

	// use template + data to get html content
	buf := new(bytes.Buffer)
	templateData, err := ioutil.ReadFile(htmlTemplateFilename)
	if err != nil {
		return tre.New(err, "reading template", "file", "htmlTemplateFilename")
	}
	t, err := template.New("SendEmail").Parse(string(templateData))
	if err != nil {
		return tre.New(err, "parsing template", "file", "htmlTemplateFilename")
	}
	err = t.ExecuteTemplate(buf, "SendEmail", data)
	if err != nil {
		return tre.New(err, "executing template", "file", "htmlTemplateFilename")
	}
	htmlContent := buf.String()

	// compose email
	m := mail.NewV3Mail()
	m.SetFrom(mail.NewEmail("Moneypenny", fromAddress))
	m.AddContent(mail.NewContent("text/plain", plainTextContent))
	m.AddContent(mail.NewContent("text/html", htmlContent))
	personalization := mail.NewPersonalization()
	for _, each := range strings.Split(toAddress, ",") {
		eachAddress := strings.TrimSpace(each)
		personalization.AddTos(mail.NewEmail("Moneypenny User", eachAddress))
	}
	personalization.Subject = subject
	m.AddPersonalizations(personalization)

	// send it away
	client := sendgrid.NewSendClient(apikey)
	resp, err := client.Send(m)
	if resp.StatusCode > http.StatusAccepted {
		return tre.New(errors.New("failed to send email"), "sendgrid failed to deliver", "status", resp.StatusCode, "body", resp.Body)
	}
	return err
}
