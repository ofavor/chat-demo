package session

import (
	"context"
	"fmt"
	"sync/atomic"

	"proto/chat"

	"github.com/ofavor/micro-lite"
	"github.com/ofavor/socket-gw/session"
	"github.com/ofavor/socket-gw/transport"
	"google.golang.org/protobuf/proto"
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
	return nil
}

// OnSessionReceived handle packet received from session
func (m *Manager) OnSessionReceived(s *session.Session, p *transport.Packet) error {
	switch p.Type {
	case PacketCustomTypeEcho: // send packet back
		s.Send(p)
	case PacketCustomTypeChat: // forward to chat service
		cr := &chat.CommandRequest{}
		if err := proto.Unmarshal(p.Body, cr); err != nil {
			return err
		}
		cr.Meta["server_id"] = m.service.Server().ID()
		cr.Meta["session_id"] = s.ID()
		cr.Meta["uid"] = s.Meta()["uid"]
		chatSvc := chat.NewChatService("chat-demo.chat", m.service.Client())
		if _, err := chatSvc.Command(context.Background(), cr); err != nil {
			return err
		}
	}
	return nil
}
