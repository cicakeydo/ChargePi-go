package grpc

import (
	"github.com/ChargePi/ChargePi-go/internal/diagnostics"
	"github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	"github.com/golang/protobuf/ptypes/empty"
)

type LogService struct {
	grpc.UnimplementedLogServer
	service diagnostics.Service
}

func NewLogService(service diagnostics.Service) *LogService {
	return &LogService{
		service: service,
	}
}

func (s *LogService) GetLogs(e *empty.Empty, server grpc.Log_GetLogsServer) error {
	// todo either a file-watcher or pipe directly from logrus (hook)?
	return nil
}

func (s *LogService) mustEmbedUnimplementedLogServer() {
}
