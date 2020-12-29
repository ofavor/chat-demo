package session

import (
	"chat-srv/chat"
	"chat-srv/chat/msg"
	"chat-srv/log"
	"context"
	"encoding/json"
	"errors"
	"gateway/backend"
)

type Handler struct {
	manager *chat.Manager
}

func NewHandler(m *chat.Manager) *Handler {
	return &Handler{
		manager: m,
	}
}

func (h *Handler) onCreateRoom(meta map[string]string, cmd *msg.CmdCreateRoom) error {
	uid := meta["uid"]
	if h.manager.IsMemberInRoom(uid) {
		return errors.New("already in a room")
	}
	room := chat.NewRoom(h.manager, cmd.Name)
	mem := chat.NewMember(uid, cmd.Nickname, meta, h.manager)
	log.Debug("Member is created:", mem)
	room.MemberJoin(mem)
	return nil
}

func (h *Handler) onJoinRoom(meta map[string]string, cmd *msg.CmdJoinRoom) error {
	uid := meta["uid"]
	if h.manager.IsMemberInRoom(uid) {
		return errors.New("already in a room")
	}
	r, err := h.manager.GetRoom(cmd.ID)
	if err != nil {
		return err
	}
	mem := chat.NewMember(uid, cmd.Nickname, meta, h.manager)
	r.MemberJoin(mem)
	return nil
}

func (h *Handler) onQuitRoom(meta map[string]string, cmd *msg.CmdQuitRoom) error {
	uid := meta["uid"]
	if !h.manager.IsMemberInRoom(uid) {
		return errors.New("not in a room")
	}
	r, err := h.manager.GetRoom(cmd.ID)
	if err != nil {
		return err
	}
	r.MemberQuit(uid)
	return nil
}

func (h *Handler) onMessage(meta map[string]string, cmd *msg.CmdMessage) error {
	uid := meta["uid"]
	r, err := h.manager.GetMemberRoom(uid)
	if err != nil {
		return err
	}
	r.Message(uid, cmd.Message)
	return nil
}

func (h *Handler) Connect(ctx context.Context, in *backend.StatusRequest, out *backend.StatusResponse) error {
	return nil
}
func (h *Handler) Disconnect(ctx context.Context, in *backend.StatusRequest, out *backend.StatusResponse) error {
	return nil
}

type Data struct {
	Type int
	Data string
}

func (h *Handler) Data(ctx context.Context, in *backend.DataRequest, out *backend.DataResponse) error {
	log.Debug("Got data request:", in)
	data := &Data{}
	if err := json.Unmarshal(in.Data, data); err != nil {
		return err
	}
	switch msg.Type(data.Type) {
	case msg.TypeCmdCreateRoom:
		cmd := &msg.CmdCreateRoom{}
		json.Unmarshal([]byte(data.Data), cmd)
		h.onCreateRoom(in.Meta, cmd)
	case msg.TypeCmdJoinRoom:
		cmd := &msg.CmdJoinRoom{}
		json.Unmarshal([]byte(data.Data), cmd)
		h.onJoinRoom(in.Meta, cmd)
	case msg.TypeCmdQuitRoom:
		cmd := &msg.CmdQuitRoom{}
		json.Unmarshal([]byte(data.Data), cmd)
		h.onQuitRoom(in.Meta, cmd)
	case msg.TypeCmdMessage:
		cmd := &msg.CmdMessage{}
		json.Unmarshal([]byte(data.Data), cmd)
		h.onMessage(in.Meta, cmd)
	default:
		log.Error("Invalid data type:", data.Type)
	}
	return nil
}
