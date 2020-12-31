package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/session"
	"sync"

	"chat-srv/chat/msg"
	"chat-srv/log"

	"github.com/go-redis/redis/v8"
	"github.com/ofavor/micro-lite"
)

// Manager chat manager
type Manager struct {
	sync.RWMutex
	service micro.Service
	rooms   map[string]*Room  // rooms
	members map[string]string // member-room relationship
	redis   *redis.Client
}

var (
	keyRoomIDCounter   = "key.room.id"
	keyMemberIDCounter = "key.member.id"
	keyRoomData        = "key.room.data.%s"
	keyMemberData      = "key.member.data.%s"
)

// NewManager create new chat manager
func NewManager(s micro.Service, rds *redis.Client) *Manager {
	return &Manager{
		service: s,
		rooms:   make(map[string]*Room),
		members: make(map[string]string),
		redis:   rds,
	}
}

func (m *Manager) GetService() micro.Service {
	return m.service
}
func (m *Manager) GetRoomServerID(id string) (string, error) {
	return m.redis.Get(context.Background(), fmt.Sprintf(keyRoomData, id)).Result()
}

func (m *Manager) GetMemberRoomID(uid string) (string, error) {
	return m.redis.Get(context.Background(), fmt.Sprintf(keyMemberData, uid)).Result()
}

// CreateRoom create new chat room by specifying name
func (m *Manager) CreateRoom(name string, uid string, nickname string, meta map[string]string) (*Room, error) {
	id, err := m.redis.Incr(context.Background(), keyRoomIDCounter).Result()
	if err != nil {
		log.Error("Generate room id error:", err)
		return nil, err
	}
	room := newRoom(fmt.Sprintf("%d", id), name)
	if _, err := m.redis.Set(context.Background(), fmt.Sprintf(keyRoomData, room.id), m.service.Server().ID(), 0).Result(); err != nil {
		log.Error("Store room data error:", err)
		return nil, err
	}
	log.Debugf("Room is created:%s", room.id)
	m.Lock()
	m.rooms[room.id] = room
	m.Unlock()

	m.JoinRoom(room.id, uid, nickname, meta)
	return room, nil
}

func (m *Manager) JoinRoom(id string, uid string, nickname string, meta map[string]string) error {
	m.RLock()
	room, _ := m.rooms[id]
	m.RUnlock()
	if _, err := m.redis.Set(context.Background(), fmt.Sprintf(keyMemberData, uid), room.id, 0).Result(); err != nil {
		log.Error("Store member data error:", err)
		return err
	}
	mem := newMember(uid, nickname, meta, m.Send)
	room.MemberJoin(mem)
	log.Debugf("Member %s join to room %s", mem.id, room.id)
	return nil
}

func (m *Manager) QuitRoom(id string, uid string) error {
	m.RLock()
	room, _ := m.rooms[id]
	m.RUnlock()
	room.MemberQuit(uid)
	m.redis.Del(context.Background(), fmt.Sprintf(keyMemberData, uid), room.id)
	log.Debugf("Member %s quit from room %s", uid, id)
	if room.IsEmpty() {
		m.Lock()
		delete(m.rooms, room.id)
		m.Unlock()
		log.Debugf("Room %s is empty, delete it", room.id)
		m.redis.Del(context.Background(), fmt.Sprintf(keyMemberData, uid), room.id)
	}
	return nil
}

func (m *Manager) Message(id string, uid string, msg string) error {
	m.RLock()
	room, _ := m.rooms[id]
	m.RUnlock()
	room.Message(uid, msg)
	log.Debugf("Member %s send message to room %s:%s", uid, id, msg)
	return nil
}

func (m *Manager) Disconnect(id string, uid string) error {
	m.RLock()
	room, _ := m.rooms[id]
	m.RUnlock()
	room.MemberDisconnect(uid)
	m.redis.Del(context.Background(), fmt.Sprintf(keyMemberData, uid), room.id)
	log.Debugf("Member %s disconnect from room %s", uid, id)
	if room.IsEmpty() {
		m.Lock()
		delete(m.rooms, room.id)
		m.Unlock()
		m.redis.Del(context.Background(), fmt.Sprintf(keyMemberData, uid), room.id)
		log.Debugf("Room %s is empty, delete it", room.id)
	}
	return nil
}

// // GetRoom get room by id
// func (m *Manager) GetRoom(id string) (*Room, error) {
// 	m.RLock()
// 	defer m.RUnlock()
// 	r, ok := m.rooms[id]
// 	if !ok {
// 		return nil, errors.New("room not found")
// 	}
// 	return r, nil
// }

// // GetMemberRoom get member's room by member id
// func (m *Manager) GetMemberRoom(id string) (*Room, error) {
// 	m.RLock()
// 	defer m.RUnlock()
// 	rid, ok := m.members[id]
// 	if !ok {
// 		return nil, errors.New("not in a room")
// 	}
// 	r, ok := m.rooms[rid]
// 	if !ok {
// 		return nil, errors.New("room not found")
// 	}
// 	return r, nil
// }

// func (m *Manager) IsMemberInRoom(id string) bool {
// 	_, ok := m.members[id]
// 	return ok
// }

func (m *Manager) Send(mem *Member, t msg.Type, n interface{}) error {
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
	return nil
}
