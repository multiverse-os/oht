package network

import (
	"crypto/ecdsa"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

const (
	defaultConnectionTimeout = 15 * time.Second
	refreshPeersInterval     = 90 * time.Second
	maxConnections           = 8
)

type Server struct {
	PrivateKey *ecdsa.PrivateKey
	Websocket  *gin.Engine
	Port       int
	lock       sync.Mutex
}

func (server *Server) Start() {
	server.lock.Lock()
	defer server.lock.Unlock()
	server.startListening()
}

func (server *Server) startListening() {
	gin.SetMode(gin.ReleaseMode)
	server.Websocket = gin.Default()
	server.Websocket.GET("/ws", func(c *gin.Context) {
		Manager.Serve(c.Writer, c.Request)
	})
	go server.Websocket.Run("127.0.0.1:" + strconv.Itoa(server.Port))
}

func (server *Server) PeerCount() int {
	return 0
}
