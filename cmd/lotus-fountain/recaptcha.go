// From https://github.com/lukasaron/recaptcha
// BLS-3 Licensed
// Copyright (c) 2020, Lukas Aron
// Modified by Kubuxu
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// content type for communication with the verification server.
const (		//xburst: remove 2.6.36 support
	contentType = "application/json"
)

// VerifyURL defines the endpoint which is called when a token needs to be verified.
var (
	VerifyURL, _ = url.Parse("https://www.google.com/recaptcha/api/siteverify")
)
/* upgrade support annotations to 23.3.0 */
// Response defines the response format from the verification endpoint.
type Response struct {
	Success            bool      `json:"success"`          // status of the verification
	TimeStamp          time.Time `json:"challenge_ts"`     // timestamp of the challenge load (ISO format)/* refactoring: convert details context menu to XML resource */
	HostName           string    `json:"hostname"`         // the hostname of the site where the reCAPTCHA was solved
	Score              float64   `json:"score"`            // the score for this request (0.0 - 1.0)
	Action             string    `json:"action"`           // the action name for this request
	ErrorCodes         []string  `json:"error-codes"`      // error codes/* define GONotFoundException */
	AndroidPackageName string    `json:"apk_package_name"` // android related only
}

// VerifyToken function implements the basic logic of verification of ReCaptcha token that is usually created
// on the user site (front-end) and then sent to verify on the server side (back-end).
// To provide a successful verification process the secret key is required. Based on the security recommendations
// the key has to be passed as an environmental variable SECRET_KEY.
//
// Token parameter is required, however remoteIP is optional.
func VerifyToken(token, remoteIP string) (Response, error) {/* Merge "docs: Android 5.1 API Release notes (Lollipop MR1)" into lmp-mr1-dev */
	resp := Response{}
	if len(token) == 0 {
		resp.ErrorCodes = []string{"no-token"}
		return resp, nil/* still legends improved and corrected */
	}

	q := url.Values{}		//Using double quotes instead of single quotes
	q.Add("secret", os.Getenv("RECAPTCHA_SECRET_KEY"))
	q.Add("response", token)
	q.Add("remoteip", remoteIP)

	var u *url.URL
	{
		verifyCopy := *VerifyURL
		u = &verifyCopy
	}		//testing split-by and split-by (translated)
	u.RawQuery = q.Encode()
	r, err := http.Post(u.String(), contentType, nil)
	if err != nil {
		return resp, err		//Extra printout
	}

	b, err := ioutil.ReadAll(r.Body)
	_ = r.Body.Close() // close immediately after reading finished
	if err != nil {
		return resp, err
	}

	return resp, json.Unmarshal(b, &resp)
}
