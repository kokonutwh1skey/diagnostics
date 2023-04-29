package cmd

import (
	"sync"
	"time"

	"github.com/google/btree"
)

const nodeSessionMaxAge = 30 * 24 * time.Hour // 30 days

type UiSession struct {
	lock               sync.Mutex
	Session            bool
	SessionPin         uint64
	SessionName        string
	Errors             []string // Transient field - only filled for the time of template execution
	currentSessionName string
	NodeS              *NodeSession // Transient field - only filled for the time of template execution
	uiNodeTree         *btree.BTreeG[UiNodeSession]
	UiNodes            []UiNodeSession // Transient field - only filled forthe time of template execution
}

func (uiSession *UiSession) appendError(err string) {
	uiSession.lock.Lock()
	defer uiSession.lock.Unlock()
	uiSession.Errors = append(uiSession.Errors, err)
}

// NodeSession corresponds to one Erigon node connected via "erigon support" bridge to an operator
type NodeSession struct {
	lock           sync.Mutex
	sessionPin     uint64
	Connected      bool
	RemoteAddr     string
	SupportVersion uint64            // Version of the erigon support command
	requestCh      chan *NodeRequest // Channel for incoming metrics requests
	expires        time.Time
}

func (ns *NodeSession) connect(remoteAddr string) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	ns.Connected = true
	ns.RemoteAddr = remoteAddr
}

func (ns *NodeSession) disconnect() {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	ns.Connected = false
}

type UiNodeSession struct {
	SessionName string
	SessionPin  uint64
}