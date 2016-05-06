package network

import (
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var stoppedError = errors.New("Gin: Webserver is being stopped")

type stoppableListener struct {
	tcpKeepAliveListener
	stop chan int
}

type WebServer struct {
	Online   bool
	Host     string
	engine   *gin.Engine
	listener *stoppableListener
}

func InitializeWebServer(r *gin.Engine, host string) (server *WebServer) {
	return &WebServer{
		Online: false,
		Host:   host,
		engine: r,
	}
}

func (wServer *WebServer) Start() error {
	if !wServer.Online {
		hServer := &http.Server{Addr: wServer.Host, Handler: wServer.engine}
		listener, err := net.Listen("tcp", wServer.Host)
		if err != nil {
			return err
		}
		wServer.listener, err = newStoppableListener(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			return err
		}
		go hServer.Serve(wServer.listener)
		wServer.Online = true
		if err != nil {
			if err != stoppedError {
				panic(err)
			}
		}
		return nil
	} else {
		return errors.New("Gin: Web server is already online")
	}
}

func (server *WebServer) Stop() bool {
	close(server.listener.stop)
	server.Online = false
	return true
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func newStoppableListener(l net.Listener) (*stoppableListener, error) {
	tcpL, ok := l.(tcpKeepAliveListener)
	if !ok {
		return nil, errors.New("Gin: Cannot wrap listener")
	}
	retval := &stoppableListener{}
	retval.tcpKeepAliveListener = tcpL
	retval.stop = make(chan int)
	return retval, nil
}

func (sl *stoppableListener) Accept() (net.Conn, error) {
	for {
		sl.SetDeadline(time.Now().Add(time.Second))
		newConn, err := sl.tcpKeepAliveListener.Accept()
		select {
		case <-sl.stop:
			return nil, stoppedError
		default:
		}
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}
		return newConn, err
	}
}
