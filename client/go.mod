module client

go 1.15

require (
	chat-srv v0.0.0-00010101000000-000000000000
	github.com/ofavor/socket-gw v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.7.0
)

replace github.com/ofavor/socket-gw => ../../socket-gw // socket-gw is undering development, use local repo

replace chat-srv => ../chat-srv

replace gateway => ../gateway
