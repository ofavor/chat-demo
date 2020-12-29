package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ofavor/socket-gw/transport"

	"chat-srv/chat/msg"
	"client/log"

	"github.com/ofavor/socket-gw/client"
)

func parseCmd(text string) []string {
	arr := strings.Split(strings.ToLower(text), " ")
	cmd := []string{}
	for _, a := range arr {
		t := strings.TrimSpace(a)
		if t != "" {
			cmd = append(cmd, t)
		}
	}
	return cmd
}

func buildPacketData(t int, data interface{}) []byte {
	j1, _ := json.Marshal(data)
	j2, _ := json.Marshal(map[string]interface{}{"type": t, "data": string(j1)})
	return j2
}

var (
	roomId = ""
)

func main() {
	client := client.NewClient(
		client.LogLevel("debug"),
		client.Address("127.0.0.1:8888"),
	)
	if err := client.Connect(); err != nil {
		log.Fatal("Connect error:", err)
	}
	if err := client.Send(transport.NewPacket(transport.PacketTypeAuth, []byte("abcd"))); err != nil {
		log.Fatal("Auth error:", err)
	}
	log.Info("Waiting for auth ACK ...")
	if _, err := client.Recv(); err != nil {
		log.Fatal(err)
	}

	log.Info("Connected!")
	go func() {
		for {
			p, err := client.Recv()
			if err != nil {
				log.Error("Receive error:", err)
				os.Exit(1)
			}
			fmt.Println(string(p.Body))
			if p.Type == 12 {
				j1 := map[string]string{}
				json.Unmarshal(p.Body, &j1)
				j2 := map[string]interface{}{}
				json.Unmarshal([]byte(j1["data"]), &j2)
				fmt.Println(j2)
			}
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		cmd := parseCmd(text)
		switch cmd[0] {
		case "create":
			cmd := &msg.CmdCreateRoom{
				Name:     cmd[1],
				Nickname: cmd[2],
			}
			data := buildPacketData(1, cmd)
			client.Send(transport.NewPacket(12, data))
		case "join":
			cmd := &msg.CmdJoinRoom{
				ID:       cmd[1],
				Nickname: cmd[2],
			}
			data := buildPacketData(2, cmd)
			client.Send(transport.NewPacket(12, data))
		case "quit":
			cmd := &msg.CmdQuitRoom{
				ID: cmd[1],
			}
			data := buildPacketData(3, cmd)
			client.Send(transport.NewPacket(12, data))
		case "send":
			cmd := &msg.CmdMessage{
				Message: cmd[1],
			}
			data := buildPacketData(4, cmd)
			client.Send(transport.NewPacket(12, data))
		}
	}
}
