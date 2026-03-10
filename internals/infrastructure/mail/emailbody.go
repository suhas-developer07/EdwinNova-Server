package mail

import (
	"bytes"
	"text/template"
	"time"
)

func BuildRegistrationEmailBody(teamName, pmName, pmEmail, pmContact string, applicationID string, createdAt time.Time) (string, error) {

	loc, _ := time.LoadLocation("Asia/Kolkata")
	indianTime := createdAt.In(loc).Format("02 Jan 2006, 03:04 PM IST")

	data := map[string]interface{}{
		"TeamName":      teamName,
		"PMName":        pmName,
		"PMEmail":       pmEmail,
		"PMContact":     pmContact,
		"ApplicationID": applicationID,
		"IndianTime":    indianTime,
	}

	tmpl, err := template.New("registration").Parse(RegistrationTemplate)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}