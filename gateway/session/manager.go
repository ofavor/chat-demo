package session

import (
	"context"
	"fmt"
	"gateway/backend"
	"gateway/log"
	"sync/atomic"

	"github.com/ofavor/micro-lite"
	"github.com/ofavor/socket-gw/session"
	"github.com/ofavor/socket-gw/transport"
)

var (
	// PacketCustomTypeEcho echo
	PacketCustomTypeEcho transport.PacketType = 11 // echo

	// PacketCustomTypeChat chat
	PacketCustomTypeChat transport.PacketType = 12 // chat
)

// Manager of session
type Manager struct {
	session.BaseHandler
	service micro.Service
}

// NewManager create new session manager
func NewManager(s micro.Service) *Manager {
	return &Manager{
		BaseHandler: session.NewBaseHandler(),
		service:     s,
	}
}

var uidCounter = int32(10)

// OnSessionAuth handle session auth procedure
func (m *Manager) OnSessionAuth(s *session.Session, p *transport.Packet) error {
	// bind uid to session
	s.Meta()["uid"] = fmt.Sprintf("%d", atomic.AddInt32(&uidCounter, 1))
	log.Debugf("Session %s meta:%v", s.ID(), s.Meta())
	return nil
}

var backendMapping = map[transport.PacketType]string{
	PacketCustomTypeChat: "chat-demo.chat",
}

// OnSessionReceived handle packet received from session
func (m *Manager) OnSessionReceived(s *session.Session, p *transport.Packet) error {
	switch p.Type {
	case PacketCustomTypeEcho: // send packet back
		s.Send(p)
	default: // forward to backend service
		sn, ok := backendMapping[p.Type]
		if !ok {
			log.Errorf("No backend service for packet type:%d", p.Type)
		} else {
			dr := &backend.DataRequest{
				Id:   s.ID(),
				Type: uint32(p.Type),
				Data: p.Body,
				Meta: make(map[string]string),
			}
			dr.Meta["server_id"] = m.service.Server().ID()
			dr.Meta["session_id"] = s.ID()
			dr.Meta["uid"] = s.Meta()["uid"]
			cli := backend.NewBackendService(sn, m.service.Client())
			if _, err := cli.Data(context.Background(), dr); err != nil {
				log.Error("Send data to backend service error:", err)
			}
		}
	}
	return nil
}

// OnSessionClosed handle session closed event
func (m *Manager) OnSessionClosed(s *session.Session) error {
	go func() {
		for _, sn := range backendMapping {
			sr := &backend.StatusRequest{
				Id:   s.ID(),
				Meta: make(map[string]string),
			}
			sr.Meta["server_id"] = m.service.Server().ID()
			sr.Meta["session_id"] = s.ID()
			sr.Meta["uid"] = s.Meta()["uid"]
			cli := backend.NewBackendService(sn, m.service.Client())
			if _, err := cli.Disconnect(context.Background(), sr); err != nil {
				log.Error("Send disconnect notification to backend service error:", err)
			}
		}
	}()
	return nil
}
