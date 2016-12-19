# go-graylog

[![GoDoc](https://godoc.org/github.com/Devatoria/go-graylog?status.svg)](https://godoc.org/github.com/Devatoria/go-graylog)

Graylog GELF messages sending using UDP, TCP or TCP/TLS, written in Golang.

# Examples

```go
package main

import (
	"time"

	"github.com/Devatoria/go-graylog"
)

func main() {
    // Initialize a new graylog client with TLS
	g, err := graylog.NewGraylogTLS(graylog.Endpoint{
		Transport: graylog.TCP,
		Address:   "localhost",
		Port:      12202,
	}, 3*time.Second, nil)
	if err != nil {
		panic(err)
	}

    // Send a message
	err = g.Send(graylog.Message{
		Version:      "1.1",
		Host:         "localhost",
		ShortMessage: "Sample test",
		FullMessage:  "Stacktrace",
		Timestamp:    time.Now().Unix(),
		Level:        1,
		Extra: map[string]string{
			"MY-EXTRA-FIELD": "extra_value",
		},
	})
    if err != nil {
        panic(err)
    }

    // Close the graylog connection
    if err := g.Close(); err != nil {
        panic(err)
    }
}
```
