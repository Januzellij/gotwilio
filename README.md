## Overview
This is an unofficial Go API library for [Twilio](http://www.twilio.com/). Gotwilio supports making voice calls, sending text messages, validating requests, creating TWiML responses, and various REST resources.

## License
Gotwilio is licensed under a BSD license.

## Installation
To install gotwilio, run `go get github.com/Januzellij/gotwilio`.

## Getting Started
Just create a Twilio client with either `NewTwilioClient(accountSid, authToken)` or `NewTwilioClientFromEnv()`, and store the accountSid and authToken in `TWILIO_ACCOUNT_SID` and `TWILIO_AUTH_TOKEN` environment variables, respectively. For security purposes, please use `NewTwilioClientFromEnv()` in any open source code.

## Docs
All documentation can be found <a href="http://godoc.org/github.com/Januzellij/gotwilio" target="_blank">here</a>

## SMS Example

```go
package main

import (
	"log"

	"github.com/Januzellij/gotwilio"
)

func main() {
	accountSid := "ABC123..........ABC123"
	authToken := "ABC123..........ABC123"
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := "+15555555555"
	to := "+15555555555"
	message := "Welcome to gotwilio!"
	_, exception, err := twilio.SendSMS(from, to, message, "", "")
	if exception != nil {
		log.Fatal(*exception)
	}
	if err != nil {
		log.Fatal(err)
	}
}
```
	
## Voice Example

```go
package main

import "github.com/Januzellij/gotwilio"

func main() {
	accountSid := "ABC123..........ABC123"
	authToken := "ABC123..........ABC123"
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := "+15555555555"
	to := "+15555555555"
	callbackParams := gotwilio.NewCallbackParameters("http://example.com")
	_, exception, err := twilio.CallWithUrlCallbacks(from, to, callbackParams)
	if exception != nil {
		log.Fatal(*exception)
	}
	if err != nil {
		log.Fatal(err)
	}
}
```

## Validate Example

```go
package main

import (
	"log"
	"net/http"

	"github.com/Januzellij/gotwilio"
)

func root(w http.ResponseWriter, r *http.Request) {
	twilio, err := NewTwilioClientFromEnv()
	if err != nil {
		// one or more environment variables are missing
		log.Fatal(err)
	}
	url := "http://example.com/"
	err = gotwilio.Validate(r, url, twilio.authToken)
	if err == nil {
		// proceed as normal, the request is from Twilio
	}
}

func main() {
	http.HandleFunc("/", root)
	http.ListenAndServe(":8080", nil)
}
```

## Twiml Response Example

```go
package main

import (
	"log"
	"os"

	"github.com/Januzellij/gotwilio"
)

func main() {
	newSay := gotwilio.Say{Text: "test", Voice: "alice"}
	newPause := gotwilio.Pause{Length: "2"}
	resp := gotwilio.NewTwimlResponse(newSay, newPause)
	err := resp.SendTwimlResponse(os.Stdout) // when using Twiml in a real web app, this would actually be written to a http.ResponseWriter.
	if err != nil {
		// something other than TWiML verbs was given to the TWiML response
		log.Fatal(err)
	}
}
```

## Usage Record Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/Januzellij/gotwilio"
)

func main() {
	twilio, err := gotwilio.NewTwilioClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	filter := &gotwilio.UsageFilter{StartDate: "2012-6-4", EndDate: "2014-1-1"}
	records, exception, err := twilio.UsageRecordsDaily(filter)
	if exception != nil {
		log.Fatal(*exception)
	}
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range records.UsageRecords {
		fmt.Printf("Category: %s, Usage: %d \n", record.Category, record.Usage)
	}
}
```