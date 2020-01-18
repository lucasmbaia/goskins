package steam

import (
	"encoding/json"
	"net/url"
	"strconv"
	"regexp"
	"errors"
	"fmt"

	"github.com/lucasmbaia/goskins/steam-api/request"
)

const (
	InventoryEndpoint = "http://steamcommunity.com/inventory/%d/%d/%d?"
)

type InventoryContext struct {
	ID         uint64 `json:"id,string"` /* Apparently context id needs at least 64 bits...  */
	AssetCount uint32 `json:"asset_count"`
	Name       string `json:"name"`

}

type Asset struct {
	AppID      uint32 `json:"appid"`
	ContextID  uint64 `json:"contextid,string"`
	AssetID    uint64 `json:"assetid,string"`
	ClassID    uint64 `json:"classid,string"`
	InstanceID uint64 `json:"instanceid,string"`
	Amount     uint64 `json:"amount,string"`
}

type Inventory struct {
	Assets              []Asset         `json:"assets"`
	Descriptions        []*EconItemDesc `json:"descriptions"`
	Success             int             `json:"success"`
	HasMore             int             `json:"more_items"`
	LastAssetID         string          `json:"last_assetid"`
	TotalInventoryCount int             `json:"total_inventory_count"`
	ErrorMsg            string          `json:"error"`
}

type InventoryItem struct {
	AppID      uint32        `json:"appid"`
	ContextID  uint64        `json:"contextid"`
	AssetID    uint64        `json:"id,string,omitempty"`
	ClassID    uint64        `json:"classid,string,omitempty"`
	InstanceID uint64        `json:"instanceid,string,omitempty"`
	Amount     uint64        `json:"amount,string"`
	Desc       *EconItemDesc `json:"-"` /* May be nil  */
}

type InventoryApps struct {
	AppID		    uint64			    `json:"appid,omitempty"`
	Name		    string			    `json:"name,omitempty"`
	AssetCount	    uint64			    `json:"asset_count,omitempty"`
	Icon		    string			    `json:"icon,omitempty"`
	Link		    string			    `json:"link,omitempty"`
	InventoryLog	    string			    `json:"inventory_logo,omitempty"`
	TradePermissions    string			    `json:"trade_permission,omitempty"`
	Contexts	    map[string]*InventoryContext    `json:"rgContexts,omitempty"`
}

var inventoryContextRegexp = regexp.MustCompile("var g_rgAppContextData = (.*?);")

type FetchInventory struct {
	sid		SteamID
	appID		uint64
	contextID	uint64
	startAssetID	uint64
	filters		[]Filter
	items		[]InventoryItem
}

func (s *Session) fetchInvetory(f *FetchInventory) (hm bool, lastAssetID uint64, err error) {
	var (
		params		url.Values
		resp		request.Response
		inventory	Inventory
		descriptions	= make(map[string]int)
	)

	params = url.Values{
		"l": {s.language},
	}

	if f.startAssetID != 0 {
		params.Set("start_assetid", strconv.FormatUint(f.startAssetID, 10))
		params.Set("count", "75")
	} else {
		params.Set("count", "100")
	}


	if resp, err = s.client.Request(request.GET, fmt.Sprintf(InventoryEndpoint, f.sid, f.appID, f.contextID) + params.Encode(), request.Options{}); err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body, &inventory); err != nil {
		return
	}

	if inventory.Success == 0 {
		if len(inventory.ErrorMsg) != 0 {
			err = errors.New(inventory.ErrorMsg)
		}

		return
	}

	for i, desc := range inventory.Descriptions {
		descriptions[fmt.Sprintf("%d_%d", desc.ClassID, desc.InstanceID)] = i
	}

	for _, asset := range inventory.Assets {
		var (
			desc	*EconItemDesc
			add	bool
			item	InventoryItem
		)

		if d, ok := descriptions[fmt.Sprintf("%d_%d", asset.ClassID, asset.InstanceID)]; ok {
			desc = inventory.Descriptions[d]
		}

		item = InventoryItem{
			AppID:      asset.AppID,
			ContextID:  asset.ContextID,
			AssetID:    asset.AssetID,
			ClassID:    asset.ClassID,
			InstanceID: asset.InstanceID,
			Amount:     asset.Amount,
			Desc:       desc,
		}

		add = true
		for _, filter := range f.filters {
			add = filter(&item)
			if !add {
				break
			}
		}

		if add {
			f.items = append(f.items, item)
		}
	}

	hm = inventory.HasMore != 0
	if !hm {
		return
	}

	if lastAssetID, err = strconv.ParseUint(inventory.LastAssetID, 10, 64); err != nil {
		return
	}

	return
}

func (s *Session) GetInventory(sid SteamID, appID, contextID uint64, tradableOnly bool) (it []InventoryItem, err error) {
	var fi = &FetchInventory{
		sid:	    sid,
		appID:	    appID,
		contextID:  contextID,
		filters:    []Filter{},
	}

	if tradableOnly {
		fi.filters = append(fi.filters, IsTradable(tradableOnly))
	}

	s.GetFilterableInventory(fi)
	it = fi.items

	return
}

func (s *Session) GetFilterableInventory(fi *FetchInventory) (err error) {
	var (
		hm	    bool
		lastAssetID uint64
	)

	fi.startAssetID = uint64(0)

	for {
		if hm, lastAssetID, err = s.fetchInvetory(fi); err != nil {
			return
		}

		if !hm {
			break
		}

		fi.startAssetID = lastAssetID
	}

	return
}

func (s *Session) GetInventoryApps(sid SteamID) (ia map[string]InventoryApps, err error) {
	var (
		resp	request.Response
		host	string
		m	[][]byte
	)

	ia = make(map[string]InventoryApps)

	host = fmt.Sprintf("https://steamcommunity.com/profiles/%s/inventory", sid.ToString())
	if resp, err = s.client.Request(request.GET, host, request.Options{}); err != nil {
		return
	}

	m = inventoryContextRegexp.FindSubmatch(resp.Body)
	if m == nil || len(m) != 2 {
		return
	}

	if err = json.Unmarshal(m[1], &ia); err != nil {
		return
	}

	return
}
