package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"./common"
	"./config"
	"./contacts"
	"./database"
	"./network"
	"./webui"

	"github.com/pborman/uuid"
)

var (
	name          = "oht"
	version_major = 0
	version_minor = 1
	version_patch = 0
	version       = fmt.Sprintf("%d.%d.%d", version_major, version_minor, version_patch)
	//daemon        = flag.Bool("d", true, "Start the process daemonized")
	console     = flag.Bool("c", true, "Start the process with a console")
	wui         = flag.Bool("wui", true, "Start the process with a web ui")
	username    = flag.String("username", "user", "Specify a username")
	peerAddress = flag.String("peer", "", "Specify a peer address for direct connection")
	listenPort  = flag.Int("listen", 12312, "Specify a listen port")
	webUIPort   = flag.String("wuiport", "8080", "Specify a webui port")
	socksPort   = flag.String("socks", "12052", "Specify a socks proxy port")
	controlPort = flag.String("control", "9555", "Specify a control port")
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	// Initialize Data Directory
	if !common.FileExist(common.DefaultDataDir()) {
		os.MkdirAll(common.DefaultDataDir(), os.ModePerm)
	}
	config.InitializeConfig()
	contacts.InitializeContacts()
	// Start Tor
	log.Println("Starting " + common.MakeName(name, version) + ":")
	log.Println("########################################")
	log.Println("Initializing Tor Process...")
	tor := network.InitializeTor(*listenPort, *socksPort, *controlPort, *webUIPort)
	// Database
	database.InitializeDatabase()
	// Define a clean shutdown process
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		tor.Shutdown()
		os.Exit(1)
	}()
	// Start P2P Networking
	go network.Manager.Start(tor.ListenPort)
	log.Printf("\nListening for peers :  " + tor.OnionHost)
	// Connect Directly To Peer (Will be required for bootstraping)
	if *peerAddress != "" {
		if match, _ := regexp.Match(":", []byte(*peerAddress)); !match {
			*peerAddress += ":12312"
		}
		log.Printf("Connecting to peer  :  " + *peerAddress)
		go network.ConnectToPeer(*peerAddress, tor.SocksPort)
	}
	// Start WebUI
	if *wui == true {
		webui.InitializeServer(tor.WebUIPort, tor.OnionWebUIHost)
		log.Printf("\nWeb UI :  " + tor.OnionWebUIHost + ":" + tor.WebUIPort)
	}
	// Start console
	if *console == true {
		log.Println("\nWelcome to " + name + " console. Type \"/help\" to learn about the available commands.")
		prompt := "whisper> "
		// Going to need a function to dump all the peers
		cli := bufio.NewScanner(os.Stdin)
		fmt.Printf(prompt)
		for cli.Scan() {
			message := network.Message{
				Id:        uuid.New(),
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
				Username:  *username,
				Body:      cli.Text()}
			// Check for commands
			if message.Body == "/help" {
				fmt.Println("Available Commands:\n")
				fmt.Println("    /peers - List all active peers (Not Implemented)")
				fmt.Println("    /whisper - Direct message peer (Not Implemented)")
				fmt.Println("    /connect - Direct connect to peer (Not Implemented)")
			} else {
				log.Println("[", message.Timestamp, "] ", message.Username, " : ", message.Body)
				network.Manager.Broadcast <- message
			}
			fmt.Printf(prompt)
		}
	}
}
