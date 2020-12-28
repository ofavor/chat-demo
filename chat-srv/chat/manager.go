package chat

import (
	"errors"
	"sync"

	"github.com/ofavor/micro-lite"
)

// Manager chat manager
type Manager struct {
	sync.RWMutex
	service micro.Service
	rooms   map[string]*room  // rooms
	members map[string]string // member-room relationship
}

// NewManager create new chat manager
func NewManager(s micro.Service) *Manager {
	return &Manager{
		service: s,
		rooms:   make(map[string]*room),
		members: make(map[string]string),
	}
}

func (m *Manager) getRoom(id string) (*room, error) {
	m.RLock()
	defer m.RUnlock()
	r, ok := m.rooms[id]
	if !ok {
		return nil, errors.New("room not found")
	}
	return r, nil
}

func (m *Manager) getMemberRoom(id string) (*room, error) {
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
