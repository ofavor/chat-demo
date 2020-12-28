package session

import (
	"context"
	"gateway/log"

	"github.com/ofavor/socket-gw/transport"
)

type sessionHandlerImpl struct {
	manager *Manager
}

// NewSessionHandler create new session handler
func NewSessionHandler(m *Manager) SessionHandler {
	return &sessionHandlerImpl{
		manager: m,
	}
}

// Send data to session
func (m *sessionHandlerImpl) Send(ctx context.Context, in *Request, out *Response) error {
	log.Debug("Sending data to session:", in)
	s, err := m.manager.GetSession(in.Id)
	if err != nil {
		return err
	}
	p := transport.NewPacket(transport.PacketType(in.Type), in.Data)
	return s.Send(p)
}

// Broadcast data to sessions
func (m *sessionHandlerImpl) Broadcast(ctx context.Context, in *Request, out *Response) error {
	log.Debug("Broadcast data to sessions:", in)
	ss, err := m.manager.GetAllSessions()
	if err != nil {
		return err
	}
	p := transport.NewPacket(transport.PacketType(in.Type), in.Data)
	for _, s := range ss {
		if err := s.Send(p); err != nil {
			log.Errorf("Broadcast to session %s error: %s", s.ID(), err)
		}
	}
	return nil
}
