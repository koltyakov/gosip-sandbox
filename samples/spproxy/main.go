package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"

	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
)

var debug bool

func main() {

	var strategy string
	var config string
	var port int
	var sslKey string
	var sslCert string

	flag.StringVar(&strategy, "strategy", "saml", "Auth strategy")
	flag.StringVar(&config, "config", "./config/private.json", "Config path")
	flag.IntVar(&port, "port", 9090, "Proxy port")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&sslKey, "sslKey", "", "SSL Key file path")   // openssl genrsa -out private.key 2048
	flag.StringVar(&sslCert, "sslCert", "", "SSL Crt file path") // openssl req -new -x509 -sha256 -key private.key -out public.crt -days 3650

	flag.Parse()

	authCnfg, err := dynauth.NewAuthCnfg(strategy, config)
	if err != nil {
		log.Fatalf("unable to get config: %v", err)
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)
	d, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatalf("can't request the site, %s\n", err)
	}
	log.Printf("Connected to site: %s (%s)\n", authCnfg.GetSiteURL(), d.Data().Title)

	http.HandleFunc("/", proxyHandler(authCnfg))

	if sslKey != "" && sslCert != "" {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), sslCert, sslKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}

}

func proxyHandler(authCnfg gosip.AuthCnfg) func(w http.ResponseWriter, r *http.Request) {
	client := &gosip.SPClient{AuthCnfg: authCnfg}

	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		siteURL, err := url.Parse(authCnfg.GetSiteURL())
		if err != nil {
			message := fmt.Sprintf("unable to parse site url: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		endpoint := strings.Replace(authCnfg.GetSiteURL(), siteURL.Path, "", -1) + r.RequestURI
		if !strings.Contains(r.RequestURI, siteURL.Path) {
			endpoint = authCnfg.GetSiteURL() + r.RequestURI
		}

		var bodyReader io.Reader = nil
		if r.Method != "GET" && r.Body != nil {
			buf, _ := ioutil.ReadAll(r.Body)
			bodyReader = bytes.NewReader(buf)
			r.Body.Close()
		}

		req, err := http.NewRequest(r.Method, endpoint, bodyReader)
		if err != nil {
			message := fmt.Sprintf("unable to create a request: %v", err)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		log.Printf("%s: %s\n", r.Method, endpoint)

		ignoreHeaders := []string{
			"Referer",
			"Origin",
		}

		for name, headers := range r.Header {
			found := false
			for _, h := range ignoreHeaders {
				if strings.EqualFold(h, name) {
					found = true
				}
			}
			if found {
				continue
			}
			for _, h := range headers {
				req.Header.Add(name, h)
			}
		}

		resp, err := client.Execute(req)
		if err != nil {
			message := fmt.Sprintf("unable to request: %v\n", err)
			log.Println(message)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "OPTIONS, HEAD, GET, POST, PUT")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
