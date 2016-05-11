package dendrite

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"time"
)

const (
	// protocol buffer messages (for definitions, see pb_defs/chord.proto)
	PbPing MsgType = iota
	PbAck
	PbErr
	PbForward
	PbJoin
	PbLeave
	PbListVnodes
	PbListVnodesResp
	PbFindSuccessors
	PbGetPredecessor
	PbProtoVnode
	PbNotify
)

// newErrorMsg is a helper to create encoded *ChordMsg (PBProtoErr) with error in it.
func (transport *ZMQTransport) newErrorMsg(msg string) *ChordMsg {
	pbmsg := &PBProtoErr{
		Error: proto.String(msg),
	}
	pbdata, _ := proto.Marshal(pbmsg)
	return &ChordMsg{
		Type: PbErr,
		Data: pbdata,
	}

}

// NewErrorMsg is a helper to create encoded *ChordMsg (PBProtoErr) with error in it.
func (transport *ZMQTransport) NewErrorMsg(msg string) *ChordMsg {
	pbmsg := &PBProtoErr{
		Error: proto.String(msg),
	}
	pbdata, _ := proto.Marshal(pbmsg)
	return &ChordMsg{
		Type: PbErr,
		Data: pbdata,
	}

}

// Encode implement's Transport's Encode() in ZMQTransport.
func (transport *ZMQTransport) Encode(mt MsgType, data []byte) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(byte(mt))
	buf.Write(data)
	return buf.Bytes()
}

// Decode implements Transport's Decode() in ZMQTransport. For request messages
// it also sets their respective handler to be called when such request comes in.
// If message type is unknown to this transport, Decode() also checks for registered
// TransportHooks and runs their Decode() implementation.
func (transport *ZMQTransport) Decode(data []byte) (*ChordMsg, error) {
	data_len := len(data)
	if data_len == 0 {
		return nil, fmt.Errorf("data too short: %d", len(data))
	}

	cm := &ChordMsg{Type: MsgType(data[0])}

	if data_len > 1 {
		cm.Data = data[1:]
	}

	// parse the data and set the handler
	switch cm.Type {
	case PbPing:
		var pingMsg PBProtoPing
		err := proto.Unmarshal(cm.Data, &pingMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoPing message - %s", err)
		}
		cm.TransportMsg = pingMsg
		cm.TransportHandler = transport.zmq_ping_handler
	case PbErr:
		var errorMsg PBProtoErr
		err := proto.Unmarshal(cm.Data, &errorMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoErr message - %s", err)
		}
		cm.TransportMsg = errorMsg
		cm.TransportHandler = transport.zmq_error_handler
	case PbForward:
		var forwardMsg PBProtoForward
		err := proto.Unmarshal(cm.Data, &forwardMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoForward message - %s", err)
		}
		cm.TransportMsg = forwardMsg
	case PbLeave:
		var leaveMsg PBProtoLeave
		err := proto.Unmarshal(cm.Data, &leaveMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoLeave message - %s", err)
		}
		cm.TransportMsg = leaveMsg
		cm.TransportHandler = transport.zmq_leave_handler
	case PbListVnodes:
		var listVnodesMsg PBProtoListVnodes
		err := proto.Unmarshal(cm.Data, &listVnodesMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoListVnodes message - %s", err)
		}
		cm.TransportMsg = listVnodesMsg
		cm.TransportHandler = transport.zmq_listVnodes_handler
	case PbListVnodesResp:
		var listVnodesRespMsg PBProtoListVnodesResp
		err := proto.Unmarshal(cm.Data, &listVnodesRespMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoListVnodesResp message - %s", err)
		}
		cm.TransportMsg = listVnodesRespMsg
	case PbFindSuccessors:
		var findSuccMsg PBProtoFindSuccessors
		err := proto.Unmarshal(cm.Data, &findSuccMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoFindSuccessors message - %s", err)
		}
		cm.TransportMsg = findSuccMsg
		cm.TransportHandler = transport.zmq_find_successors_handler
	case PbGetPredecessor:
		var getPredMsg PBProtoGetPredecessor
		err := proto.Unmarshal(cm.Data, &getPredMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoGetPredecessor message - %s", err)
		}
		cm.TransportMsg = getPredMsg
		cm.TransportHandler = transport.zmq_get_predecessor_handler
	case PbNotify:
		var notifyMsg PBProtoNotify
		err := proto.Unmarshal(cm.Data, &notifyMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoNotify message - %s", err)
		}
		cm.TransportMsg = notifyMsg
		cm.TransportHandler = transport.zmq_notify_handler
	case PbProtoVnode:
		var vnodeMsg PBProtoVnode
		err := proto.Unmarshal(cm.Data, &vnodeMsg)
		if err != nil {
			return nil, fmt.Errorf("error decoding PBProtoVnode message - %s", err)
		}
		cm.TransportMsg = vnodeMsg
	default:
		// maybe a TransportHook should handle this?
		for _, hook := range transport.hooks {
			if hook_cm, err := hook.Decode(data); err != nil {
				_, ok := err.(ErrHookUnknownType)
				if ok {
					// this hook knows nothing about this message type, try next one
					continue
				}
				return nil, err
			} else {
				// hook is handling this!
				return hook_cm, nil
			}
		}
		return nil, fmt.Errorf("error decoding message - unknown request type %x", cm.Type)
	}

	return cm, nil
}

