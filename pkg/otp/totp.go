package otp

import (
	"fmt"
	"image/png"
	"os"

	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func newOpts(accountName string) totp.GenerateOpts {
	return totp.GenerateOpts{
		Issuer:      "hirnag",
		AccountName: accountName,
	}
}

func GenerateKey(accountName string) (string, error) {
	opts := newOpts(accountName)
	key, err := totp.Generate(opts)
	if err != nil {
		return "", err
	}
	fmt.Printf("%+v\n", key)

	// create the output file
	qrCode, err := key.Image(180, 180)
	if err != nil {
		return "", err
	}
	file, _ := os.Create("QRCODE.png")
	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return key.Secret(), nil
}

func VerifyToken(passcode string, secret string) (bool, error) {
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
