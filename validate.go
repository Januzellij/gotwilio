package gotwilio

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"sort"
)

// sortedFormString takes POST parameters and returns the params concatenated like "keyvaluekeyvaluevalue"
func sortedFormString(f url.Values) string {
	// move the dictionary into two slices for sorting
	keys := make([]string, len(f))
	values := make([][]string, len(f))

	i := 0
	for k, v := range f {
		keys[i] = k
		values[i] = v
		i++
	}

	// params must be sorted in alphabetical order
	sort.Strings(keys)

	// concatenate the params together
	var b bytes.Buffer
	for _, val := range keys {
		b.WriteString(val)
		for _, value := range f[val] {
			b.WriteString(value)
		}
	}

	return b.String()
}

// Validate checks if an http request actually came from Twilio, and is not faked.
// Validate uses directions from https://www.twilio.com/docs/security
func Validate(r *http.Request, url, authToken string) (bool, error) {
	var urlString string

	// if the request is a POST request, get the string of the form
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			return false, err
		}
		rawForm := r.PostForm
		formString := sortedFormString(rawForm)
		urlString = url + formString
	} else {
		urlString = url
	}

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(urlString))

	var b bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &b)
	macBytes := mac.Sum(nil)
	encoder.Write(macBytes)
	encoder.Close()

	// fetch the given twilio signature and compare against ours
	twilioSig := r.Header.Get("X-Twilio-Signature")
	if twilioSig == "" {
		return false, errors.New("Request is missing X-Twilio-Signature header")
	}
	return bytes.Equal([]byte(twilioSig), b.Bytes()), nil
}