// getVnodeHandler returns registered local vnode handler, if one is found for given vnode.
func (transport *ZMQTransport) getVnodeHandler(dest *Vnode) (VnodeHandler, error) {
	h, ok := transport.table[dest.String()]
	if ok {
		return h.handler, nil
	}
	return nil, fmt.Errorf("local vnode handler not found")
}

// GetVnodeHandler returns registered local vnode handler, if one is found for given vnode.
func (transport *ZMQTransport) GetVnodeHandler(vnode *Vnode) (VnodeHandler, bool) {
	handler, err := transport.getVnodeHandler(vnode)
	if err != nil {
		return nil, false
	}
	return handler, true
}

// Register registers a VnodeHandler within ZMQTransport.
func (transport *ZMQTransport) Register(vnode *Vnode, handler VnodeHandler) {
	transport.lock.Lock()
	transport.table[vnode.String()] = &localHandler{vn: vnode, handler: handler}
	transport.lock.Unlock()
}

// ListVnodes - client request. Implements Transport's ListVnodes() in ZQMTransport.
func (transport *ZMQTransport) ListVnodes(host string) ([]*Vnode, error) {
	error_c := make(chan error, 1)
	resp_c := make(chan []*Vnode, 1)

	go func() {
		req_sock, err := transport.zmq_context.NewSocket(zmq.REQ)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:ListVnodes - newsocket error - %s", err)
			return
		}
		req_sock.SetRcvtimeo(2 * time.Second)
		req_sock.SetSndtimeo(2 * time.Second)

		defer req_sock.Close()
		err = req_sock.Connect("tcp://" + host)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:ListVnodes - connect error - %s", err)
			return
		}
		// Build request protobuf
		req := new(PBProtoListVnodes)
		reqData, _ := proto.Marshal(req)
		encoded := transport.Encode(PbListVnodes, reqData)
		_, err = req_sock.SendBytes(encoded, 0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::ListVnodes - error while sending request - %s", err)
			return
		}

		// read response and decode it
		resp, err := req_sock.RecvBytes(0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::ListVnodes - error while reading response - %s", err)
			return
		}
		decoded, err := transport.Decode(resp)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::ListVnodes - error while decoding response - %s", err)
			return
		}

		switch decoded.Type {
		case PbErr:
			pbMsg := decoded.TransportMsg.(PBProtoErr)
			error_c <- fmt.Errorf("ZMQ::ListVnodes - got error response - %s", pbMsg.GetError())
		case PbListVnodesResp:
			pbMsg := decoded.TransportMsg.(PBProtoListVnodesResp)
			vnodes := make([]*Vnode, len(pbMsg.GetVnodes()))
			for idx, pbVnode := range pbMsg.GetVnodes() {
				vnodes[idx] = VnodeFromProtobuf(pbVnode)
			}
			resp_c <- vnodes
			return
		default:
			// unexpected response
			error_c <- fmt.Errorf("ZMQ::ListVnodes - unexpected response")
			return
		}
	}()

	select {
	case <-time.After(transport.clientTimeout):
		return nil, fmt.Errorf("ZMQ::ListVnodes - command timed out!")
	case err := <-error_c:
		return nil, err
	case resp_vnodes := <-resp_c:
		return resp_vnodes, nil
	}

}

