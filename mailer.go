package mailer

// Email describes attributes that an email has
type Email struct {
	Subject  string
	HTMLBody string
	TextBody string
	CharSet  string
}

// Destination describes who an email will be sent to. We add a level
// of nesting to make it easier to add support for CC and BCC later
type Destination struct {
	ToAddresses []string
}

// SendEmailInput combines an email with send information, and is
// meant to be passed to a Mailer
type SendEmailInput struct {
	Destination Destination
	Sender      string
	Email       Email
}

// SendEmailResult returns any implementation-specific details that should
// be returned to the user
type SendEmailResult interface {
	// ID returns the ID that can be used to identify this email. If
	// the ID isn't natively a string, either 1) it will be stringified or 2)
	// this library will be updated to handle type differences in a more
	// elegant way.
	ID() string
}

// Mailer describes a common interface for sending emails
type Mailer interface {
	SendEmail(input *SendEmailInput) (*SendEmailResult, error)
}
