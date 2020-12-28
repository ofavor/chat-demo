package chat

import (
	"chat-srv/chat/msg"
	"context"
	"encoding/json"
	"errors"
)

type chatHandlerImpl struct {
	manager *Manager
}

// NewChatHandler create new chat handler
func NewChatHandler(m *Manager) ChatHandler {
	return &chatHandlerImpl{
		manager: m,
	}
}

func (m *chatHandlerImpl) onCreateRoom(meta map[string]string, cmd *msg.CmdCreateRoom) error {
	// svid := meta["server_id"]
	// ssid := meta["session_id"]
	// uid := meta["uid"]
	// if m.manager.isMemberInRoom(uid) {
	// 	// TODO
	// }
	return nil
}

func (m *chatHandlerImpl) onJoinRoom(meta map[string]string, cmd *msg.CmdJoinRoom) error {
	// r, err := m.manager.getRoom(cmd.ID)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (m *chatHandlerImpl) onQuitRoom(meta map[string]string, cmd *msg.CmdQuitRoom) error {
	// r, err := m.manager.getRoom(cmd.ID)
	// if err != nil {
	// 	return err
	// }
	// svID := meta["server_id"]
	// ssID := meta["session_id"]
	// r.removeMember(svID + ":" + ssID)
	return nil
}

func (m *chatHandlerImpl) onMessage(meta map[string]string, cmd *msg.CmdMessage) error {
	// svID := meta["server_id"]
	// ssID := meta["session_id"]
	// mid := svID + ":" + ssID
	// r, err := m.manager.getMemberRoom(mid)
	// if err != nil {
	// 	return err
	// }
	// r.message(mid, cmd.Text)
	return nil
}

// Command handle client chat command
func (m *chatHandlerImpl) Command(ctx context.Context, in *CommandRequest, out *CommandResponse) error {
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
