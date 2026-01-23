package grpc_consumer

import "google.golang.org/grpc"

const (
	ConsumerAssistant = "GrpcConsumerAssistant"
)

type GrpcConsumerConfig struct {
	Host          string
	SecondTimeout int
}

type GrpcConsumerService interface {
	GrpcConsumerType() string
	BuildConfig() *GrpcConsumerConfig
	NewClient(conn *grpc.ClientConn) error
}
