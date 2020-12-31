package chat

import (
	"chat-srv/chat/msg"
	"sync"
	"time"
)

// Room chat room
type Room struct {
	sync.RWMutex
	id      string
	name    string
	members map[string]*Member
}

var roomIDCounter = int32(100)

// newRoom create new room
func newRoom(id string, name string) *Room {
	r := &Room{
		id:      id,
		name:    name,
		members: make(map[string]*Member),
	}
	return r
}

// MemberJoin join
func (r *Room) MemberJoin(mem *Member) {
	r.Lock()
	r.members[mem.id] = mem
	r.Unlock()

	for _, mm := range r.members {
		if mm.id != mem.id {
			// send notification
			n := &msg.NotifyJoinRoom{
				ID: r.id,
				Room: msg.Room{
					ID:   r.id,
					Name: r.name,
				},
				Who: msg.Member{
					ID:       mem.id,
					Nickname: mem.nickname,
				},
				At: time.Now().Unix(),
			}
			mm.Send(msg.TypeNotifyJoinRoom, n)
		}
	}
}

// MemberQuit quit
func (r *Room) MemberQuit(id string) {
	r.Lock()
	mem, _ := r.members[id]
	delete(r.members, id)
	r.Unlock()

	for _, mm := range r.members {
		if mm.id != mem.id {
			// send notification
			n := &msg.NotifyQuitRoom{
				ID: r.id,
				Room: msg.Room{
					ID:   r.id,
					Name: r.name,
				},
				Who: msg.Member{
					ID:       mem.id,
					Nickname: mem.nickname,
				},
				At: time.Now().Unix(),
			}
			mm.Send(msg.TypeNotifyQuitRoom, n)
		}
	}
}

// MemberDisconnect disconnect
func (r *Room) MemberDisconnect(id string) {
	r.Lock()
	mem, _ := r.members[id]
	delete(r.members, id)
	r.Unlock()

	for _, mm := range r.members {
		if mm.id != mem.id {
			// send notification
			n := &msg.NotifyQuitRoom{
				ID: r.id,
				Room: msg.Room{
					ID:   r.id,
					Name: r.name,
				},
				Who: msg.Member{
					ID:       mem.id,
					Nickname: mem.nickname,
				},
				At: time.Now().Unix(),
			}
			mm.Send(msg.TypeNotifyDisconnected, n)
		}
	}
}

// Message send message
func (r *Room) Message(id string, txt string) {
	mem, _ := r.members[id]
	for _, mm := range r.members {
		if mm.id != mem.id {
			// send notification
			n := &msg.NotifyMessage{
				ID: r.id,
				Room: msg.Room{
					ID:   r.id,
					Name: r.name,
				},
				Who: msg.Member{
					ID:       mem.id,
					Nickname: mem.nickname,
				},
				Message: txt,
				At:      time.Now().Unix(),
			}
			mm.Send(msg.TypeNotifyMessage, n)
		}
	}
}

// IsEmpty check if the room is empty
func (r *Room)IsEmpty() bool {
	r.RLock()
	defer r.RUnlock()
	return len(r.members) == 0
}