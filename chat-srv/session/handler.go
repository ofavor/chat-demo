package session

import (
	"chat-srv/chat"
	"chat-srv/chat/msg"
	"chat-srv/log"
	"chat-srv/tracer"
	"context"
	"encoding/json"
	"errors"
	"gateway/backend"
	"strconv"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"

	"github.com/ofavor/micro-lite/client"
	"github.com/ofavor/micro-lite/client/selector"
)

type Handler struct {
	manager *chat.Manager
}

func NewHandler(m *chat.Manager) *Handler {
	return &Handler{
		manager: m,
	}
}

// Connect handle session connect event
func (h *Handler) Connect(ctx context.Context, in *backend.StatusRequest, out *backend.StatusResponse) error {
	// nothing to do
	return nil
}

// Disconnect handle session disconnect event
func (h *Handler) Disconnect(ctx context.Context, in *backend.StatusRequest, out *backend.StatusResponse) error {
	uid := in.Meta["uid"]
	rid, _ := h.manager.GetMemberRoomID(uid)
	svid, err := h.manager.GetRoomServerID(rid)
	log.Debugf("Current server id:%s <==> Room server id:%s", h.manager.GetService().Server().ID(), svid)
	if err != nil {
		log.Error("Failed to get room server id:", err)
		return err
	}
	if svid != h.manager.GetService().Server().ID() { // forward to another server
		log.Debugf("Forwad status request to other server")
		cli := backend.NewBackendService("chat-demo.chat", h.manager.GetService().Client())
		cli.Disconnect(ctx, in, client.WithSelectOption(selector.WithIDFilter([]string{svid})))
		return nil
	}
	h.manager.Disconnect(rid, uid)
	return nil
}

// Data handle gateway data request
func (h *Handler) Data(ctx context.Context, in *backend.DataRequest, out *backend.DataResponse) error {
	log.Debug("Got data request:", in)
	tcx := in.Meta["trace_ctx"]
	spanCtx := &model.SpanContext{}
	json.Unmarshal([]byte(tcx), spanCtx)
	span := tracer.StartSpan("chat-srv.ondata", zipkin.Parent(*spanCtx))
	defer span.Finish()

	data := map[string]string{}
	if err := json.Unmarshal(in.Data, &data); err != nil {
		log.Error("Unmarshal data request error:", err)
		return err
	}
	rid, ok := data["room_id"]
	log.Debugf("Data request room_id:%s", rid)
	if ok {
		svid, err := h.manager.GetRoomServerID(rid)
		log.Debugf("Current server id:%s <==> Room server id:%s", h.manager.GetService().Server().ID(), svid)
		if err != nil {
			log.Error("Failed to get room server id:", err)
			return err
		}
		if svid != h.manager.GetService().Server().ID() { // forward to another server
			log.Debugf("Forwad data request to other server")
			cli := backend.NewBackendService("chat-demo.chat", h.manager.GetService().Client())
			cli.Data(ctx, in, client.WithSelectOption(selector.WithIDFilter([]string{svid})))
			return nil
		}
	}
	t, ok := data["type"]
	if !ok {
		log.Error("Bad request data, 'type' is required")
		return errors.New("bad request data")
	}
	span.Tag("data_type", t)
	typ, err := strconv.Atoi(t)
	if err != nil {
		log.Error("Bad request data, 'type' error:", err)
		return err
	}
	switch msg.Type(typ) {
	case msg.TypeCmdCreateRoom:
		h.manager.CreateRoom(data["name"], in.Meta["uid"], data["nickname"], in.Meta)
	case msg.TypeCmdJoinRoom:
		h.manager.JoinRoom(rid, in.Meta["uid"], data["nickname"], in.Meta)
	case msg.TypeCmdQuitRoom:
		h.manager.QuitRoom(rid, in.Meta["uid"])
	case msg.TypeCmdMessage:
		h.manager.Message(rid, in.Meta["uid"], data["message"])
	default:
		log.Error("Unsupported request data type:", t)
		return errors.New("Bad request data type")
	}
	return nil
}
