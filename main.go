package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/hirnag/testgo-googleauthenticator/pkg/otp"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))
var totpSecret string
var hotpSecret string
var hotpCounter uint64 = 1

func main() {
	var err error
	totpSecret, err = otp.GenerateTOTPKey("username")
	if err != nil {
		panic(err)
	}
	hotpSecret, err = otp.GenerateHOTPKey("username")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/Authenticate", authenticate)

	fmt.Println("server starting...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if err := templates.Execute(res, nil); err != nil {
		panic(err)
	}
}

func authenticate(res http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	defer func() {
		if err := templates.Execute(res, data); err != nil {
			panic(err)
		}
	}()

	// ParseForm
	req.ParseForm()
	fmt.Println(req.Form)
	if len(req.Form["passcode"]) == 0 {
		data["Message"] = "Invalid parameter, must be input pass code."
		return
	}
	passcode := req.Form["passcode"][0]
	data["Code"] = passcode
	if len(req.Form["algorithm"]) == 0 {
		data["Message"] = "Invalid parameter, must be select algorithm."
		return
	}
	algorithm := req.Form["algorithm"][0]

	// Verify Pass Code
	var result bool
	var err error
	switch algorithm {
	case "totp":
		result, err = otp.VerifyTOTPToken(passcode, totpSecret)
		break
	case "hotp":
		result, err = otp.VerifyHOTPToken(passcode, hotpCounter, hotpSecret)
		if result {
			hotpCounter++
		}
		break
	}
	if err != nil {
		data["Message"] = err.Error()
	}
	data["Result"] = result
	fmt.Println(result)
}
