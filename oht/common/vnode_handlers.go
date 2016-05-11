package dendrite

import (
	"bytes"
)

/*
 VnodeHandler interface defines methods (from Transport interface) that are to be called in vnode context.
 Transports use this interface to avoid duplicate implementations. localVnode implements this interface.
*/
type VnodeHandler interface {
	FindSuccessors([]byte, int) ([]*Vnode, *Vnode, error) // args: key, limit # returns: succs, forward, error
	FindRemoteSuccessors(int) ([]*Vnode, error)
	GetPredecessor() (*Vnode, error)
	Notify(*Vnode) ([]*Vnode, error)
}

// localHandler is a handler object connecting a VnodeHandler and Vnode.
type localHandler struct {
	vn      *Vnode
	handler VnodeHandler
}

// FindSuccessors implements Transport's FindSuccessors() in vnode context.
func (vn *localVnode) FindSuccessors(key []byte, limit int) ([]*Vnode, *Vnode, error) {
	// check if we have direct successor for requested key
	succs := make([]*Vnode, 0)
	max_vnodes := min(limit, len(vn.successors))
	if bytes.Compare(key, vn.Id) == 0 || between(vn.Id, vn.successors[0].Id, key, true) {
		for i := 0; i < max_vnodes; i++ {
			if vn.successors[i] == nil {
				continue
			}
			succs = append(succs, &Vnode{
				Id:   vn.successors[i].Id,
				Host: vn.successors[i].Host,
			})
		}
		return succs, nil, nil
	}

	// if finger table has been initialized - forward request to closest finger
	// otherwise forward to my successor

	forward_vn := vn.closest_preceeding_finger(key)

	// if we got ourselves back, that's it - I'm the successor
	if bytes.Compare(forward_vn.Id, vn.Id) == 0 {
		succs = append(succs, &Vnode{
			Id:   vn.Id,
			Host: vn.Host,
		})
		for i := 1; i < max_vnodes; i++ {
			if vn.successors[i-1] == nil {
				break
			}
			succs = append(succs, vn.successors[i-1])
		}
		return succs, nil, nil
	}
	//log.Printf("findsuccessor (%X) forwarding to %X\n", vn.Id, forward_vn.Id)
	return nil, forward_vn, nil
}

// GetPredecessor implements Transport's GetPredecessor() in vnode context.
func (vn *localVnode) GetPredecessor() (*Vnode, error) {
	if vn.predecessor == nil {
		return nil, nil
	}
	return vn.predecessor, nil
}

// Notify is invoked when a Vnode gets notified.
func (vn *localVnode) Notify(maybe_pred *Vnode) ([]*Vnode, error) {
	// Check if we should update our predecessor
	if vn.predecessor == nil || between(vn.predecessor.Id, vn.Id, maybe_pred.Id, false) {
		var real_pred *Vnode

		if vn.predecessor == nil {
			if vn.old_predecessor != nil {
				// need to check against old predecessor here
				real_pred = vn.old_predecessor
			} else {
				real_pred = vn.predecessor
			}
		} else {
			real_pred = vn.predecessor
		}

		// before emiting anything, lets update our remotes
		vn.updateRemoteSuccessors()

		if real_pred == nil || between(real_pred.Id, vn.Id, maybe_pred.Id, false) {
			ctx := &EventCtx{
				EvType:        EvPredecessorJoined,
				Target:        &vn.Vnode,
				PrimaryItem:   maybe_pred,
				SecondaryItem: vn.old_predecessor,
			}
			vn.ring.emit(ctx)
		} else {
			ctx := &EventCtx{
				EvType:        EvPredecessorLeft,
				Target:        &vn.Vnode,
				PrimaryItem:   maybe_pred,
				SecondaryItem: vn.old_predecessor,
			}
			vn.ring.emit(ctx)
		}

		// maybe we're just joining and one of our local vnodes is closer to us than this predecessor
		vn.ring.Logf(LogInfo, "vn.Notify() - setting new predecessor for %x: %x\n", vn.Id, maybe_pred.Id)
		vn.predecessor = maybe_pred
	}

	// Return our successors list
	return vn.successors, nil
}

// FindRemoteSuccessors returns up to 'limit' successor vnodes,
// that are unique and do not reside on same physical node as vnode.
func (vn *localVnode) FindRemoteSuccessors(limit int) ([]*Vnode, error) {
	remote_succs := make([]*Vnode, 0)
	for _, succ := range vn.remote_successors {
		if succ == nil {
			continue
		}
		remote_succs = append(remote_succs, succ)
	}
	return remote_succs, nil
}
