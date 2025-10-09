package mailer

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ouiasy/golang-auth/conf"
	"github.com/ouiasy/golang-auth/crypto"
	"github.com/ouiasy/golang-auth/models"
	"github.com/resend/resend-go/v2"
)

var (
	ErrorMaxFrequencyLimit error = errors.New("frequency limit reached")
	configFile                   = ""
)

//go:embed templates/*.html
var templateFS embed.FS

type EmailClient struct {
	*resend.Client
	config *conf.GlobalConfiguration
}

func NewEmailClient(config *conf.GlobalConfiguration) *EmailClient {
	client := resend.NewClient(config.Mail.ResendApiKey)
	return &EmailClient{client, config}
}

type ConfirmationEmailData struct {
	UserName  string
	VerifyUrl string
}

func (c *EmailClient) SendConfirmationEmail(
	txx *sqlx.Tx, user *models.User, emailTo string, maxFreq time.Duration,
) error {
	if user.ConfirmationSentAt != nil && !user.ConfirmationSentAt.Add(maxFreq).Before(time.Now()) {
		return ErrorMaxFrequencyLimit
	}

	secToken, err := crypto.GenerateSecureToken()
	if err != nil {
		return fmt.Errorf("error generating confirmation token: %s", err)
	}

	if err := c.parseHtmlandSendConfirmation(user, emailTo, secToken); err != nil {
		return fmt.Errorf("error while sending confirmation email: %w", err)
	}

	now := time.Now()

	query := `UPDATE app.users SET confirmation_token = $1, confirmation_sent_at = $2 where id = $3`
	_, err = txx.Exec(query, secToken, now, user.ID)
	if err != nil {
		return fmt.Errorf("error setting confirmation token: %s", err)
	}

	return nil
}

func (c *EmailClient) parseHtmlandSendConfirmation(user *models.User, emailTo string, secToken string) error {
	// Parse the embedded template
	tmpl, err := template.ParseFS(templateFS, "templates/confirmation.html")
	if err != nil {
		return fmt.Errorf("error parsing email template: %w", err)
	}

	// Prepare template data
	verifyURL := fmt.Sprintf("%s/verify?token=%s", c.config.App.Host, secToken)
	data := ConfirmationEmailData{
		UserName:  user.Username,
		VerifyUrl: verifyURL,
	}

	// Execute template
	var htmlBody bytes.Buffer
	if err := tmpl.Execute(&htmlBody, data); err != nil {
		return fmt.Errorf("error executing email template: %w", err)
	}

	// Send email using Resend SDK
	params := &resend.SendEmailRequest{
		From:    c.config.Mail.ResendFromEmail,
		To:      []string{emailTo},
		Subject: "メールアドレスの確認",
		Html:    htmlBody.String(),
	}

	fmt.Println("sending to ", emailTo)

	_, err = c.Client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("error sending confirmation email: %w", err)
	}

	return nil
}
