package chat

import (
	"chat-srv/chat/msg"
)

// Member chat member
type Member struct {
	id        string
	nickname  string
	meta      map[string]string
	transport Transport
}

// Transport send
type Transport interface {
	Send(mem *Member, t msg.Type, n interface{})
}

// NewMember create new member
func NewMember(id string, nn string, meta map[string]string, trans Transport) *Member {
	return &Member{
		id:        id,
		nickname:  nn,
		meta:      meta,
		transport: trans,
	}
}

// Send data to member
func (m *Member) Send(t msg.Type, n interface{}) {
	m.transport.Send(m, t, n)
}
