package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func UnaryGrpcRequestTimeInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		reqTime := time.Since(start)
		logrus.Infof("request time: %v: %v", info.FullMethod, reqTime)
		return resp, err
	}
}

func UnaryRequestTimeInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		reqTime := time.Since(start)
		logrus.Infof("request time: %v: %v", method, reqTime)
		return err
	}
}
