package steam

import (
	"encoding/json"
	"regexp"
	"errors"
	"time"
	"net/url"
	"strconv"
	"net/http"
	"fmt"
	"github.com/lucasmbaia/goskins/steam-api/request"
)

type EconItemDesc struct {
	ClassID         uint64        `json:"classid,string"`    // for matching with EconItem
	InstanceID      uint64        `json:"instanceid,string"` // for matching with EconItem
	Tradable        int           `json:"tradable"`
	BackgroundColor string        `json:"background_color"`
	IconURL         string        `json:"icon_url"`
	IconLargeURL    string        `json:"icon_url_large"`
	IconDragURL     string        `json:"icon_drag_url"`
	Name            string        `json:"name"`
	NameColor       string        `json:"name_color"`
	MarketName      string        `json:"market_name"`
	MarketHashName  string        `json:"market_hash_name"`
	Comodity        bool          `json:"comodity"`
	Actions         []*EconAction `json:"actions"`
	Tags            []*EconTag    `json:"tags"`
	Descriptions    []*EconDesc   `json:"descriptions"`
}

type EconAction struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

type EconItem struct {
	AssetID    uint64 `json:"assetid,string,omitempty"`
	InstanceID uint64 `json:"instanceid,string,omitempty"`
	ClassID    uint64 `json:"classid,string,omitempty"`
	AppID      uint32 `json:"appid"`
	ContextID  uint64 `json:"contextid,string"`
	Amount     uint16 `json:"amount,string"`
	Missing    bool   `json:"missing,omitempty"`
}

type EconDesc struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	Color string `json:"color"`
}

