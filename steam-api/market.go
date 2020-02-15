package steam

import (
	"github.com/lucasmbaia/goskins/steam-api/request"
	"net/url"
	//"net/http"
	"fmt"
)

const (
	getAppPriceInfo	= "https://api.steampowered.com/ISteamEconomy/GetMarketPrices/v1/?"
)

func (s *Session) GetAppPriceInfo(sid uint64, appid uint64) (err error) {
	var (
		params	url.Values
		//u	*url.URL
		resp	request.Response
	)

	//if u, err = url.Parse(steamCommunity); err != nil {
	//	return
	//}

	//s.client.SetCookies(u, []*http.Cookie{
	//	{Name: "sessionid", Value: s.sessionID},
	//})

	params = url.Values{
		"key":	    {s.apiKey},
		"appid":   {"17588643249"},
	}

	fmt.Println(getAppPriceInfo + params.Encode())
	if resp, err = s.client.Request(request.GET, getAppPriceInfo + params.Encode(), request.Options{
		//Headers:    map[string]string{
		//	"X-Requested-With": httpXRequestedWithValue,
		//	"User-Agent":	    httpUserAgentValue,
		//	"Accept":	    httpAcceptValue,
		//	"Host":		    "steamcommunity.com",
		//},
	}); err != nil {
		return
	}

	fmt.Println(string(resp.Body), resp.Code)
	return
}
