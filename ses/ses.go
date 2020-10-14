package ses

import (
	//go get -u github.com/aws/aws-sdk-go

	"time"

	"github.com/Synthesis-AI-Dev/mailer"
	"github.com/Synthesis-AI-Dev/mailer/lib/generic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	log "github.com/sirupsen/logrus"
)

// SES encapsulates an Amazon SES client and provides an implementation
// for SendEmail
type SES struct {
	client *ses.SES
}

// Result contains the ID assigned to this email by SES
type Result struct {
	SESOutput *ses.SendEmailOutput
}

// ID returns the ID of this email as assigned by SES
func (s Result) ID() string {
	return *s.SESOutput.MessageId
}

// New initializes a session with the provided aws config and returns a new
// SES client instance
func New(config aws.Config) *SES {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            config,
		SharedConfigState: session.SharedConfigEnable,
	}))

	s := SES{client: ses.New(session)}

	return &s
}

// SendEmail sends an email using the parameters provided through SendEmailInput
func (s *SES) SendEmail(input *mailer.SendEmailInput) (mailer.SendEmailResult, error) {
	to := input.Destination.ToAddresses

	var rec []*string
	for i := 0; i < len(to); i++ {
		rec = append(rec, aws.String(to[i]))
	}

	// Construct ses.SendEmailInput from input params passed via SendEmail
	i := &ses.SendEmailInput{
		Destination: &ses.Destination{ToAddresses: rec},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(input.Email.CharSet),
					Data:    aws.String(input.Email.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(input.Email.CharSet),
					Data:    aws.String(input.Email.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(input.Email.CharSet),
				Data:    aws.String(input.Email.Subject),
			},
		},
		Source: aws.String(input.Sender),
	}

	var result *ses.SendEmailOutput
	generic.Retry(3, 3*time.Second, func() error {
		var err error

		result, err = s.client.SendEmail(i)
		// if the error was returned by AWS then it'll almost certainly
		// send the same error again if we retry, so send a stop error
		if _, ok := err.(awserr.Error); ok {
			log.WithFields(log.Fields{
				"send_email_input": i,
			}).Error("send email fail")
			return generic.NewStopErr(err)
		}

		return nil
	})

	return Result{SESOutput: result}, nil
}
