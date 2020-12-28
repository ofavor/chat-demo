module chat-srv

go 1.15

require (
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/ofavor/micro-lite v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/ofavor/micro-lite => ../../micro-lite // micro-lite is undering development, use local repo
