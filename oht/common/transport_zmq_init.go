package dendrite

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"os"
	"sync"
	"time"
)

type controlType int

const (
	workerShutdownReq controlType = iota
	workerShutdownAllowed
	workerShutdownDenied
	workerShutdownConfirm
	workerRegisterReq
	workerRegisterAllowed
	workerRegisterDenied
	workerCtlShutdown
)

// ZMQTransport implements Transport interface using ZeroMQ for communication.
type ZMQTransport struct {
	lock              *sync.Mutex
	minHandlers       int
	maxHandlers       int
	incrHandlers      int
	activeRequests    int
	ring              *Ring
	table             map[string]*localHandler
	clientTimeout     time.Duration
	ClientTimeout     time.Duration
	control_c         chan *workerComm
	dealer_sock       *zmq.Socket
	router_sock       *zmq.Socket
	zmq_context       *zmq.Context
	ZMQContext        *zmq.Context
	workerIdleTimeout time.Duration
	hooks             []TransportHook
	Logger            *log.Logger
}

// RegisterHook registers TransportHook within ZMQTransport.
func (t *ZMQTransport) RegisterHook(h TransportHook) {
	t.hooks = append(t.hooks, h)
}

/*
	InitZMQTransport creates ZeroMQ transport.

	It multiplexes incoming connections which are then processed in separate go routines (workers).
	Multiplexer spawns go routines as needed, but 10 worker routines are created on startup.
	Every request times out after provided timeout duration. ZMQ pattern is:
		zmq.ROUTER(incoming) -> proxy -> zmq.DEALER -> [zmq.REP(worker), zmq.REP...]
*/
func InitZMQTransport(hostname string, timeout time.Duration, logger *log.Logger) (Transport, error) {
	// use default logger if one is not provided
	if logger == nil {
		logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	// initialize ZMQ Context
	context, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}

	// setup router and bind() to tcp address for clients to connect to
	router_sock, err := context.NewSocket(zmq.ROUTER)
	if err != nil {
		return nil, err
	}
	err = router_sock.Bind("tcp://" + hostname)
	if err != nil {
		return nil, err
	}

	// setup dealer
	dealer_sock, err := context.NewSocket(zmq.DEALER)
	if err != nil {
		return nil, err
	}
	err = dealer_sock.Bind("inproc://dendrite-zmqdealer")
	if err != nil {
		return nil, err
	}
	poller := zmq.NewPoller()
	poller.Add(router_sock, zmq.POLLIN)
	poller.Add(dealer_sock, zmq.POLLIN)

	transport := &ZMQTransport{
		lock:              new(sync.Mutex),
		clientTimeout:     timeout,
		ClientTimeout:     timeout,
		minHandlers:       10,
		maxHandlers:       1024,
		incrHandlers:      10,
		activeRequests:    0,
		workerIdleTimeout: 10 * time.Second,
		table:             make(map[string]*localHandler),
		control_c:         make(chan *workerComm),
		dealer_sock:       dealer_sock,
		router_sock:       router_sock,
		zmq_context:       context,
		ZMQContext:        context,
		hooks:             make([]TransportHook, 0),
		Logger:            logger,
	}

	go zmq.Proxy(router_sock, dealer_sock, nil)
	// Scheduler goroutine keeps track of running workers
	// It spawns new ones if needed, and cancels ones that are idling
	go func() {
		sched_ticker := time.NewTicker(60 * time.Second)
		workers := make(map[*workerComm]bool)
		// fire up initial set of workers
		for i := 0; i < transport.minHandlers; i++ {
			go transport.zmq_worker()
		}
		for {
			select {
			case comm := <-transport.control_c:
				// worker sent something...
				msg := <-comm.worker_out
				switch {
				case msg == workerRegisterReq:
					if len(workers) == transport.maxHandlers {
						comm.worker_in <- workerRegisterDenied
						logger.Println("[DENDRITE][INFO]: TransportListener - max number of workers reached")
						continue
					}
					if _, ok := workers[comm]; ok {
						// worker already registered
						continue
					}
					comm.worker_in <- workerRegisterAllowed
					workers[comm] = true
					logger.Println("[DENDRITE][INFO]: TransportListener - registered new worker, total:", len(workers))
				case msg == workerShutdownReq:
					//logger.Println("Got shutdown req")
					if len(workers) > transport.minHandlers {
						comm.worker_in <- workerShutdownAllowed
						for _ = range comm.worker_out {
							// wait until worker closes the channel
						}
						delete(workers, comm)
					} else {
						comm.worker_in <- workerShutdownDenied
					}
				}
			case <-sched_ticker.C:
				// check if requests are piling up and start more workers if that's the case
				if transport.activeRequests > 3*len(workers) {
					for i := 0; i < transport.incrHandlers; i++ {
						go transport.zmq_worker()
					}
				}
			}
		}
	}()
	return transport, nil
}

