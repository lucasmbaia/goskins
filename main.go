package main

import (
	"github.com/lucasmbaia/goskins/steam-api"
	"bufio"
	"fmt"
	"os"
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

	apps, err := s.GetInventoryApps(s.GetSteamID())
	if err != nil {
		panic(err)
	}

	for _, app := range apps {
		fmt.Printf("************ %s ****************\n", app.Name)
		for _, context := range app.Contexts {
			fmt.Printf("-- Items on %d %d (count %d)\n", app.AppID, context.ID, context.AssetCount)
			inven, err := s.GetInventory(s.GetSteamID(), app.AppID, context.ID, false)
			if err != nil {
				panic(err)
			}

			for _, item := range inven {
				fmt.Printf("Item: %s = %d\n", item.Desc.MarketHashName, item.AssetID)
			}
		}
	}
}
