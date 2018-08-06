package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/hirnag/testgo-googleauthenticator/pkg/otp"
)

var templates = template.Must(template.ParseFiles("templates/index.html"))
var secret string

func main() {
	var err error
	secret, err = otp.GenerateKey("username")
	if err != nil {
		panic(err)
	}

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
	req.ParseForm()
	fmt.Println(req.Form)

	// Verify Pass Code
	result, err := otp.VerifyToken(req.Form["passcode"][0], secret)
	fmt.Println(result)
	data := map[string]interface{}{"Code": req.Form["passcode"][0], "Result": result}
	if err != nil {
		data["Message"] = err.Error()
	}

	if err := templates.Execute(res, data); err != nil {
		panic(err)
	}
}
