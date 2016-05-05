package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"./oht"
)

var (
	wui         = flag.Bool("wui", true, "Start the process with a web ui")
	username    = flag.String("username", "user", "Specify a username")
	peerAddress = flag.String("peer", "", "Specify a peer address for direct connection")
	listenPort  = flag.String("listen", "12312", "Specify a listen port")
	webUIPort   = flag.String("wuiport", "8080", "Specify a webui port")
	socksPort   = flag.String("socks", "12052", "Specify a socks proxy port")
	controlPort = flag.String("control", "9555", "Specify a control port")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	oht := oht.NewOHT(*listenPort, *socksPort, *controlPort, *webUIPort)
	log.Println("Starting " + oht.Interface.ClientInfo())
	log.Printf("\nListening for peers (Websockets) :  " + oht.Interface.TorOnionHost())
	// Connect Directly To Known Peer And Join Ring
	if *peerAddress != "" {
		if match, _ := regexp.Match(":", []byte(*peerAddress)); !match {
			*peerAddress += ":12312"
		}
		log.Printf("Connecting to peer (Websockets)  :  " + *peerAddress)
		oht.Interface.ConnectToPeer(*peerAddress)
	}
	// Start WebUI
	if *wui == true {
		oht.Interface.WebUIStart()
	}
	// Console UI
	log.Println("\nWelcome to " + oht.Interface.ClientName() + " console. Type \"/help\" to learn about the available commands.")
	prompt := "oht> "
	cli := bufio.NewScanner(os.Stdin)
	fmt.Printf(prompt)
	username := *username
	for cli.Scan() {
		body := cli.Text()
		if body == "/help" {
			fmt.Println("COMMANDS:")
			fmt.Println("  CONFIG:")
			fmt.Println("    /config                      - List configuration values (Not Implemented)")
			fmt.Println("    /set [config] [option]       - Change configuration options (Not Implemented)")
			fmt.Println("    /unset [config]              - Unset configuration option (Not Implemented)")
			fmt.Println("    /webui [start]               - Control webui (Not Implemented)")
			fmt.Println("\n  TOR:")
			fmt.Println("    /tor [start|stop]            - Start or stop tor process (Not Implemented)")
			fmt.Println("    /newtor                      - Obtain new Tor identity (Not Implemented)")
			fmt.Println("    /newonion                    - Obtain new onion address (Not Implemented)")
			fmt.Println("\n  NETWORK:")
			fmt.Println("    /peers                       - List all connected peers (Not Implemented)")
			fmt.Println("    /successor                   - Next peer in identifier ring (Not Implemented)")
			fmt.Println("    /predecessor                 - Previous peer in identifier ring (Not Implemented)")
			fmt.Println("    /ftable                      - List ftable peers (Not Implemented)")
			fmt.Println("    /create                      - Create new ring (Not Implemented)")
			fmt.Println("    /connect [onion address|id]  - Join to ring with peer (Not Implemented)")
			fmt.Println("    /lookup [id]                 - Find onion address of account with id (Not Implemented)")
			fmt.Println("    /ping [onion address|id]     - Ping peer (Not Implemented)")
			fmt.Println("    /ringcast [message]          - Message every peer in ring (Not Implemented)")
			fmt.Println("\n  DHT:")
			fmt.Println("    /put [key] [value]           - Put key and value into database (Not Implemented)")
			fmt.Println("    /get [key]                   - Get value of key (Not Implemented)")
			fmt.Println("    /delete [key]                - Delete value of key (Not Implemented)")
			fmt.Println("\n  ACCOUNT:")
			fmt.Println("    /accounts                    - List all accounts (Not Implemented)")
			fmt.Println("    /generate                    - Generate new account key pair (Not Implemented)")
			fmt.Println("    /delete                      - Delete an account key pair (Not Implemented)")
			fmt.Println("    /sign [id] [message]         - Sign with account key pair (Not Implemented)")
			fmt.Println("    /verify [id] [message]       - Verify a signed message with keypair (Not Implemented)")
			fmt.Println("    /encrypt [id] [message]      - Encrypt a message with keypair (Not Implemented)")
			fmt.Println("    /decrypt [id] [message]      - Decrypt a message with keypair (Not Implemented)")
			fmt.Println("\n  CONTACTS:")
			fmt.Println("    /contacts                    - List all saved contacts (Not Implemented)")
			fmt.Println("    /request [id] [message]      - Request account to add your id to their contacts (Not Implemented)")
			fmt.Println("    /add [id]                    - Add account to contacts (Not Implemented)")
			fmt.Println("    /rm [id]                     - Remove account from contacts (Not Implemented)")
			fmt.Println("    /whisper [id] [message]      - Direct message peer (Not Implemented)")
			fmt.Println("    /contactcast [message]       - Message all contacts (Not Implemented)")
			fmt.Println("\n  CHANNELS:")
			fmt.Println("    /channels                    - List all known channels (Not Implemented)")
			fmt.Println("    /channel                     - Generates a new channel (Not Implemented)")
			fmt.Println("    /join [id]                   - Join channel with id (Not Implemented)")
			fmt.Println("    /leave [id]                  - Leave channel with id (Not Implemented)")
			fmt.Println("    /channelcast [id] [message]  - Message all channel subscribers (Not Implemented)")
			fmt.Println("\n    /quit\n")
		} else if body == "/q" || body == "/quit" {
			oht.Interface.Stop()
		} else {
			oht.Interface.RingCast(username, body)
		}
		fmt.Printf(prompt)
	}
}