// FindSuccessors - client request. Implements Transport's FindSuccessors() in ZQMTransport.
func (transport *ZMQTransport) FindSuccessors(remote *Vnode, limit int, key []byte) ([]*Vnode, error) {
	error_c := make(chan error, 1)
	resp_c := make(chan []*Vnode, 1)
	forward_c := make(chan *Vnode, 1)

	go func() {
		req_sock, err := transport.zmq_context.NewSocket(zmq.REQ)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:FindSuccessors - newsocket error - %s", err)
			return
		}
		req_sock.SetRcvtimeo(2 * time.Second)
		req_sock.SetSndtimeo(2 * time.Second)

		defer req_sock.Close()
		err = req_sock.Connect("tcp://" + remote.Host)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:FindSuccessors - connect error - %s", err)
			return
		}
		// Build request protobuf
		req := &PBProtoFindSuccessors{
			Dest:  remote.ToProtobuf(),
			Key:   key,
			Limit: proto.Int32(int32(limit)),
		}
		reqData, _ := proto.Marshal(req)
		encoded := transport.Encode(PbFindSuccessors, reqData)
		_, err = req_sock.SendBytes(encoded, 0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::FindSuccessors - error while sending request - %s", err)
			return
		}

		// read response and decode it
		resp, err := req_sock.RecvBytes(0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::FindSuccessors - error while reading response - %s %X", err, remote.Id)
			return
		}
		decoded, err := transport.Decode(resp)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::FindSuccessors - error while decoding response - %s", err)
			return
		}

		switch decoded.Type {
		case PbErr:
			pbMsg := decoded.TransportMsg.(PBProtoErr)
			error_c <- fmt.Errorf("ZMQ::FindSuccessors - got error response - %s", pbMsg.GetError())
			return
		case PbForward:
			pbMsg := decoded.TransportMsg.(PBProtoForward)
			forward_c <- VnodeFromProtobuf(pbMsg.GetVnode())
			return
		case PbListVnodesResp:
			pbMsg := decoded.TransportMsg.(PBProtoListVnodesResp)
			vnodes := make([]*Vnode, len(pbMsg.GetVnodes()))
			for idx, pbVnode := range pbMsg.GetVnodes() {
				vnodes[idx] = VnodeFromProtobuf(pbVnode)
			}
			resp_c <- vnodes
			return
		default:
			// unexpected response
			error_c <- fmt.Errorf("ZMQ::FindSuccessors - unexpected response")
			return
		}
	}()

	select {
	case <-time.After(transport.clientTimeout):
		return nil, fmt.Errorf("ZMQ::FindSuccessors - command timed out!")
	case err := <-error_c:
		return nil, err
	case new_remote := <-forward_c:
		return transport.FindSuccessors(new_remote, limit, key)
	case resp_vnodes := <-resp_c:
		return resp_vnodes, nil
	}
}

// GetPredecessor - client request. Implements Transport's GetPredecessor() in ZQMTransport.
func (transport *ZMQTransport) GetPredecessor(remote *Vnode) (*Vnode, error) {
	error_c := make(chan error, 1)
	resp_c := make(chan *Vnode, 1)

	go func() {
		req_sock, err := transport.zmq_context.NewSocket(zmq.REQ)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:GetPredecessor - newsocket error - %s", err)
			return
		}
		req_sock.SetRcvtimeo(2 * time.Second)
		req_sock.SetSndtimeo(2 * time.Second)

		defer req_sock.Close()
		err = req_sock.Connect("tcp://" + remote.Host)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:GetPredecessor - connect error - %s", err)
			return
		}
		// Build request protobuf
		req := &PBProtoGetPredecessor{
			Dest: remote.ToProtobuf(),
		}
		reqData, _ := proto.Marshal(req)
		encoded := transport.Encode(PbGetPredecessor, reqData)
		_, err = req_sock.SendBytes(encoded, 0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:GetPredecessor - error while sending request - %s", err)
			return
		}

		// read response and decode it
		resp, err := req_sock.RecvBytes(0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::GetPredecessor - error while reading response - %s", err)
			return
		}
		decoded, err := transport.Decode(resp)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::GetPredecessor - error while decoding response - %s", err)
			return
		}

		switch decoded.Type {
		case PbErr:
			pbMsg := decoded.TransportMsg.(PBProtoErr)
			error_c <- fmt.Errorf("ZMQ::GetPredecessor - got error response - %s", pbMsg.GetError())
			return
		case PbProtoVnode:
			pbMsg := decoded.TransportMsg.(PBProtoVnode)
			resp_c <- VnodeFromProtobuf(&pbMsg)
			return
		default:
			// unexpected response
			error_c <- fmt.Errorf("ZMQ::GetPredecessor - unexpected response")
			return
		}
	}()

	select {
	case <-time.After(transport.clientTimeout):
		return nil, fmt.Errorf("ZMQ::GetPredecessor - command timed out!")
	case err := <-error_c:
		return nil, err
	case resp_vnode := <-resp_c:
		return resp_vnode, nil
	}
}

