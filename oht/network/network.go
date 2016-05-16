package network

import (
	"log"
	"net/url"
)

func main() {
	Connect("http://cool.com")
	Connect("oht:cool.com")
	Connect("ricochet:testos.com")
}

func Listen(protocol, listenHost, listenPort string) {
	if protocol == "http" {
		log.Println("http")
		transports.InitializeHTTP(listenHost, listenPort)
	} else if protocol == "oht" {
		log.Println("oht")
		transports.InitializeTCP(listenHost, listenPort)
	} else if protocol == "ricochet" {
		log.Println("ricochet")
		transports.InitializeTCP(listenHost, listenPort)
	}
}

func Connect(peerUrl string) {
	// Break up the URL into
	// protocol :// hostname, throw out any non onion hosted addresses
	parsedUrl, err := url.Parse(peerUrl)
	if err != nil {
		log.Println("Network: Error! ", err)
	}
	if parsedUrl.Scheme == "http" {
		log.Println("HTTP scheme")
		log.Println(parsedUrl.Host)
	} else if parsedUrl.Scheme == "ricochet" {
		log.Println("RICOCHET scheme")
		log.Println(parsedUrl.Opaque)
		//
	} else if parsedUrl.Scheme == "oht" {
		log.Println("OHT scheme")
		log.Println(parsedUrl.Opaque)
	}
}
