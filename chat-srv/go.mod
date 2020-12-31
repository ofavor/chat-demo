module chat-srv

go 1.15

require (
	gateway v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v8 v8.4.4
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.1.2
	github.com/ofavor/micro-lite v0.0.0-20201228025941-028dbd235880
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/grpc/examples v0.0.0-20201226181154-53788aa5dcb4 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace gateway => ../gateway

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/ofavor/micro-lite => ../../micro-lite // micro-lite is undering development, use local repo

replace github.com/ofavor/socket-gw => ../../socket-gw // socket-gw is undering development, use local repo
