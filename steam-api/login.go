package steam

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"errors"
	"fmt"
	"strings"
	"math/big"
	"time"

	"github.com/lucasmbaia/goskins/steam-api/request"
)

const (
	httpXRequestedWithValue = "com.valvesoftware.android.steam.community"
	httpUserAgentValue      = "Mozilla/5.0 (Linux; U; Android 4.1.1; en-us; Google Nexus 4 - 4.1.1 - API 16 - 768x1280 Build/JRO03S) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30"
	httpAcceptValue         = "text/javascript, text/html, application/xml, text/xml, */*"

	steamCommunity		= "https://steamcommunity.com"
)

type LoginResponse struct {
	Success      bool   `json:"success,omitempty"`
	PublicKeyMod string `json:"publickey_mod,omitempty"`
	PublicKeyExp string `json:"publickey_exp,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	TokenGID     string `json:"token_gid,omitempty"`
}

type OAuth struct {
	SteamID       SteamID `json:"steamid,string"`
	Token         string  `json:"oauth_token"`
	WGToken       string  `json:"wgtoken"`
	WGTokenSecure string  `json:"wgtoken_secure"`
	WebCookie     string  `json:"webcookie"`
}

type LoginSession struct {
	Success           bool   `json:"success"`
	LoginComplete     bool   `json:"login_complete"`
	RequiresTwoFactor bool   `json:"requires_twofactor"`
	Message           string `json:"message"`
	RedirectURI       string `json:"redirect_uri"`
	OAuthInfo         string `json:"oauth"`
}

type Session struct {
	client	    *request.Client
	oauth	    OAuth
	sessionID   string
	deviceID    string
	umqID	    string
	language    string
	apiKey	    string
}

func NewSession() (s *Session, err error) {
	var client *request.Client

	if client, err = request.NewClient(); err != nil {
		return
	}

	return &Session {
		client:	    client,
		language:   "english",
	}, nil
}

func (s *Session) doLogin(response *LoginResponse, accountName, password, twoFactorCode string) (err error) {
	var (
		n	    big.Int
		exp	    int64
		rsaPKCS	    []byte
		values	    string
		resp	    request.Response
		host	    string
		ls	    LoginSession
		rb	    = make([]byte, 6)
		sessionID   []byte
		cookies	    []*http.Cookie
		u	    *url.URL
		sum	    [md5.Size]byte
	)

	n.SetString(response.PublicKeyMod, 16)

	if exp, err = strconv.ParseInt(response.PublicKeyExp, 16, 32); err != nil {
		return
	}

	if rsaPKCS, err = rsa.EncryptPKCS1v15(rand.Reader, &rsa.PublicKey{N: &n, E: int(exp)}, []byte(password)); err != nil {
		return
	}

	values = url.Values{
		"captcha_text":      {""},
		"captchagid":        {"-1"},
		"emailauth":         {""},
		"emailsteamid":      {""},
		"password":          {base64.StdEncoding.EncodeToString(rsaPKCS)},
		"remember_login":    {"true"},
		"rsatimestamp":      {response.Timestamp},
		"twofactorcode":     {twoFactorCode},
		"username":          {accountName},
		"oauth_client_id":   {"DE45CD61"},
		"oauth_scope":       {"read_profile write_profile read_client write_client"},
		"loginfriendlyname": {"#login_emailauth_friendlyname_mobile"},
		"donotcache":        {strconv.FormatInt(time.Now().Unix()*1000, 10)},
	}.Encode()

	host = fmt.Sprintf("https://steamcommunity.com/login/dologin/?%s", values)
	if resp, err = s.client.Request(request.POST, host, request.Options{
		Headers:    map[string]string{
			"X-Requested-With": httpXRequestedWithValue,
			"Referer":	    "https://steamcommunity.com/mobilelogin?oauth_client_id=DE45CD61&oauth_scope=read_profile%20write_profile%20read_client%20write_client",
			"User-Agent":	    httpUserAgentValue,
			"Accept":	    httpAcceptValue,
		},
	}); err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body, &ls); err != nil {
		return
	}

	if !ls.Success {
		if ls.RequiresTwoFactor {
			err = errors.New("Invalid TwoFactor Code")
			return
		}

		err = errors.New(ls.Message)
		return
	}

	if err = json.Unmarshal([]byte(ls.OAuthInfo), &s.oauth); err != nil {
		return
	}

	if _, err = rand.Read(rb); err != nil {
		return
	}

	sessionID = make([]byte, hex.EncodedLen(len(rb)))
	hex.Encode(sessionID, rb)
	s.sessionID = string(sessionID)

	if u, err = url.Parse(steamCommunity); err != nil {
		return
	}

	cookies = s.client.GetCookies(u)
	for _, cookie := range cookies {
		if cookie.Name == "mobileClient" || cookie.Name == "mobileClientVersion" || cookie.Name == "steamCountry" || strings.Contains(cookie.Name, "steamMachineAuth") {
			cookie.MaxAge = -1
		}
	}

	sum = md5.Sum([]byte(accountName + password))
	s.deviceID = fmt.Sprintf("android:%x-%x-%x-%x-%x", sum[:2], sum[2:4], sum[4:6], sum[6:8], sum[8:10])

	s.client.SetCookies(u, append(cookies, &http.Cookie{
		Name:	"sessionid",
		Value:	s.deviceID,
	}))

	return
}

func (s *Session) loginRequest(accountName, password string) (lr *LoginResponse, err error) {
	var (
		u	*url.URL
		resp	request.Response
		host	string
	)

	lr = &LoginResponse{}
	if u, err = url.Parse(steamCommunity); err != nil {
		return
	}

	s.client.SetCookies(u, []*http.Cookie{
		{Name: "mobileClientVersion", Value: "0 (2.1.3)"},
		{Name: "mobileClient", Value: "android"},
		{Name: "Steam_Language", Value: s.language},
		{Name: "timezoneOffset", Value: "0,0"},
	})

	host = fmt.Sprintf("https://steamcommunity.com/login/getrsakey?username=%s", accountName)
	if resp, err = s.client.Request(request.POST, host, request.Options{
		Headers:    map[string]string{
			"X-Requested-With": httpXRequestedWithValue,
			"Referer":	    "https://steamcommunity.com/mobilelogin?oauth_client_id=DE45CD61&oauth_scope=read_profile%20write_profile%20read_client%20write_client",
			"User-Agent":	    httpUserAgentValue,
			"Accept":	    httpAcceptValue,
		},
	}); err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body, lr); err != nil {
		return
	}

	if !lr.Success {
		err = errors.New("Invalid Username")
	}

	return
}

func (s *Session) Login(accountName, password, twoFactorCode string) (err error) {
	var response *LoginResponse

	if response, err = s.loginRequest(accountName, password); err != nil {
		return
	}

	err = s.doLogin(response, accountName, password, twoFactorCode)
	return
}

func (s *Session) GetSteamID() SteamID {
	return s.oauth.SteamID
}
