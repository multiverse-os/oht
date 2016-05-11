package dendrite

import (
	"sync"
)

// LocalTransport implements Transport interface, but is used for communicating between local vnodes.
type LocalTransport struct {
	host   string
	remote Transport
	lock   sync.RWMutex
	table  map[string]*localHandler
}

// InitLocalTransport initializes LocalTransport.
func InitLocalTransport(remote Transport) Transport {
	lt := &LocalTransport{
		remote: remote,
		table:  make(map[string]*localHandler),
	}
	return lt
}

// RegisterHook does nothing in local transport. Just satisfying interface.
func (lt *LocalTransport) RegisterHook(th TransportHook) {
}

// Decode does nothing in local transport. Just satisfying interface.
func (lt *LocalTransport) Decode(raw []byte) (*ChordMsg, error) {
	return nil, nil
}

// Encode does nothing in local transport. Just satisfying interface.
func (lt *LocalTransport) Encode(msgtype MsgType, data []byte) []byte {
	return nil
}

// Register registers a VnodeHandler within local and remote transports.
func (lt *LocalTransport) Register(vnode *Vnode, handler VnodeHandler) {
	// Register local instance
	lt.lock.Lock()
	lt.host = vnode.Host
	lt.table[vnode.String()] = &localHandler{vnode, handler}
	lt.lock.Unlock()

	// Register with remote transport
	lt.remote.Register(vnode, handler)
}

func (lt *LocalTransport) getVnodeHandler(vnode *Vnode) (VnodeHandler, bool) {
	lt.lock.Lock()
	defer lt.lock.Unlock()
	h, ok := lt.table[vnode.String()]
	if ok {
		return h.handler, ok
	}
	return nil, ok
}

// GetVnodeHandler returns registered local vnode handler, if one is found for given vnode.
func (lt *LocalTransport) GetVnodeHandler(vnode *Vnode) (VnodeHandler, bool) {
	return lt.getVnodeHandler(vnode)
}

// FindSuccessors implements Transport's FindSuccessors() in local transport.
func (lt *LocalTransport) FindSuccessors(vn *Vnode, limit int, key []byte) ([]*Vnode, error) {
	// Look for it locally
	handler, ok := lt.getVnodeHandler(vn)
	// If it exists locally, handle it
	if ok {
		succs, forward_vn, err := handler.FindSuccessors(key, limit)
		if err != nil {
			return nil, err
		}
		if forward_vn != nil {
			return lt.FindSuccessors(forward_vn, limit, key)
		}
		return succs, nil
	}

	// Pass onto remote
	return lt.remote.FindSuccessors(vn, limit, key)
}

// ListVnodes implements Transport's ListVnodes() in local transport.
func (lt *LocalTransport) ListVnodes(host string) ([]*Vnode, error) {
	// Check if this is a local host
	if host == lt.host {
		// Generate all the local clients
		res := make([]*Vnode, 0, len(lt.table))

		// Build list
		lt.lock.RLock()
		for _, v := range lt.table {
			res = append(res, v.vn)
		}
		lt.lock.RUnlock()

		return res, nil
	}

	// Pass onto remote
	return lt.remote.ListVnodes(host)
}

// Ping implements Transport's Ping() in local transport.
func (lt *LocalTransport) Ping(vn *Vnode) (bool, error) {
	// Look for it locally
	_, ok := lt.getVnodeHandler(vn)
	if ok {
		return true, nil
	}
	// ping remote
	return lt.remote.Ping(vn)
}

// GetPredecessor implements Transport's GetPredecessor() in local transport.
func (lt *LocalTransport) GetPredecessor(vn *Vnode) (*Vnode, error) {
	local_vn, ok := lt.getVnodeHandler(vn)
	if ok {
		return local_vn.GetPredecessor()
	}
	return lt.remote.GetPredecessor(vn)
}

// Notify implements Transport's Notify() in local transport.
func (lt *LocalTransport) Notify(dest, self *Vnode) ([]*Vnode, error) {
	// Look for it locally
	handler, ok := lt.getVnodeHandler(dest)

	// If it exists locally, handle it
	if ok {
		return handler.Notify(self)
	}

	// Pass onto remote
	return lt.remote.Notify(dest, self)
}
