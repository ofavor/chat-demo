package chat

import (
	"context"
	"errors"
	"sync"

	"chat-srv/chat/msg"

	"proto/session"

	"github.com/ofavor/micro-lite"
)

// Manager chat manager
type Manager struct {
	sync.RWMutex
	service micro.Service
	rooms   map[string]*Room  // rooms
	members map[string]string // member-room relationship
}

// NewManager create new chat manager
func NewManager(s micro.Service) *Manager {
	return &Manager{
		service: s,
		rooms:   make(map[string]*Room),
		members: make(map[string]string),
	}
}

func (m *Manager) getRoom(id string) (*Room, error) {
	m.RLock()
	defer m.RUnlock()
	r, ok := m.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}
	return r, nil
}

func (m *Manager) getMemberRoom(id string) (*Room, error) {
	m.RLock()
	defer m.RUnlock()
	rid, ok := m.members[id]
	if !ok {
		return nil, errors.New("not in a room")
	}
	r, ok := m.rooms[rid]
	if !ok {
		return nil, errors.New("room not found")
	}
	return r, nil
}

func (m *Manager) isMemberInRoom(id string) bool {
	_, ok := m.members[id]
	return ok
}

func (m *Manager) OnRoomCreated(r *Room) {
	m.Lock()
	m.rooms[r.id] = r
	m.Unlock()
}

func (m *Manager) OnMemberJoin(r *Room, mem *Member) {
	m.Lock()
	m.members[mem.id] = r.id
	m.Unlock()
}

func (m *Manager) OnMemberQuit(r *Room, mem *Member) {
	m.Lock()
	delete(m.rooms, r.id)
	delete(m.members, mem.id)
	m.Unlock()
}

func (m *Manager) OnMessage(r *Room, mem *Member, txt string) {
	// TODO
}

func (m *Manager) Send(mem *Member, t msg.Type, n interface{}) {
	cli := session.NewSessionService("chat-demo.gateway", m.service.Client())
	ssid := mem.meta["session_id"]
	in := &session.Request{
		Id:   ssid,
		Type: 12,
		// TODO
	}
	cli.Send(context.Background(), in)
}
