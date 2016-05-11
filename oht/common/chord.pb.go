package dendrite

import proto "github.com/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

// PBProtoVnode represents Vnode structure.
type PBProtoVnode struct {
	Id               []byte  `protobuf:"bytes,1,req,name=id" json:"id,omitempty"`
	Host             *string `protobuf:"bytes,2,req,name=host" json:"host,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PBProtoVnode) Reset()         { *m = PBProtoVnode{} }
func (m *PBProtoVnode) String() string { return proto.CompactTextString(m) }
func (*PBProtoVnode) ProtoMessage()    {}

func (m *PBProtoVnode) GetId() []byte {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *PBProtoVnode) GetHost() string {
	if m != nil && m.Host != nil {
		return *m.Host
	}
	return ""
}

// PBProtoPing is simple structure for pinging remote vnodes.
type PBProtoPing struct {
	Version          *int64 `protobuf:"varint,1,req,name=version" json:"version,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PBProtoPing) Reset()         { *m = PBProtoPing{} }
func (m *PBProtoPing) String() string { return proto.CompactTextString(m) }
func (*PBProtoPing) ProtoMessage()    {}

func (m *PBProtoPing) GetVersion() int64 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

// PBProtoAck is generic response message with boolean 'ok' state.
type PBProtoAck struct {
	Version          *int64 `protobuf:"varint,1,req,name=version" json:"version,omitempty"`
	Ok               *bool  `protobuf:"varint,2,req,name=ok" json:"ok,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PBProtoAck) Reset()         { *m = PBProtoAck{} }
func (m *PBProtoAck) String() string { return proto.CompactTextString(m) }
func (*PBProtoAck) ProtoMessage()    {}

func (m *PBProtoAck) GetVersion() int64 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *PBProtoAck) GetOk() bool {
	if m != nil && m.Ok != nil {
		return *m.Ok
	}
	return false
}

// PBProtoErr defines error message.
type PBProtoErr struct {
	Error            *string `protobuf:"bytes,2,req,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *PBProtoErr) Reset()         { *m = PBProtoErr{} }
func (m *PBProtoErr) String() string { return proto.CompactTextString(m) }
func (*PBProtoErr) ProtoMessage()    {}

func (m *PBProtoErr) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

// PBProtoForward is sent to caller if request should be forwarded to another vnode.
type PBProtoForward struct {
	Vnode            *PBProtoVnode `protobuf:"bytes,1,req,name=vnode" json:"vnode,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *PBProtoForward) Reset()         { *m = PBProtoForward{} }
func (m *PBProtoForward) String() string { return proto.CompactTextString(m) }
func (*PBProtoForward) ProtoMessage()    {}

func (m *PBProtoForward) GetVnode() *PBProtoVnode {
	if m != nil {
		return m.Vnode
	}
	return nil
}

// PBProtoLeave (not used)
type PBProtoLeave struct {
	Source           *PBProtoVnode `protobuf:"bytes,1,req,name=source" json:"source,omitempty"`
	Dest             *PBProtoVnode `protobuf:"bytes,2,req,name=dest" json:"dest,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *PBProtoLeave) Reset()         { *m = PBProtoLeave{} }
func (m *PBProtoLeave) String() string { return proto.CompactTextString(m) }
func (*PBProtoLeave) ProtoMessage()    {}

func (m *PBProtoLeave) GetSource() *PBProtoVnode {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *PBProtoLeave) GetDest() *PBProtoVnode {
	if m != nil {
		return m.Dest
	}
	return nil
}

// PBProtoListVnodes - request the list of vnodes from remote vnode.
type PBProtoListVnodes struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *PBProtoListVnodes) Reset()         { *m = PBProtoListVnodes{} }
func (m *PBProtoListVnodes) String() string { return proto.CompactTextString(m) }
func (*PBProtoListVnodes) ProtoMessage()    {}

// PBProtoListVnodesResp is a structure for returning multiple vnodes to a caller.
type PBProtoListVnodesResp struct {
	Vnodes           []*PBProtoVnode `protobuf:"bytes,1,rep,name=vnodes" json:"vnodes,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *PBProtoListVnodesResp) Reset()         { *m = PBProtoListVnodesResp{} }
func (m *PBProtoListVnodesResp) String() string { return proto.CompactTextString(m) }
func (*PBProtoListVnodesResp) ProtoMessage()    {}

func (m *PBProtoListVnodesResp) GetVnodes() []*PBProtoVnode {
	if m != nil {
		return m.Vnodes
	}
	return nil
}

// PBProtoFindSuccessors is a structure to request successors for a key.
type PBProtoFindSuccessors struct {
	Key              []byte        `protobuf:"bytes,1,req,name=key" json:"key,omitempty"`
	Dest             *PBProtoVnode `protobuf:"bytes,2,req,name=dest" json:"dest,omitempty"`
	Limit            *int32        `protobuf:"varint,3,opt,name=limit" json:"limit,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *PBProtoFindSuccessors) Reset()         { *m = PBProtoFindSuccessors{} }
func (m *PBProtoFindSuccessors) String() string { return proto.CompactTextString(m) }
func (*PBProtoFindSuccessors) ProtoMessage()    {}

func (m *PBProtoFindSuccessors) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *PBProtoFindSuccessors) GetDest() *PBProtoVnode {
	if m != nil {
		return m.Dest
	}
	return nil
}

func (m *PBProtoFindSuccessors) GetLimit() int32 {
	if m != nil && m.Limit != nil {
		return *m.Limit
	}
	return 0
}

// PBProtoGetPredecessor - request immediate predecessor from vnode.
type PBProtoGetPredecessor struct {
	Dest             *PBProtoVnode `protobuf:"bytes,1,req,name=dest" json:"dest,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *PBProtoGetPredecessor) Reset()         { *m = PBProtoGetPredecessor{} }
func (m *PBProtoGetPredecessor) String() string { return proto.CompactTextString(m) }
func (*PBProtoGetPredecessor) ProtoMessage()    {}

func (m *PBProtoGetPredecessor) GetDest() *PBProtoVnode {
	if m != nil {
		return m.Dest
	}
	return nil
}

// PBProtoNotify is a message to notify the remote vnode of origin's existence.
type PBProtoNotify struct {
	Dest             *PBProtoVnode `protobuf:"bytes,1,req,name=dest" json:"dest,omitempty"`
	Vnode            *PBProtoVnode `protobuf:"bytes,2,req,name=vnode" json:"vnode,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *PBProtoNotify) Reset()         { *m = PBProtoNotify{} }
func (m *PBProtoNotify) String() string { return proto.CompactTextString(m) }
func (*PBProtoNotify) ProtoMessage()    {}

func (m *PBProtoNotify) GetDest() *PBProtoVnode {
	if m != nil {
		return m.Dest
	}
	return nil
}

func (m *PBProtoNotify) GetVnode() *PBProtoVnode {
	if m != nil {
		return m.Vnode
	}
	return nil
}

func init() {
}
