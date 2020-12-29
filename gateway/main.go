package main

import (
	"flag"
	"fmt"
	"gateway/log"
	"gateway/session"
	"strings"
	"time"

	"github.com/ofavor/micro-lite"
	gw "github.com/ofavor/socket-gw"
)

func main() {
	regAddrs := flag.String("registry_addrs", "127.0.0.1:2379", "registry address list, splitted by ','")
	flag.Parse()

	fmt.Println("registry address:", *regAddrs)

	log.SetLevel("debug")
	log.Info("Gateway is starting ...")

	service := micro.NewService(
		micro.LogLevel("debug"),
		micro.Name("chat-demo.gateway"),
		micro.Address(":9999"),
		micro.RegistryAddrs(strings.Split(*regAddrs, ",")),
	)

	sessionMgr := session.NewManager(service)
	gw := gw.NewGateway(
		gw.LogLevel("debug"),
		gw.Address(":8888"),
		gw.SessionAuth(true),
		gw.SessionHandler(sessionMgr),
	)
	go func() {
		for range time.Tick(5 * time.Second) {
			log.Infof("Current session count:%d", sessionMgr.GetTotalSessions())
		}
	}()
	if err := gw.Run(); err != nil {
		log.Fatal("Gateway run error:", err)
	}

	// register rpc handler
	session.RegisterSessionHandler(service.Server(), session.NewSessionHandler(sessionMgr))

	if err := service.Run(); err != nil {
		log.Fatal("service run error:", err)
	}
}
