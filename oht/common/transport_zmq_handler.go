package dendrite

import (
	"github.com/golang/protobuf/proto"
)

func (transport *ZMQTransport) zmq_ping_handler(request *ChordMsg, w chan *ChordMsg) {
	pbPongMsg := &PBProtoPing{
		Version: proto.Int64(1),
	}
	pbPong, _ := proto.Marshal(pbPongMsg)
	pong := &ChordMsg{
		Type: PbPing,
		Data: pbPong,
	}
	w <- pong
}

func (transport *ZMQTransport) zmq_listVnodes_handler(request *ChordMsg, w chan *ChordMsg) {
	pblist := new(PBProtoListVnodesResp)
	for _, handler := range transport.table {
		h, _ := transport.getVnodeHandler(handler.vn)
		local_vn := h.(*localVnode)
		for _, vnode := range local_vn.ring.vnodes {
			pblist.Vnodes = append(pblist.Vnodes, vnode.ToProtobuf())
		}
		break
	}
	pbdata, err := proto.Marshal(pblist)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::ListVnodesHandler - failed to marshal response - " + err.Error())
		w <- errorMsg
		return
	}
	w <- &ChordMsg{
		Type: PbListVnodesResp,
		Data: pbdata,
	}
	return
}

func (transport *ZMQTransport) zmq_find_successors_handler(request *ChordMsg, w chan *ChordMsg) {
	pbMsg := request.TransportMsg.(PBProtoFindSuccessors)
	key := pbMsg.GetKey()
	dest := VnodeFromProtobuf(pbMsg.GetDest())

	// make sure destination vnode exists locally
	local_vn, err := transport.getVnodeHandler(dest)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::FindSuccessorsHandler - " + err.Error())
		w <- errorMsg
		return
	}
	succs, forward_vn, err := local_vn.FindSuccessors(key, int(pbMsg.GetLimit()))
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::FindSuccessorsHandler - " + err.Error())
		w <- errorMsg
		return
	}

	// if forward_vn is not set, return the list
	if forward_vn == nil {
		pblist := new(PBProtoListVnodesResp)
		for _, s := range succs {
			pblist.Vnodes = append(pblist.Vnodes, s.ToProtobuf())
		}
		pbdata, err := proto.Marshal(pblist)
		if err != nil {
			errorMsg := transport.newErrorMsg("ZMQ::FindSuccessorsHandler - failed to marshal response - " + err.Error())
			w <- errorMsg
			return
		}
		w <- &ChordMsg{
			Type: PbListVnodesResp,
			Data: pbdata,
		}
		return
	}
	// send forward response
	pbfwd := &PBProtoForward{Vnode: forward_vn.ToProtobuf()}
	pbdata, err := proto.Marshal(pbfwd)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::FindSuccessorsHandler - failed to marshal forward response - " + err.Error())
		w <- errorMsg
		return
	}
	w <- &ChordMsg{
		Type: PbForward,
		Data: pbdata,
	}
}

func (transport *ZMQTransport) zmq_get_predecessor_handler(request *ChordMsg, w chan *ChordMsg) {
	pbMsg := request.TransportMsg.(PBProtoGetPredecessor)
	dest := VnodeFromProtobuf(pbMsg.GetDest())

	// make sure destination vnode exists locally
	local_vn, err := transport.getVnodeHandler(dest)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::GetPredecessorHandler - " + err.Error())
		w <- errorMsg
		return
	}

	pred, err := local_vn.GetPredecessor()
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::GetPredecessorHandler - " + err.Error())
		w <- errorMsg
		return
	}
	pbpred := &PBProtoVnode{}
	if pred != nil {
		pbpred.Id = pred.Id
		pbpred.Host = proto.String(pred.Host)
	}
	pbdata, err := proto.Marshal(pbpred)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::GetPredecessorHandler - Failed to marshal response - " + err.Error())
		w <- errorMsg
		return
	}

	w <- &ChordMsg{
		Type: PbProtoVnode,
		Data: pbdata,
	}

}

// handle Notify() request
func (transport *ZMQTransport) zmq_notify_handler(request *ChordMsg, w chan *ChordMsg) {
	pbMsg := request.TransportMsg.(PBProtoNotify)
	dest := VnodeFromProtobuf(pbMsg.GetDest())

	// make sure destination vnode exists locally
	local_vn, err := transport.getVnodeHandler(dest)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::NotifyHandler - " + err.Error())
		w <- errorMsg
		return
	}
	pred := VnodeFromProtobuf(pbMsg.GetVnode())
	succ_list, err := local_vn.Notify(pred)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::NotifyHandler - " + err.Error())
		w <- errorMsg
		return
	}
	pblist := new(PBProtoListVnodesResp)
	for _, succ := range succ_list {
		if succ == nil {
			break
		}
		pblist.Vnodes = append(pblist.Vnodes, succ.ToProtobuf())
	}

	pbdata, err := proto.Marshal(pblist)
	if err != nil {
		errorMsg := transport.newErrorMsg("ZMQ::Notify - Failed to marshal response - " + err.Error())
		w <- errorMsg
		return
	}
	w <- &ChordMsg{
		Type: PbListVnodesResp,
		Data: pbdata,
	}
}

func (transport *ZMQTransport) zmq_leave_handler(request *ChordMsg, w chan *ChordMsg) {

}
func (transport *ZMQTransport) zmq_error_handler(request *ChordMsg, w chan *ChordMsg) {

}
