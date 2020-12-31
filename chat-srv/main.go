package main

import (
	"chat-srv/chat"
	"chat-srv/log"
	"chat-srv/session"
	"chat-srv/tracer"
	"flag"
	"fmt"
	"gateway/backend"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/ofavor/micro-lite"
)

func main() {
	regAddrs := flag.String("registry_addrs", "127.0.0.1:2379", "registry address list, splitted by ','")
	rdsAddr := flag.String("redis_addr", "127.0.0.1:6379", "redis server address")
	tracerAddr := flag.String("tracer_addr", "127.0.0.1:9411", "tracer(zipkin) address")
	flag.Parse()

	fmt.Println("registry address:", *regAddrs)
	fmt.Println("tracer address:", *tracerAddr)
	tracer.Init(*tracerAddr)

	log.SetLevel("debug")
	log.Info("Gateway is starting ...")
	service := micro.NewService(
		micro.LogLevel("debug"),
		micro.Name("chat-demo.chat"),
		micro.RegistryAddrs(strings.Split(*regAddrs, ",")),
	)

	rds := redis.NewClient(&redis.Options{Addr: *rdsAddr})

	chatMgr := chat.NewManager(service, rds)

	// register rpc
	backend.RegisterBackendHandler(service.Server(), session.NewHandler(chatMgr))

	if err := service.Run(); err != nil {
		log.Fatal("service run error:", err)
	}
}
