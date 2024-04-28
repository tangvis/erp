package define

import (
	"fmt"
	"net/mail"
)

type MailInfo struct {
	Operator string
	To       []string
	Subject  string
	Content  map[string]interface{}
	Template string
}

func (m *MailInfo) Validate() error {
	if m.Operator == "" {
		return fmt.Errorf("operator is required")
	}
	if len(m.To) == 0 {
		return fmt.Errorf("at least one receiver is required")
	}
	for _, addr := range m.To {
		if _, err := mail.ParseAddress(addr); err != nil {
			return fmt.Errorf("%s wrong mail address format, %w", addr, err)
		}
	}
	return nil
}
