package chat

import (
	"chat-srv/chat/msg"
	"context"
	"encoding/json"
	"errors"
	"proto/chat"
)

type chatHandlerImpl struct {
	manager *Manager
}

// NewChatHandler create new chat handler
func NewChatHandler(m *Manager) chat.ChatHandler {
	return &chatHandlerImpl{
		manager: m,
	}
}

func (m *chatHandlerImpl) onCreateRoom(meta map[string]string, cmd *msg.CmdCreateRoom) error {
	// svid := meta["server_id"]
	// ssid := meta["session_id"]
	uid := meta["uid"]
	if m.manager.isMemberInRoom(uid) {
		return errors.New("already in a room")
	}
	room := NewRoom(m.manager, cmd.Name)
	mem := NewMember(uid, cmd.Nickname, meta, m.manager)
	room.MemberJoin(mem)
	return nil
}

func (m *chatHandlerImpl) onJoinRoom(meta map[string]string, cmd *msg.CmdJoinRoom) error {
	uid := meta["uid"]
	if m.manager.isMemberInRoom(uid) {
		return errors.New("already in a room")
	}
	r, err := m.manager.getRoom(cmd.ID)
	if err != nil {
		return err
	}
	mem := NewMember(uid, cmd.Nickname, meta, m.manager)
	r.MemberJoin(mem)
	return nil
}

func (m *chatHandlerImpl) onQuitRoom(meta map[string]string, cmd *msg.CmdQuitRoom) error {
	uid := meta["uid"]
	if !m.manager.isMemberInRoom(uid) {
		return errors.New("not in a room")
	}
	r, err := m.manager.getRoom(cmd.ID)
	if err != nil {
		return err
	}
	r.MemberQuit(uid)
	return nil
}

func (m *chatHandlerImpl) onMessage(meta map[string]string, cmd *msg.CmdMessage) error {
	uid := meta["uid"]
	r, err := m.manager.getMemberRoom(uid)
	if err != nil {
		return err
	}
	r.Message(uid, cmd.Message)
	return nil
}

// Command handle client chat command
func (m *chatHandlerImpl) Command(ctx context.Context, in *chat.CommandRequest, out *chat.CommandResponse) error {
	switch msg.Type(in.Type) {
	case msg.TypeCmdCreateRoom:
		cmd := &msg.CmdCreateRoom{}
		if err := json.Unmarshal([]byte(in.Data), cmd); err != nil {
			return err
		}
		return m.onCreateRoom(in.Meta, cmd)
	case msg.TypeCmdJoinRoom:
		cmd := &msg.CmdJoinRoom{}
		if err := json.Unmarshal([]byte(in.Data), cmd); err != nil {
			return err
		}
		return m.onJoinRoom(in.Meta, cmd)
	case msg.TypeCmdQuitRoom:
		cmd := &msg.CmdQuitRoom{}
		if err := json.Unmarshal([]byte(in.Data), cmd); err != nil {
			return err
		}
		return m.onQuitRoom(in.Meta, cmd)
	case msg.TypeCmdMessage:
		cmd := &msg.CmdMessage{}
		if err := json.Unmarshal([]byte(in.Data), cmd); err != nil {
			return err
		}
		return m.onMessage(in.Meta, cmd)
	default:
		return errors.New("invalid command type")
	}
}
