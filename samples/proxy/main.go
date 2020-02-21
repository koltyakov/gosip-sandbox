package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip"

	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
)

var debug bool

func main() {

	var strategy string
	var config string
	var port int
	var sslKey string
	var sslCert string

	flag.StringVar(&strategy, "strategy", "", "Auth strategy")
	flag.StringVar(&config, "config", "", "Config path")
	flag.IntVar(&port, "port", 9090, "Proxy port")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&sslKey, "sslKey", "", "SSL Key file path")   // openssl genrsa -out private.key 2048
	flag.StringVar(&sslCert, "sslCert", "", "SSL Crt file path") // openssl req -new -x509 -sha256 -key private.key -out public.crt -days 3650

	flag.Parse()

	auth, err := dynauth.NewAuthCnfg(strategy, config)
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	http.HandleFunc("/", proxyHandler(auth))

	if sslKey != "" && sslCert != "" {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), sslCert, sslKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}

}

func proxyHandler(ctx gosip.AuthCnfg) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		client := &gosip.SPClient{
			AuthCnfg: ctx,
		}

		siteURL, err := url.Parse(ctx.GetSiteURL())
		if err != nil {
			message := fmt.Sprintf("unable to parse site url: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		endpoint := strings.Replace(ctx.GetSiteURL(), siteURL.Path, "", -1) + r.RequestURI
		if strings.Contains(r.RequestURI, siteURL.Path) == false {
			endpoint = ctx.GetSiteURL() + r.RequestURI
		}

		req, err := http.NewRequest(r.Method, endpoint, r.Body)
		if err != nil {
			message := fmt.Sprintf("unable to create a request: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		for name, headers := range r.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				req.Header.Add(name, h)
			}
		}

		if debug {
			fmt.Printf("requesting endpoint: %s\n", endpoint)
		}
		resp, err := client.Execute(req)
		if err != nil {
			message := fmt.Sprintf("unable to request: %v\n", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		for name, headers := range resp.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				w.Header().Add(name, h)
			}
		}

		w.WriteHeader(resp.StatusCode)

		_, _ = io.Copy(w, resp.Body)
	}
}
