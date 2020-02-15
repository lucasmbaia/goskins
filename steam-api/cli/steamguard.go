package main

import (
	"github.com/lucasmbaia/goskins/steam-api"
	"encoding/json"
	"strings"
	"bufio"
	"fmt"
	"os"
)

func getKeyboard(msg string) (text string, err error) {
	var reader = bufio.NewReader(os.Stdin)

	for text == "" {
		fmt.Printf("%s: ", msg)

		if text, err = reader.ReadString('\n'); err != nil {
			return
		}

		text = strings.TrimSpace(text)
	}

	return
}

func main() {
	var (
		s		*steam.Session
		err		error
		username	string
		password	string
		mailCode	string
		smsCode		string
		deviceID	string
		authenticator	steam.Authenticator
		body		[]byte
		file		*os.File
	)

	if s, err = steam.NewSession(); err != nil {
		panic(err)
	}

	if username, err = getKeyboard("Username"); err != nil {
		panic(err)
	}

	if password, err = getKeyboard("Password"); err != nil {
		panic(err)
	}

	s.Login(username, password, "")

	if mailCode, err = getKeyboard("Mail Code"); err != nil {
		panic(err)
	}

	if err = s.Login(username, password, mailCode); err != nil {
		panic(err)
	}

	if file, err = os.Create(fmt.Sprintf("%s.json", username)); err != nil {
		panic(err)
	}
	defer file.Close()

	if deviceID, err = s.NewDeviceID(); err != nil {
		panic(err)
	}

	if authenticator, err = s.AddAuthenticator(deviceID); err != nil {
		panic(err)
	}

	if smsCode, err = getKeyboard("Sms Code"); err != nil {
		panic(err)
	}

	if s.FinalizeAddAuthenticator(smsCode, authenticator.Response.SharedSecret); err != nil {
		panic(err)
	}

	fmt.Printf("Response: %v\n", authenticator.Response)

	if body, err = json.Marshal(authenticator.Response); err != nil {
		panic(err)
	}

	if _, err = file.Write(body); err != nil {
		panic(err)
	}

	return
}
