package testgo_2factorauth

import "fmt"

var apiURL = "https://chart.googleapis.com/chart?chs=%sx%s&cht=qr&chl=%s"

func GenerateFromGoogleAPI(size string, data string) {
	url := fmt.Sprintf(apiURL, size, size, data)

}