type workerComm struct {
	worker_in  chan controlType // worker's input channel for two way communication with scheduler
	worker_out chan controlType // worker's output channel for two way communication with scheduler
	worker_ctl chan controlType // worker's control channel for communication with scheduler
}

func (transport *ZMQTransport) zmq_worker() {
	// setup REP socket
	rep_sock, err := transport.zmq_context.NewSocket(zmq.REP)
	if err != nil {
		transport.Logger.Println("[DENDRITE][ERROR]: TransportListener worker failed to create REP socket", err)
		return
	}
	err = rep_sock.Connect("inproc://dendrite-zmqdealer")
	if err != nil {
		transport.Logger.Println("[DENDRITE][ERROR]: TransportListener worker failed to connect to dealer", err)
		return
	}

	// setup communication channels with scheduler
	worker_in := make(chan controlType, 1)
	worker_out := make(chan controlType, 1)
	worker_ctl := make(chan controlType, 1)
	comm := &workerComm{
		worker_in:  worker_in,
		worker_out: worker_out,
		worker_ctl: worker_ctl,
	}
	// notify scheduler that we're up
	worker_out <- workerRegisterReq
	transport.control_c <- comm
	v := <-worker_in
	if v == workerRegisterDenied {
		return
	}

	// setup socket read channel
	rpc_req_c := make(chan *ChordMsg)
	rpc_response_c := make(chan *ChordMsg)
	poller := zmq.NewPoller()
	poller.Add(rep_sock, zmq.POLLIN)
	cancel_c := make(chan bool, 1)
	// read from socket and emit data, or stop if canceled
	go func() {
	MAINLOOP:
		for {
			// poll for 5 seconds, but then see if we should be canceled
			sockets, _ := poller.Poll(5 * time.Second)
			for _, socket := range sockets {
				rawmsg, err := socket.Socket.RecvBytes(0)
				if err != nil {
					transport.Logger.Println("[DENDRITE][ERROR]: TransportListener error while reading from REP, ", err)
					continue
				}
				// decode raw data
				decoded, err := transport.Decode(rawmsg)
				if err != nil {
					errorMsg := transport.newErrorMsg("Failed to decode request - " + err.Error())
					encoded := transport.Encode(errorMsg.Type, errorMsg.Data)
					socket.Socket.SendBytes(encoded, 0)
					continue
				}
				// if transportHandler is nil, this can not be a valid request
				if decoded.TransportHandler == nil {
					errorMsg := transport.newErrorMsg("Invalid request, unknown handler")
					encoded := transport.Encode(errorMsg.Type, errorMsg.Data)
					socket.Socket.SendBytes(encoded, 0)
					continue
				}
				rpc_req_c <- decoded
				// wait for response
				response := <-rpc_response_c
				encoded := transport.Encode(response.Type, response.Data)
				socket.Socket.SendBytes(encoded, 0)
			}
			// check for cancel request
			select {
			case <-cancel_c:
				break MAINLOOP
			default:
				break
			}
		}
	}()

	// read from socket and process request -- OR
	// shutdown if scheduler wants me to -- OR
	// request shutdown from scheduler if idling and exit if allowed
	ticker := time.NewTicker(transport.workerIdleTimeout)
	for {
		select {
		case request := <-rpc_req_c:
			// handle request
			request.TransportHandler(request, rpc_response_c)
			// restart idle timer
			ticker.Stop()
			ticker = time.NewTicker(transport.workerIdleTimeout)

		case controlMsg := <-comm.worker_ctl:
			if controlMsg == workerCtlShutdown {
				close(comm.worker_out)
				cancel_c <- true
				close(cancel_c)
				transport.Logger.Println("[DENDRITE][INFO]: TransportListener: worker shutdown")
				return
			}
		case <-ticker.C:
			// we're idling, lets request shutdown
			comm.worker_out <- workerShutdownReq
			transport.control_c <- comm
			v := <-comm.worker_in
			if v == workerShutdownAllowed {
				transport.Logger.Println("[DENDRITE][INFO]: TransportListener: worker shutdown due to idle state")
				close(comm.worker_out)
				cancel_c <- true
				close(cancel_c)
				return
			}
		}
	}
}
