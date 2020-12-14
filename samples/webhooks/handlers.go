package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/koltyakov/gosip/api"
)

// notificationsHandler handler to process webhooks notifications
// https://docs.microsoft.com/en-us/sharepoint/dev/apis/webhooks/sharepoint-webhooks-using-azure-functions
func notificationsHandler(w http.ResponseWriter, r *http.Request) {
	validationIoken := r.URL.Query().Get("validationtoken")

	if len(validationIoken) != 0 {
		log.Printf("Subscription validation is received: %s\n", validationIoken)
		fmt.Fprint(w, validationIoken)
		return
	}

	defer func() { _ = r.Body.Close() }()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	var whsInfo struct {
		Value []*api.WebhookInfo `json:"value"`
	}

	if err := json.Unmarshal(data, &whsInfo); err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	if len(whsInfo.Value) == 0 {
		fmt.Printf("empty changes\n")
		return
	}

	// Processing changes in a separate goroutine
	go trackChanges(whsInfo.Value)
}

// subscribeHandler handler to subscribe list notifications
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	listName := r.URL.Query().Get("listName")
	if len(listName) == 0 {
		http.Error(w, "Url parameter 'listName' is missing", 400)
		return
	}

	webhookURL := getNotificationsURL(r)
	if len(webhookURL) == 0 {
		http.Error(w, "Url parameter 'notificationsUrl' is missing", 400)
		return
	}

	expiration := time.Now().Add(10 * time.Minute)

	sub, err := sp.Web().Lists().GetByTitle(listName).Subscriptions().Add(webhookURL, expiration, "test webhook")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

// unsubscribeHandler handler to unsubscribe list notifications
func unsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	listName := r.URL.Query().Get("listName")
	if len(listName) == 0 {
		http.Error(w, "Url parameter 'listName' is missing", 400)
		return
	}

	subscriptionID := r.URL.Query().Get("subscriptionId")
	if len(subscriptionID) == 0 {
		http.Error(w, "Url parameter 'subscriptionId' is missing", 400)
		return
	}

	list := sp.Web().Lists().GetByTitle(listName)
	if err := list.Subscriptions().GetByID(subscriptionID).Delete(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{ "result": "ok" }`)
}

// subscriptionsHandler handler for getting a list of subscriptions
func subscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	listName := r.URL.Query().Get("listName")
	if len(listName) == 0 {
		http.Error(w, "Url parameter 'listName' is missing", 400)
		return
	}

	subs, err := sp.Web().Lists().GetByTitle(listName).Subscriptions().Get()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if subs == nil {
		fmt.Fprint(w, "[]")
		return
	}

	json.NewEncoder(w).Encode(subs)
}
