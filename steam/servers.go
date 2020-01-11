package steam

import (
	"encoding/json"
	"math/rand"
	"time"
	"fmt"

	"github.com/lucasmbaia/goskins/steam/request"
)

const (
	getCmList = "https://api.steampowered.com/ISteamDirectory/GetCMList/v1/?cellId=0"
)

type serversFields struct {
	Response struct {
		ServerList  []string	`json:"serverlist,omitempty"`
		WebSockets  []string	`json:"serverlist_websockets,omitempty"`
		Message	    string	`json:"message,omitempty"`
		Result	    uint32	`json:"result,omitempty"`
	} `json:"response,omitempty"`
}

type steamServers struct {
	servers	    []string
	webSockets  []string
}

func GetServersSteam() (s steamServers, err error) {
	var (
		resp	request.Response
		sf	serversFields
	)

	if resp, err = request.Request(request.GET, getCmList, &request.Options{}); err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body, &sf); err != nil {
		return
	}

	if resp.Code != 200 {
		err = fmt.Errorf("Failed to get steam directory, code: %d", resp.Code)
		return
	}

	if sf.Response.Result != 1 {
		err = fmt.Errorf("Failed to get steam directory, result: %v, message: %v\n", sf.Response.Result, sf.Response.Message)
		return
	}

	if len(sf.Response.ServerList) == 0 {
		err = fmt.Errorf("Steam returned zero servers for steam directory request\n")
		return
	}

	s.servers = sf.Response.ServerList
	s.webSockets = sf.Response.WebSockets
	return
}

func (s steamServers) GetRandomServer() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return s.servers[rand.Intn(len(s.servers))]
}
