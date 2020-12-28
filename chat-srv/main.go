package main

import (
	csc "chat-srv/chat"
	"chat-srv/log"
	"flag"
	"fmt"
	"proto/chat"
	"strings"

	"github.com/ofavor/micro-lite"
)

func main() {
	regAddrs := flag.String("registry_addrs", "127.0.0.1:2379", "registry address list, splitted by ','")
	flag.Parse()

	fmt.Println("registry address:", *regAddrs)

	log.SetLevel("debug")
	log.Info("Gateway is starting ...")
	service := micro.NewService(
		micro.LogLevel("debug"),
		micro.Name("chat-demo.chat"),
		micro.RegistryAddrs(strings.Split(*regAddrs, ",")),
	)

	chatMgr := csc.NewManager(service)
	// register rpc
	chat.RegisterChatHandler(service.Server(), csc.NewChatHandler(chatMgr))

	if err := service.Run(); err != nil {
		log.Fatal("service run error:", err)
	}
}
