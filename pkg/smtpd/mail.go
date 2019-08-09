package smtpd

import (
	"fmt"
	"strings"
)

// Envelope contains the sender address (originator or return-path).
type Envelope struct {
	Sender     string   `json:"sender"`
	Recipients []string `json:"recipients"`
}

// Mail object contains an envelope and content as described in RFC 5321 §2.3.1.
type Mail struct {
	Envelope Envelope `json:"envelope"`
	Content  []string `json:"content"`
}

func (m Mail) String() string {
	return fmt.Sprintf("MAIL FROM:%v\nRCPT TO:%v\n%v", m.Envelope.Sender, strings.Join(m.Envelope.Recipients, ", "), strings.Join(m.Content, "\n"))
}
