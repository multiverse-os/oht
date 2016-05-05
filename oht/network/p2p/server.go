package p2p

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultConnectionTimeout = 15 * time.Second
	refreshPeersInterval     = 90 * time.Second
	maxConnections           = 8
)

type Server struct {
	WebSocket *gin.Engine
	lock      sync.Mutex
	Port      string
	Manager   *Manager
}

func (server *Server) Start() {
	server.lock.Lock()
	defer server.lock.Unlock()
	gin.SetMode(gin.ReleaseMode)
	server.WebSocket = gin.Default()
	server.WebSocket.GET("/ws", func(c *gin.Context) {
		server.Manager.Serve(c.Writer, c.Request)
	})
	go server.WebSocket.Run("127.0.0.1:" + server.Port)
}

func (server *Server) PeerCount() int {
	return 0
}
