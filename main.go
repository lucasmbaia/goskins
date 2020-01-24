package main

import (
	"github.com/lucasmbaia/goskins/steam-api"
	"bufio"
	"fmt"
	"os"
	//"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Code: ")
	code, _ := reader.ReadString('\n')

	var s *steam.Session
	var err error

	if s, err = steam.NewSession(); err != nil {
		panic(err)
	}

	if err = s.Login("", "", code); err != nil {
		panic(err)
	}

	//s.GetPartnerInventory(76561198072709722, 730, 2, "https://steamcommunity.com/tradeoffer/new/?partner=112443994")

	var sid steam.SteamID
	sid.ParseDefaults(112443994)

	rto := &steam.RequestTraderOffer{
		NewVersion: true,
		Version:    2,
		Them:	    steam.OfferItems{
			Assets:	[]steam.EconItem{
			    {AppID: 730, ContextID: 2, Amount: 1, AssetID: 17588643249},
			},
			Currency:   make([]struct{}, 0),
		},
		Me: steam.OfferItems{
			Assets:	[]steam.EconItem{},
			Currency:   make([]struct{}, 0),
		},
	}

	err = s.SendTraderOffer(rto, sid, "https://steamcommunity.com/tradeoffer/new/?partner=112443994&token=6pSdV2m_", "6pSdV2m_")
	//err = s.SendTraderOffer(rto, sid, "https://steamcommunity.com/tradeoffer/new/?partner=112443994", "6pSdV2m_")
	if err != nil {
		panic(err)
	}
	/*_, err = s.GetWebApiKey()
	if err != nil {
		panic(err)
	}

	resp, err := s.GetTradeOffers(steam.TradeFilterSentOffers, time.Now())
	if err != nil {
		panic(err)
	}

	for _, offer := range resp.SentOffers {
		fmt.Println("***************** OFFER *************************")
		fmt.Println(offer.Partner)
		var sid steam.SteamID
		sid.ParseDefaults(offer.Partner)

		fmt.Println("Offer ID: ", offer.ID)
		fmt.Println("Partner: ", uint64(sid))
		fmt.Println("ReceiptID: ", offer.ReceiptID)
		fmt.Println("Message: ", offer.Message)
		fmt.Println("State: ", offer.State)
		fmt.Println("Created: ", offer.Created)

		for _, recv := range offer.RecvItems{
			fmt.Println(recv)
		}
	}*/

	/*apps, err := s.GetInventoryApps(s.GetSteamID())
	if err != nil {
		panic(err)
	}

	for _, app := range apps {
		fmt.Printf("************ %s ****************\n", app.Name)
		for _, context := range app.Contexts {
			fmt.Printf("-- Items on %d %d (count %d)\n", app.AppID, context.ID, context.AssetCount)
			inven, err := s.GetInventory(76561198037342607, app.AppID, context.ID, true)
			if err != nil {
				panic(err)
			}

			for _, item := range inven {
				fmt.Printf("Item: %s = %d\n", item.Desc.MarketHashName, item.AssetID)
				fmt.Println("ClassID: ", item.ClassID)
				fmt.Println("InstanceID: ", item.InstanceID)
				fmt.Println("Amount: ", item.Amount)
			}
		}
	}*/
}
