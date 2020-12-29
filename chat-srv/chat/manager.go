package chat

import (
	"context"
	"encoding/json"
	"errors"
	"gateway/session"
	"sync"

	"chat-srv/chat/msg"
	"chat-srv/log"

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

// GetRoom get room by id
func (m *Manager) GetRoom(id string) (*Room, error) {
	m.RLock()
	defer m.RUnlock()
	r, ok := m.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}
	return r, nil
}

// GetMemberRoom get member's room by member id
func (m *Manager) GetMemberRoom(id string) (*Room, error) {
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

func (m *Manager) IsMemberInRoom(id string) bool {
	_, ok := m.members[id]
	return ok
}

func (m *Manager) OnRoomCreated(r *Room) {
	log.Debugf("Room %s is created:", r.id)
	m.Lock()
	m.rooms[r.id] = r
	m.Unlock()
}

func (m *Manager) OnMemberJoin(r *Room, mem *Member) {
	log.Debugf("Member %s is joined into Room %s:", mem.id, r.id)
	m.Lock()
	m.members[mem.id] = r.id
	m.Unlock()
	log.Debug("Current room members:", m.members)
}

func (m *Manager) OnMemberQuit(r *Room, mem *Member) {
	m.Lock()
	delete(m.rooms, r.id)
	delete(m.members, mem.id)
	m.Unlock()
	log.Debug("Current room members:", m.members)
}

func (m *Manager) OnMessage(r *Room, mem *Member, txt string) {
	// TODO
}

func (m *Manager) Send(mem *Member, t msg.Type, n interface{}) {
	cli := session.NewSessionService("chat-demo.gateway", m.service.Client())
	ssid := mem.meta["session_id"]
	j1, _ := json.Marshal(n)
	j2, _ := json.Marshal(map[string]interface{}{"type": t, "data": string(j1)})
	in := &session.Request{
		Id:   ssid,
		Type: 12,
		Data: j2,
	}
	cli.Send(context.Background(), in)
}
