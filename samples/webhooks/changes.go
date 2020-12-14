package main

import (
	"fmt"

	"github.com/koltyakov/gosip/api"
)

// Tokens in memory cache
var tokens = NewCache()

func trackChanges(whs []*api.WebhookInfo) {
	for _, c := range whs {
		changeTokenStart, _ := tokens.Load(c.SubscriptionID)
		web, err := sp.Site().WebByID(c.WebID)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		changeTokenEnd, err := web.Lists().GetByID(c.Resource).Changes().GetCurrentToken()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		fmt.Printf("Changes from '%s' till '%s'\n", changeTokenStart, changeTokenEnd)
		tokens.Store(c.SubscriptionID, changeTokenEnd)
		changes, err := web.Lists().GetByID(c.Resource).Changes().GetChanges(&api.ChangeQuery{
			ChangeTokenStart: changeTokenStart,
			ChangeTokenEnd:   changeTokenEnd,
			Item:             true,
			Add:              true,
			Restore:          true,
			DeleteObject:     true,
			Update:           true,
		})
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		for _, c := range changes.Data() {
			fmt.Printf("Item ID: %d, Change Type: %s\n", c.ItemID, web.Changes().GetChangeType(c.ChangeType))
		}
		checkNextPage := true
		for checkNextPage {
			changes, err = changes.GetNextPage()
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}
			cc := changes.Data()
			for _, c := range cc {
				fmt.Printf("Item ID: %d, Change Type: %s\n", c.ItemID, web.Changes().GetChangeType(c.ChangeType))
			}
			if len(cc) == 0 {
				checkNextPage = false
			}
		}
	}
}
