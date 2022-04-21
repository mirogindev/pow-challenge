package protocol

import (
	"errors"
	"fmt"
	"strings"
)

type Message struct {
	Header  string
	Payload []byte
}

func (m *Message) Encode() string {
	if m.Payload == nil {
		return fmt.Sprintf("%s##empty", m.Header)
	}

	return fmt.Sprintf("%s##%s", m.Header, m.Payload)
}

func DecodeMessage(s string) (*Message, error) {
	if s == "" {
		return nil, errors.New("message must not be empty")
	}
	m := &Message{}
	parsed := strings.Split(strings.TrimSpace(s), "##")
	if len(parsed) == 2 && parsed[1] != "" {
		m.Payload = []byte(strings.TrimSpace(parsed[1]))
	} else {
		m.Payload = []byte("empty")
	}
	m.Header = strings.TrimSpace(parsed[0])
	return m, nil
}
