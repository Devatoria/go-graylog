package graylog

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"time"

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
	Client    *net.Conn
	TLSClient *tls.Conn
}

// Message represents a GELF formated message
type Message struct {
	Version      string            `json:"version"`
	Host         string            `json:"host"`
	ShortMessage string            `json:"short_message"`
	FullMessage  string            `json:"full_message,omitempty"`
	Timestamp    int64             `json:"timestamp,omitempty"`
	Level        uint              `json:"level,omitempty"`
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

// NewGraylogTLS instanciates a new graylog connection with TLS, using the given endpoint
func NewGraylogTLS(e Endpoint, timeout time.Duration, config *tls.Config) (*Graylog, error) {
	c, err := tls.DialWithDialer(&net.Dialer{Timeout: timeout}, string(e.Transport), fmt.Sprintf("%s:%d", e.Address, e.Port), config)
	if err != nil {
		return nil, err
	}

	return &Graylog{TLSClient: c}, nil
}

// Send writes the given message to the given graylog client
func (g *Graylog) Send(m Message) error {
	data, err := prepareMessage(m)
	if err != nil {
		return err
	}

	// Check if TLS client is instanciated, otherwise send using the classic client
	if g.TLSClient != nil {
		_, err = (*g.TLSClient).Write(data)
	} else {
		_, err = (*g.Client).Write(data)
	}

	return err
}

// Close closes the opened connections of the given client
func (g *Graylog) Close() error {
	if g.TLSClient != nil {
		if err := (*g.TLSClient).Close(); err != nil {
			return err
		}
	}

	if g.Client != nil {
		if err := (*g.Client).Close(); err != nil {
			return err
		}
	}

	return nil
}

// prepareMessage marshal the given message, add extra fields and append EOL symbols
func prepareMessage(m Message) ([]byte, error) {
	// Marshal the GELF message in order to get base JSON
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		return []byte{}, err
	}

	// Parse JSON in order to dynamically edit it
	c, err := gabs.ParseJSON(jsonMessage)
	if err != nil {
		return []byte{}, err
	}

	// Loop on extra fields and inject them into JSON
	for key, value := range m.Extra {
		_, err = c.Set(value, fmt.Sprintf("_%s", key))
		if err != nil {
			return []byte{}, err
		}
	}

	// Append the \n\0 sequence to the given message in order to indicate
	// to graylog the end of the message
	data := append(c.Bytes(), '\n', byte(0))

	return data, nil
}
