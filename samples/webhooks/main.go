package main

import (
	"log"
	"net/http"
	"os"
)

func init() {

}

func main() {
	listenAddr := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddr = ":" + val
	}

	http.HandleFunc("/api/notifications", notificationsHandler)
	http.HandleFunc("/api/subscribe", subscribeHandler)
	http.HandleFunc("/api/unsubscribe", unsubscribeHandler)
	http.HandleFunc("/api/subscriptions", subscriptionsHandler)

	log.Printf("Server has been started at http://127.0.0.1%s", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
