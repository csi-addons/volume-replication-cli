package client

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/kubernetes-csi/csi-lib-utils/connection"
	"github.com/kubernetes-csi/csi-lib-utils/rpc"
)

// Connect to the GRPC client
func Connect(address string) (*grpc.ClientConn, error) {
	return connection.Connect(address, nil, connection.OnConnectionLoss(connection.ExitOnConnectionLoss()))
}

// Probe the GRPC client once
func Probe(conn *grpc.ClientConn, timeout time.Duration) error {
	return rpc.ProbeForever(conn, timeout)
}

// GetDriverName gets the driver name from the CSI driver
func GetDriverName(conn *grpc.ClientConn, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return rpc.GetDriverName(ctx, conn)
}