// Notify - client request. Implements Transport's Notify() in ZQMTransport.
func (transport *ZMQTransport) Notify(remote, self *Vnode) ([]*Vnode, error) {
	error_c := make(chan error, 1)
	resp_c := make(chan []*Vnode, 1)

	go func() {
		req_sock, err := transport.zmq_context.NewSocket(zmq.REQ)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:GetPredecessor - newsocket error - %s", err)
			return
		}
		req_sock.SetRcvtimeo(2 * time.Second)
		req_sock.SetSndtimeo(2 * time.Second)

		defer req_sock.Close()
		err = req_sock.Connect("tcp://" + remote.Host)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ:GetPredecessor - connect error - %s", err)
			return
		}

		// Build request protobuf
		req := &PBProtoNotify{
			Dest:  remote.ToProtobuf(),
			Vnode: self.ToProtobuf(),
		}
		reqData, _ := proto.Marshal(req)
		encoded := transport.Encode(PbNotify, reqData)
		_, err = req_sock.SendBytes(encoded, 0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::Notify - error while sending request - %s", err)
			return
		}

		// read response and decode it
		resp, err := req_sock.RecvBytes(0)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::Notify - error while reading response - %s", err)
			return
		}
		decoded, err := transport.Decode(resp)
		if err != nil {
			error_c <- fmt.Errorf("ZMQ::Notify - error while decoding response - %s", err)
			return
		}

		switch decoded.Type {
		case PbErr:
			pbMsg := decoded.TransportMsg.(PBProtoErr)
			error_c <- fmt.Errorf("ZMQ::Notify - got error response - %s", pbMsg.GetError())
			return
		case PbListVnodesResp:
			pbMsg := decoded.TransportMsg.(PBProtoListVnodesResp)
			vnodes := make([]*Vnode, len(pbMsg.GetVnodes()))
			for idx, pbVnode := range pbMsg.GetVnodes() {
				vnodes[idx] = VnodeFromProtobuf(pbVnode)
			}
			resp_c <- vnodes
			return
		default:
			// unexpected response
			error_c <- fmt.Errorf("ZMQ::Notify - unexpected response")
			return
		}
	}()

	select {
	case <-time.After(transport.clientTimeout):
		return nil, fmt.Errorf("ZMQ::Notify - command timed out!")
	case err := <-error_c:
		return nil, err
	case resp_vnode := <-resp_c:
		return resp_vnode, nil
	}
}

// Ping - client request. Implements Transport's Ping() in ZQMTransport.
func (transport *ZMQTransport) Ping(remote_vn *Vnode) (bool, error) {
	req_sock, err := transport.zmq_context.NewSocket(zmq.REQ)
	if err != nil {
		return false, err
	}
	defer req_sock.Close()

	err = req_sock.Connect("tcp://" + remote_vn.Host)
	if err != nil {
		return false, err
	}
	req_sock.SetRcvtimeo(2 * time.Second)
	req_sock.SetSndtimeo(2 * time.Second)

	PbPingMsg := &PBProtoPing{
		Version: proto.Int64(1),
	}
	PbPingData, _ := proto.Marshal(PbPingMsg)
	encoded := transport.Encode(PbPing, PbPingData)
	_, err = req_sock.SendBytes(encoded, 0)
	if err != nil {
		return false, err
	}
	resp, err := req_sock.RecvBytes(0)
	if err != nil {
		return false, err
	}
	decoded, err := transport.Decode(resp)
	if err != nil {
		return false, err
	}
	pongMsg := new(PBProtoPing)
	err = proto.Unmarshal(decoded.Data, pongMsg)
	if err != nil {
		return false, err
	}
	return true, nil
}
