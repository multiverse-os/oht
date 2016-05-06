package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"net/http"
	"time"

	"../../types"
)

const (
	writeWait      = 15 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type EventFunc func(Manager *Manager, Peer *Peer)

type Manager struct {
	PrivateKey      *ecdsa.PrivateKey
	Server          *Server
	MaxPeers        int
	MaxPendingPeers int
	Peers           map[*Peer]bool
	Broadcast       chan types.Message
	Receive         chan types.Message
	Register        chan *Peer
	Unregister      chan *Peer
	OnConnect       EventFunc
	OnClose         EventFunc
	LastLookup      time.Time
}

func InitializeP2PManager(port string) *Manager {
	return &Manager{
		// Need to add a new message with custom type to send out
		// Info about current node, like onion address
		Server:          &Server{Port: port},
		MaxPeers:        8,
		MaxPendingPeers: 8,
		Broadcast:       make(chan types.Message, maxMessageSize),
		Receive:         make(chan types.Message, maxMessageSize),
		Register:        make(chan *Peer, maxMessageSize),
		Unregister:      make(chan *Peer, maxMessageSize),
		Peers:           make(map[*Peer]bool, maxMessageSize),
		OnConnect:       nil,
		OnClose:         nil,
		LastLookup:      time.Now(),
	}
}

func (manager *Manager) Start() {
	manager.Server.Start()
	for {
		select {
		case p := <-manager.Register:
			manager.Peers[p] = true
			log.Println("P2P: Peer connection established: ", p.OnionHost)
			fmt.Printf("oht> ")
			if manager.OnConnect != nil {
				go manager.OnConnect(manager, p)
			}
		case p := <-manager.Unregister:
			if _, ok := manager.Peers[p]; ok {
				delete(manager.Peers, p)
				close(p.Send)
				if manager.OnClose != nil {
					go manager.OnClose(manager, p)
				}
				p.WebSocket.Close()
			}
		case m := <-manager.Broadcast:
			for p := range manager.Peers {
				select {
				case p.Send <- m:
				default:
					close(p.Send)
					delete(manager.Peers, p)
				}
			}
		case m := <-manager.Receive:
			fmt.Println("")
			fmt.Println("[", m.Timestamp, "] ", m.Username, " : ", m.Body)
			fmt.Printf("oht> ")
		}
	}
}

func (manager *Manager) Stop() {
}

func (manager *Manager) DumpPeers() {
	for p := range manager.Peers {
		log.Println("Active Peers")
		log.Println("Connection: ", p.OnionHost)
	}
}

// Serve handles websocket requests from the peer
func (manager *Manager) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	p := &Peer{Send: make(chan types.Message, 256), WebSocket: ws}
	manager.Register <- p
	go p.writeMessages()
	p.readMessages()
}
