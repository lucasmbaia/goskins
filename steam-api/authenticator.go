package steam

import (
	"encoding/json"
	"crypto/rand"
	"encoding/hex"
	"crypto/sha1"
	"net/url"
	"errors"
	"strconv"
	"fmt"

	"github.com/lucasmbaia/goskins/steam-api/request"
)

const (
	addAuthenticator	= "https://api.steampowered.com/ITwoFactorService/AddAuthenticator/v0001"
	finalizeAuthenticator	= "https://api.steampowered.com/ITwoFactorService/FinalizeAddAuthenticator/v0001"
)

type AuthenticatorResponse struct {
	SharedSecret	string	`json:"shared_secret"`
	SerialNumber	string	`json:"serial_number"`
	RevocationCode	string	`json:"revocation_code"`
	URI		string	`json:"uri"`
	ServerTime	string	`json:"server_time"`
	AccountName	string	`json:"account_name"`
	TokenGID	string	`json:"token_gid"`
	IdentitySecret	string	`json:"identity_secret"`
	Secret		string	`json:"secret_1"`
	Status		int	`json:"status"`
}

type Authenticator struct {
	Response    AuthenticatorResponse   `json:"response"`
}

type FinalizeAuthenticatorResponse struct {
	Response struct {
		Status     int       `json:"status"`
		ServerTime timestamp `json:"server_time"`
		WantMore   bool      `json:"want_more"`
		Success    bool      `json:"success"`
	} `json:"response"`
}

func (s *Session) NewDeviceID() (id string, err error) {
	var buf = make([]byte, 8)

	if _, err = rand.Read(buf); err != nil {
		return
	}

	var sum = sha1.Sum(buf)
	id = fmt.Sprintf("android:%s", hex.EncodeToString(sum[:]))

	return
}

func (s *Session) AddAuthenticator(device string) (auth Authenticator, err error) {
	var (
		params	url.Values
		resp	request.Response
	)

	params = url.Values{
		"access_token":		[]string{s.oauth.Token},
		"steamid":		[]string{s.oauth.SteamID.ToString()},
		"authenticator_type":	[]string{"1"},
		"device_identifier":	[]string{device},
		"sms_phone_id":		[]string{"1"},
	}

	if resp, err = s.client.Request(request.POST, addAuthenticator, request.Options{
		Headers:    map[string]string{
			"User-Agent":	    httpUserAgentValue,
			"Accept":	    httpAcceptValue,
			"Referrer":	    "https://steamcommunity.com/mobilelogin?oauth_client_id=DE45CD61&oauth_scope=read_profile%20write_profile%20read_client%20write_client",
			"Content-Type":	    "application/x-www-form-urlencoded; charset=UTF-8",
		},
		BodyS:	params.Encode(),
	}); err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body, &auth); err != nil {
		return
	}

	if auth.Response.Status != 1 {
		err = errors.New(fmt.Sprintf("Protocol error: Response.Status was %d, expected 1", auth.Response.Status))
	}

	return
}

func (s *Session) FinalizeAddAuthenticator(smsCode, secret string) (err error) {
	var (
		params	url.Values
		resp	request.Response
		code	string
		finish	FinalizeAuthenticatorResponse
	)

	params = url.Values{
		"access_token":		[]string{s.oauth.Token},
		"activation_code":	[]string{smsCode},
		"steamid":		[]string{s.oauth.SteamID.ToString()},
	}

	for tries := 0; tries <= 30; tries++ {
		if tries == 0 {
			params.Set("authenticator_code", "")
		} else {
			if code, err = s.GenerateSteamGuardCode(secret); err != nil {
				return
			}

			params.Set("authenticator_code", code)
		}

		params.Set("authenticator_time", strconv.FormatInt(s.GetSteamTime().Unix(), 10))

		if resp, err = s.client.Request(request.POST, finalizeAuthenticator, request.Options{
			Headers:    map[string]string{
				"User-Agent":	    httpUserAgentValue,
				"Accept":	    httpAcceptValue,
				"Referrer":	    "https://steamcommunity.com/mobilelogin?oauth_client_id=DE45CD61&oauth_scope=read_profile%20write_profile%20read_client%20write_client",
				"Content-Type":	    "application/x-www-form-urlencoded; charset=UTF-8",
			},
			BodyS:	params.Encode(),
		}); err != nil {
			return
		}

		if err = json.Unmarshal(resp.Body, &finish); err != nil {
			return
		}

		if finish.Response.Status == 89 {
			err = errors.New("Code 89")
			return
		}

		if !finish.Response.Success {
			err = errors.New("Success is false")
			return
		}

		if !finish.Response.WantMore {
			return
		}

		if tries == 30 {
			err = errors.New("Max retry excedid")
		}
	}

	return
}
