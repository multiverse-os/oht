package network

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"time"
)

type P2PConfig struct {
	MaxPeers        int
	MaxPendingPeers int
	MaxQueueSize    int
}

type Server struct {
	Transport string
	OnionHost string
	Port      string
}

type EventFunc func(Manager *Manager, Peer *Peer)

type Manager struct {
	Config          *P2PConfig
	PrivateKey      *ecdsa.PrivateKey
	Servers         []*Server
	MaxPeers        int
	MaxPendingPeers int
	MaxQueueSize    int
	Peers           map[*Peer]bool
	Broadcast       chan Message
	Receive         chan Message
	Register        chan *Peer
	Unregister      chan *Peer
	OnConnect       EventFunc
	OnClose         EventFunc
	LastActivity    time.Time
}

func InitializeP2PManager(config *P2PConfig) *Manager {
	return &Manager{
		Config:          config,
		MaxPeers:        config.MaxPeers,
		MaxPendingPeers: config.MaxPendingPeers,
		MaxQueueSize:    config.MaxQueueSize,
		Broadcast:       make(chan Message, config.MaxQueueSize),
		Receive:         make(chan Message, config.MaxQueueSize),
		Register:        make(chan *Peer, config.MaxQueueSize),
		Unregister:      make(chan *Peer, config.MaxQueueSize),
		Peers:           make(map[*Peer]bool, config.MaxQueueSize),
		OnConnect:       nil,
		OnClose:         nil,
		LastActivity:    time.Now(),
	}
}

func (manager *Manager) Start() {
	for {
		select {
		case p := <-manager.Register:
			manager.Peers[p] = true
			log.Println("P2P: Peer connection established: ", p.Config.OnionHost)
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

//func (manager *Manager) DumpPeers() {
//	for p := range manager.Peers {
//		log.Println("Active Peers")
//	}
//}