type EconTag struct {
	InternalName string `json:"internal_name"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	CategoryName string `json:"category_name"`
}

type PartnerInventory struct {
	Success		string	`json:"success,omitempty"`
	RgInventory	map[string]struct{
		ID	    string  `json:"id,omitempty"`
		ClassID	    string  `json:"classid,omitempty"`
		InstanceID  string  `json:"instanceid,omitempty"`
		Amount	    string  `json:"amount,omitempty"`
		Pos	    int	    `json:"pos,omitempty"`
	} `json:"rgInventory,omitempty"`
	RgDescription	map[string]struct{
		EconItemDesc
	} `json:"rgDescriptions"`
	More		bool	`json:"more,omitempty"`
	MoreStart	bool    `json:"more_start,omitempty"`
}

type TradeOffer struct {
	ID                 uint64      `json:"tradeofferid,string"`
	Partner            uint32      `json:"accountid_other"`
	ReceiptID          uint64      `json:"tradeid,string"`
	RecvItems          []*EconItem `json:"items_to_receive"`
	SendItems          []*EconItem `json:"items_to_give"`
	Message            string      `json:"message"`
	State              uint8       `json:"trade_offer_state"`
	ConfirmationMethod uint8       `json:"confirmation_method"`
	Created            int64       `json:"time_created"`
	Updated            int64       `json:"time_updated"`
	Expires            int64       `json:"expiration_time"`
	EscrowEndDate      int64       `json:"escrow_end_date"`
	RealTime           bool        `json:"from_real_time_trade"`
	IsOurOffer         bool        `json:"is_our_offer"`
}

type TradeOfferResponse struct {
	Offer          *TradeOffer     `json:"offer"`                 // GetTradeOffer
	SentOffers     []*TradeOffer   `json:"trade_offers_sent"`     // GetTradeOffers
	ReceivedOffers []*TradeOffer   `json:"trade_offers_received"` // GetTradeOffers
	Descriptions   []*EconItemDesc `json:"descriptions"`          // GetTradeOffers
}

type APIResponse struct {
	Inner *TradeOfferResponse `json:"response"`
}

type OfferItems struct {
	Assets	    []EconItem  `json:"assets"`
	Currency    []struct{}	`json:"currency"`
	Ready	    bool	`json:"ready"`
}

type RequestTraderOffer struct {
	NewVersion  bool	`json:"newversion"`
	Version	    int		`json:"version"`
	Me	    OfferItems	`json:"me"`
	Them	    OfferItems	`json:"them"`
}

const (
	TradeFilterNone             = iota
	TradeFilterSentOffers       = 1 << 0
	TradeFilterRecvOffers       = 1 << 1
	TradeFilterActiveOnly       = 1 << 3
	TradeFilterHistoricalOnly   = 1 << 4
	TradeFilterItemDescriptions = 1 << 5
)

var (
	// receiptExp matches JSON in the following form:
	//      oItem = {"id":"...",...}; (Javascript code)
	receiptExp    = regexp.MustCompile("oItem =\\s(.+?});")
	myEscrowExp   = regexp.MustCompile("var g_daysMyEscrow = (\\d+);")
	themEscrowExp = regexp.MustCompile("var g_daysTheirEscrow = (\\d+);")
	errorMsgExp   = regexp.MustCompile("<div id=\"error_msg\">\\s*([^<]+)\\s*</div>")
	offerInfoExp  = regexp.MustCompile("token=([a-zA-Z0-9-_]+)")

	getTradeOffer     = "https://api.steampowered.com/IEconService/GetTradeOffer/v1/?"
	getTradeOffers    = "https://api.steampowered.com/IEconService/GetTradeOffers/v1/?"
	declineTradeOffer = "https://api.steampowered.com/IEconService/DeclineTradeOffer/v1/"
	cancelTradeOffer  = "https://api.steampowered.com/IEconService/CancelTradeOffer/v1/"
	partnerInventory  = "https://steamcommunity.com/tradeoffer/new/partnerinventory/?sessionid=%s&partner=%d&appid=%d&contextid=%d"
	sendTradeOffer	    = "https://steamcommunity.com/tradeoffer/new/send"

	ErrReceiptMatch        = errors.New("unable to match items in trade receipt")
	ErrCannotAcceptActive  = errors.New("unable to accept a non-active trade")
	ErrCannotFindOfferInfo = errors.New("unable to match data from trade offer url")
)

func tbit(bits uint32, bit uint32) bool {
	return (bits & bit) == bit
}

func (s *Session) GetPartnerInventory(sid, appID, contextID uint64, referer string) (pi PartnerInventory, err error) {
	var (
		u	*url.URL
		resp	request.Response
	)

	if u, err = url.Parse(steamCommunity); err != nil {
		return
	}

	s.client.SetCookies(u, []*http.Cookie{
		{Name: "bCompletedTradeOfferTutorial", Value: "true"},
		{Name: "sessionid", Value: s.sessionID},
	})

	if resp, err = s.client.Request(request.GET, fmt.Sprintf(partnerInventory, s.sessionID, sid, appID, contextID), request.Options{
		Headers:    map[string]string{
			"X-Requested-With": httpXRequestedWithValue,
			"Referer":	    referer,
			"User-Agent":	    httpUserAgentValue,
			"Accept":	    httpAcceptValue,
			"Host":		    "steamcommunity.com",
		},
	}); err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body, &pi); err != nil {
		return
	}

	return
}

func (s *Session) SendTraderOffer(rto *RequestTraderOffer, sid SteamID, partner, token string) (err error) {
	var (
		body	[]byte
		u	*url.URL
		params	string
		resp	request.Response
	)

	if body, err = json.Marshal(rto); err != nil {
		return
	}

	fmt.Println(string(body))

	if u, err = url.Parse(steamCommunity); err != nil {
		return
	}

	s.client.SetCookies(u, []*http.Cookie{
		{Name: "bCompletedTradeOfferTutorial", Value: "true"},
		{Name: "sessionid", Value: s.sessionID},
	})

	params = url.Values{
		"sessionid":                 {s.sessionID},
		"serverid":                  {"1"},
		"partner":                   {sid.ToString()},
		"tradeoffermessage":         {""},
		"json_tradeoffer":           {string(body)},
		"captcha":		     {""},
		"trade_offer_create_params": {"{\"trade_offer_access_token\":\"" + token + "\"}"},
	}.Encode()

	fmt.Println(params)
	fmt.Println(sid.ToString())

	if resp, err = s.client.Request(request.POST, sendTradeOffer, request.Options{
		Headers:    map[string]string{
			//"X-Requested-With": httpXRequestedWithValue,
			"Referer":	    partner,
			//"User-Agent":	    httpUserAgentValue,
			//"Accept":	    httpAcceptValue,
			//"Host":		    "steamcommunity.com",
			//"Origin":	    "https://steamcommunity.com",
			"Content-Type":	    "application/x-www-form-urlencoded",
		},
		BodyS:	    params,
	}); err != nil {
		return
	}

	fmt.Println(string(resp.Body), resp.Code)
	return
}

func (s *Session) GetTradeOffers(filter uint32, t time.Time) (traders *TradeOfferResponse, err error) {
	var (
		params	url.Values
		resp	request.Response
		api	APIResponse
	)

	params = url.Values{
		"key": {s.apiKey},
	}

	if tbit(filter, TradeFilterSentOffers) {
		params.Set("get_sent_offers", "1")
	}

	if tbit(filter, TradeFilterRecvOffers) {
		params.Set("get_received_offers", "1")
	}

	if tbit(filter, TradeFilterActiveOnly) {
		params.Set("active_only", "1")
	}

	if tbit(filter, TradeFilterItemDescriptions) {
		params.Set("get_descriptions", "1")
	}

	if tbit(filter, TradeFilterHistoricalOnly) {
		params.Set("historical_only", "1")
		params.Set("time_historical_cutoff", strconv.FormatInt(t.Unix(), 10))
	}

	fmt.Println(getTradeOffers + params.Encode())
	if resp, err = s.client.Request(request.GET, getTradeOffers + params.Encode(), request.Options{}); err != nil {
		return
	}

	fmt.Println(string(resp.Body))
	if err = json.Unmarshal(resp.Body, &api); err != nil {
		return
	}
	traders = api.Inner

	return
}
