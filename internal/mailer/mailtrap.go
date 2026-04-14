package mailer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

const resendAPIURL = "https://api.resend.com/emails"

type mailTrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (mailTrapClient, error) {
	if apiKey == "" {
		return mailTrapClient{}, errors.New("api key is required")
	}

	return mailTrapClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m mailTrapClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return -1, err
	}

	payload := map[string]any{
		"from":    m.fromEmail,
		"to":      []string{email},
		"subject": subject.String(),
		"html":    body.String(),
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest(http.MethodPost, resendAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return -1, err
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		return resp.StatusCode, fmt.Errorf("mailtrap API error: status %d, body: %s", resp.StatusCode, respBody.String())
	}

	return resp.StatusCode, nil
}
