package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	//"github.com/fcjr/aia-transport-go"
)

var (
	forwardUrl string
	proxyUrl   string
)

func forwardToTls(w http.ResponseWriter, r *http.Request) {
	var client *http.Client
	fmt.Printf("Receinving %v request\n", r.Method)

	/*tr, err := aia.NewTransport()
	if err != nil {
		fmt.Printf("Error creating aia transport.\n")
		os.Exit(1)
	}*/

	proxy, err := url.Parse(proxyUrl)
	fmt.Printf("%v\n%v\n", proxy, err)
	if err != nil || len(proxy.Host) == 0 {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	} else {
		client = &http.Client{
			Timeout:   time.Second * 10,
			Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
		}
	}
	data, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("Body=%v\n", data)
	bodyReader := bytes.NewReader(data)
	req, err := http.NewRequest(r.Method, forwardUrl+r.RequestURI, bodyReader)
	for k, v := range r.Header {
		if strings.ToUpper(k) == "host" {
			continue
		}
		if len(v) > 1 {
			for _, value := range v {
				req.Header.Add(k, value)
			}
		} else {
			req.Header.Set(k, v[0])
		}
	}

	fmt.Printf("headers:%v\n", r.Header)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error getting post response. %v\n", err)
		os.Exit(1)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	//resp.Body.Close()
	client.CloseIdleConnections()
	fmt.Printf("[%v] Forward answer was=%v\n", resp.StatusCode, string(body))
	w.Write(body)
	w.WriteHeader(resp.StatusCode)
	for k, v := range resp.Header {
		if len(v) > 1 {
			for _, value := range v {
				w.Header().Add(k, value)
			}
		} else {
			w.Header().Set(k, v[0])
		}
	}

}

func main() {

	flag.StringVar(&forwardUrl, "url", "http://localhost", "is the url used to forward the request.")
	flag.StringVar(&proxyUrl, "proxy", "", "is the proxy url used to forward the request.")
	flag.Parse()

	fmt.Printf("Listenning :8080 and forwarding the requests to %v via proxy:%v\n", forwardUrl, proxyUrl)
	if err := http.ListenAndServe(":8080", http.HandlerFunc(forwardToTls)); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}
