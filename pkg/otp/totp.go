package otp

import (
	"fmt"
	"image/png"
	"os"

	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func newTOTPOpts(accountName string) totp.GenerateOpts {
	return totp.GenerateOpts{
		Issuer:      "hirnag",
		AccountName: accountName,
	}
}

func GenerateTOTPKey(accountName string) (string, error) {
	opts := newTOTPOpts(accountName)
	key, err := totp.Generate(opts)
	if err != nil {
		return "", err
	}
	fmt.Println(key.String())

	// create the output file
	qrCode, err := key.Image(180, 180)
	if err != nil {
		return "", err
	}
	file, _ := os.Create("QR_TOTP.png")
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return key.Secret(), nil
}

func VerifyTOTPToken(passcode string, secret string) (bool, error) {
	return totp.ValidateCustom(
		passcode,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
}
