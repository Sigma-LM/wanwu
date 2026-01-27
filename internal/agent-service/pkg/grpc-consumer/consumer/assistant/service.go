package assistant

import (
	assistant_service "github.com/UnicomAI/wanwu/api/proto/assistant-service"
	"github.com/UnicomAI/wanwu/internal/agent-service/pkg/config"
	grpc_consumer "github.com/UnicomAI/wanwu/internal/agent-service/pkg/grpc-consumer"
	"google.golang.org/grpc"
)

var client assistant_service.AssistantServiceClient
var service = &ConsumerAssistantService{}

func init() {
	grpc_consumer.AddGrpcConsumerContainer(service)
}

type ConsumerAssistantService struct {
}

func (s *ConsumerAssistantService) GrpcConsumerType() string {
	return grpc_consumer.ConsumerAssistant
}

func (s *ConsumerAssistantService) BuildConfig() *grpc_consumer.GrpcConsumerConfig {
	return &grpc_consumer.GrpcConsumerConfig{
		Host:          config.GetConfig().Microservices.Assistant.Host,
		SecondTimeout: 30,
	}
}

func (s *ConsumerAssistantService) NewClient(conn *grpc.ClientConn) error {
	// grpc clients
	client = assistant_service.NewAssistantServiceClient(conn)
	return nil
}

func GetClient() assistant_service.AssistantServiceClient {
	return client
}
