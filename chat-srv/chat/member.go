package chat

import (
	"chat-srv/chat/msg"
)

// Member chat member
type Member struct {
	id       string
	nickname string
	meta     map[string]string
	sendFunc MemberSend
}

type MemberSend func(mem *Member, t msg.Type, n interface{}) error

// newMember create new member
func newMember(id string, nn string, meta map[string]string, send MemberSend) *Member {
	return &Member{
		id:       id,
		nickname: nn,
		meta:     meta,
		sendFunc: send,
	}
}

// Send data to member
func (m *Member) Send(t msg.Type, n interface{}) {
	m.sendFunc(m, t, n)
}
