package graylog

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/Jeffail/gabs"
)

// Transport represents a transport type enum
type Transport string

const (
	UDP Transport = "udp"
	TCP Transport = "tcp"
)

// Endpoint represents a graylog endpoint
type Endpoint struct {
	Transport Transport
	Address   string
	Port      uint
}

// Graylog represents an established graylog connection
type Graylog struct {
	Client *net.Conn
}

// Message represents a GELF formated message
type Message struct {
	Version      string            `json:"version"`
	Host         string            `json:"host"`
	ShortMessage string            `json:"short_message"`
	FullMessage  string            `json:"full_message"`
	Timestamp    int64             `json:"timestamp"`
	Level        uint              `json:"level"`
	Extra        map[string]string `json:"-"`
}

// NewGraylog instanciates a new graylog connection using the given endpoint
func NewGraylog(e Endpoint) (*Graylog, error) {
	c, err := net.Dial(string(e.Transport), fmt.Sprintf("%s:%d", e.Address, e.Port))
	if err != nil {
		return nil, err
	}

	return &Graylog{Client: &c}, nil
}

// Send sends the given GELF message, injecting extra fields, prefixing them with an underscore
func (g *Graylog) Send(m Message) error {
	// Marshal the GELF message in order to get base JSON
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// Parse JSON in order to dynamically edit it
	c, err := gabs.ParseJSON(jsonMessage)
	if err != nil {
		return err
	}

	// Loop on extra fields and inject them into JSON
	for key, value := range m.Extra {
		_, err = c.Set(value, fmt.Sprintf("_%s", key))
		if err != nil {
			return err
		}
	}

	// Write data to socket
	_, err = (*g.Client).Write(c.Bytes())

	return err
}
