package grpcclient

import (
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/traceid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newConn(cfg config.ServiceConfig) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		cfg.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(traceid.TraceClientInterceptor),
	)
}
