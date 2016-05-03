package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"./accounts"
	"./common"
	"./config"
	"./contacts"
	"./crypto"
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
	webUIPort   = flag.Int("wuiport", 8080, "Specify a webui port")
	socksPort   = flag.Int("socks", 12052, "Specify a socks proxy port")
	controlPort = flag.Int("control", 9555, "Specify a control port")
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
	// Unencrypted Account System Prototype For Low Security And Cases Where User Input Is Undesirable
	// This will be useful for assigning a key to the server struct under network for handshakes
	unencryptedKeyStore := crypto.NewKeyStorePlain(common.DefaultDataDir())
	unencryptedAccountManager := accounts.NewManager(unencryptedKeyStore)
	unencryptedAccount, _ := unencryptedAccountManager.NewAccount("password")
	log.Println("unencrypted account: " + unencryptedAccount.Address.Hex())
	// Encrypted Account System Prototype For Encryption And Signatures
	// This needs a secure password input, should build more fluid way to interact with accoutns
	encryptedKeyStore := crypto.NewKeyStorePassphrase(common.DefaultDataDir(), crypto.KDFStandard)
	encryptedAccountManager := accounts.NewManager(encryptedKeyStore)
	encryptedAccount, _ := encryptedAccountManager.NewAccount("password")
	log.Println("encrypted account:   " + encryptedAccount.Address.Hex())

	// Start P2P Networking
	go network.Manager.Start(tor.ListenPort)
	log.Printf("\nListening for peers (Websockets) :  " + tor.OnionHost)
	// Connect Directly To Peer (Will be required for bootstraping)
	if *peerAddress != "" {
		if match, _ := regexp.Match(":", []byte(*peerAddress)); !match {
			*peerAddress += ":12312"
		}
		log.Printf("Connecting to peer (Websockets)  :  " + *peerAddress)
		go network.ConnectToPeer(*peerAddress, tor.SocksPort)
	}
	// Start WebUI
	if *wui == true {
		webui.InitializeServer(tor.OnionWebUIHost, tor.WebUIPort)
		log.Printf("\nWeb UI :  " + tor.OnionWebUIHost + ":" + strconv.Itoa(tor.WebUIPort))
	}
	// Start console
	if *console == true {
		log.Println("\nWelcome to " + name + " console. Type \"/help\" to learn about the available commands.")
		prompt := "oht> "
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
			// This should be replaced with a better system but this works during early development
			if message.Body == "/help" {
				fmt.Println("COMMANDS:\n")
				fmt.Println("    /config               - List configuration values (Not Implemented)")
				fmt.Println("  DHT NETWORK:\n")
				fmt.Println("    /peers                - List all connected peers (Not Implemented)")
				fmt.Println("    /successor            - Next peer in identifier ring (Not Implemented)")
				fmt.Println("    /predecessor          - Previous peer in identifier ring (Not Implemented)")
				fmt.Println("    /ftable               - List ftable peers (Not Implemented)")
				fmt.Println("    /connect [ip address] - Direct connect to peer (Not Implemented)")
				fmt.Println("  ACCOUNTS:\n")
				fmt.Println("    /contacts             - List all saved contacts (Not Implemented)")
				fmt.Println("    /add     [user]       - Add account to contacts (Not Implemented)")
				fmt.Println("    /rm      [user]       - Remove account from contacts (Not Implemented)")
				fmt.Println("    /whisper [user]       - Direct message peer (Not Implemented)")
				fmt.Println("\n    /quit")
			} else if message.Body == "/q" || message.Body == "/quit" {
				os.Exit(0)
			} else {
				log.Println("[", message.Timestamp, "] ", message.Username, " : ", message.Body)
				network.Manager.Broadcast <- message
			}
			fmt.Printf(prompt)
		}
	}
}
