package chat

import (
	"chat-srv/chat/msg"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Room chat room
type Room struct {
	sync.RWMutex
	id      string
	name    string
	handler RoomHandler
	members map[string]*Member
}

// RoomHandler handler for room
type RoomHandler interface {
	OnRoomCreated(r *Room)

	OnMemberJoin(r *Room, m *Member)

	OnMemberQuit(r *Room, m *Member)

	OnMessage(r *Room, m *Member, txt string)
}

var roomIDCounter = int32(100)

// NewRoom create new room
func NewRoom(h RoomHandler, name string) *Room {
	r := &Room{
		id:      fmt.Sprintf("%d", atomic.AddInt32(&roomIDCounter, 1)),
		name:    name,
		handler: h,
		members: make(map[string]*Member),
	}
	h.OnRoomCreated(r)
	return r
}

// MemberJoin join
func (r *Room) MemberJoin(mem *Member) {
	r.Lock()
	r.members[mem.id] = mem
	r.Unlock()

	r.handler.OnMemberJoin(r, mem)

	for _, mm := range r.members {
		if mm.id != mem.id {
			// send notification
			n := &msg.NotifyJoinRoom{
				ID: r.id,
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

	r.handler.OnMemberQuit(r, mem)

	for _, mm := range r.members {
		if mm.id != mem.id {
			// send notification
			n := &msg.NotifyQuitRoom{
				ID: r.id,
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

// Message send message
func (r *Room) Message(id string, txt string) {
	mem, _ := r.members[id]
	for _, mm := range r.members {
		if mm.id != mem.id {
			// send notification
			n := &msg.NotifyMessage{
				ID: r.id,
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
	r.handler.OnMessage(r, mem, txt)
}
