# go-graylog
Graylog GELF messages sending using UDP, TCP or TCP/TLS, written in Golang.

# Examples

```go
package main

import (
	"time"

	"github.com/Devatoria/go-graylog"
)

func main() {
	g, err := graylog.NewGraylog(graylog.Endpoint{
		Transport: graylog.UDP,
		Address:   "localhost",
		Port:      2202,
	})
	if err != nil {
		panic(err)
	}

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
}
```
