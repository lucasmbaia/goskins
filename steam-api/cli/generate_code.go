package main

import (
	"github.com/lucasmbaia/goskins/steam-api"
	"encoding/json"
	"io/ioutil"
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
		s	    *steam.Session
		err	    error
		sg	    steam.SteamGuard
		file	    *os.File
		username    string
		body	    []byte
		code	    string
	)

	if s, err = steam.NewSession(); err != nil {
		panic(err)
	}

	if username, err = getKeyboard("Username"); err != nil {
		panic(err)
	}

	if file, err = os.Open(fmt.Sprintf("./%s.json", username)); err != nil {
		panic(err)
	}
	defer file.Close()

	if body, err = ioutil.ReadAll(file); err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, &sg); err != nil {
		panic(err)
	}

	if code, err = s.GenerateSteamGuardCode(sg.SharedSecret); err != nil {
		panic(err)
	}

	fmt.Printf("Code: %s\n", code)
	return
}
