package main

import (
	"flag"
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/koltyakov/gosip/cpass"
)

func main() {

	var rawSecret string
	var mode string

	flag.StringVar(&rawSecret, "secret", "", "Raw secret string")
	flag.StringVar(&mode, "mode", "encode", "Mode: encode/decode")
	flag.Parse()

	crypt := cpass.Cpass("")

	if rawSecret == "" {
		fmt.Printf("Password to encode: ")
		pass, _ := gopass.GetPasswdMasked()
		secret, _ := crypt.Encode(fmt.Sprintf("%s", pass))
		fmt.Println(secret)
		return
	}

	if mode == "encode" {
		secret, _ := crypt.Encode(rawSecret)
		fmt.Println(secret)
		return
	}

	if mode == "decode" {
		secret, _ := crypt.Decode(rawSecret)
		fmt.Println(secret)
		return
	}

}
