package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// package name for the liveness service.
const livenessServicePrefix = "/liveness."

// stream interceptor is used to print a log on each stream RPC.
func (b *Bootstrap) allStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	// log the method being called
	b.Logger.Info("stream rpc called", zap.String("method", info.FullMethod))

	// proceed to the actual handler
	return handler(srv, ss)
}

// allUnaryInterceptor interceptor checks the status of a service before running the gRPC function.
func (b *Bootstrap) allUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// if status is true, allow all services to proceed
	if b.Consensus.Memory.GetStatus() {
		b.Logger.Info("rpc called", zap.String("method", info.FullMethod))
		return handler(ctx, req)
	}

	// if status is false, only allow services in the liveness package
	if len(info.FullMethod) > len(livenessServicePrefix) && info.FullMethod[:len(livenessServicePrefix)] == livenessServicePrefix {
		b.Logger.Info("rpc called", zap.String("method", info.FullMethod))
		return handler(ctx, req) // allow liveness service to proceed
	}

	// block all other services
	return nil, status.Error(13, "service is not responding") // return an error for blocked services
}
