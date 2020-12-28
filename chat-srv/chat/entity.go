package chat

import (
	sync "sync"

	"github.com/google/uuid"
)

type room struct {
	sync.RWMutex
	id      string
	name    string
	members map[string]*member
}

func newRoom(name string) *room {
	return &room{
		id:      uuid.New().String(),
		name:    name,
		members: make(map[string]*member),
	}
}

func (r *room) addMember(mem *member) {
	r.Lock()
	defer r.Unlock()
	r.members[mem.id] = mem
}

func (r *room) removeMember(id string) {
	r.Lock()
	defer r.Unlock()
	delete(r.members, id)
}

func (r *room) message(mid string, txt string) {
	for id, m := range r.members {
		if id != mid {
			m.send(txt)
		}
	}
}

type member struct {
	id       string
	nickname string
}

func newMember(id string, nn string) *member {
	return &member{
		id:       id,
		nickname: nn,
	}
}

func (m *member) send(txt string) {

}
