package otp

import (
	"fmt"
	"image/png"
	"os"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/hotp"
)

func newHOTPOpts(accountName string) hotp.GenerateOpts {
	return hotp.GenerateOpts{
		Issuer:      "hirnag",
		AccountName: accountName,
	}
}

func GenerateHOTPKey(accountName string) (string, error) {
	opts := newHOTPOpts(accountName)
	key, err := hotp.Generate(opts)
	if err != nil {
		return "", err
	}
	fmt.Printf("%+v\n", key)

	// create the output file
	qrCode, err := key.Image(180, 180)
	if err != nil {
		return "", err
	}
	file, _ := os.Create("QR_HOTP.png")
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return key.Secret(), nil
}

func VerifyHOTPToken(passcode string, counter uint64, secret string) (bool, error) {
	fmt.Printf("hotp counter: %v\n", counter)
	return hotp.ValidateCustom(
		passcode,
		counter,
		secret,
		hotp.ValidateOpts{
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
}
