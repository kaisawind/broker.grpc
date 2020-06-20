package client

import (
	"net"

	pb "github.com/kaisawind/message/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GRPCClient is a client interface for grpc service.
type GRPCClient interface {
	Host() string
	Port() string
	Message() pb.MessageClient
	Close() error
}

type grpcClient struct {
	host    string
	port    string
	conn    *grpc.ClientConn
	message pb.MessageClient
}

// NewGRPCClient create a new load-balanced client to talk to the grpc service.
func NewGRPCClient(host, port string) (GRPCClient, error) {
	address := net.JoinHostPort(host, port)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logrus.Errorf("did not connect: %v", err)
		return nil, err
	}
	return &grpcClient{
		host:    host,
		port:    port,
		conn:    conn,
		message: pb.NewMessageClient(conn),
	}, nil
}

// Host grpc host
func (c *grpcClient) Host() string {
	return c.host
}

// Port grpc server port
func (c *grpcClient) Port() string {
	return c.port
}

// Close close connection with grpc server
func (c *grpcClient) Close() error {
	if c.conn != nil {
		logrus.Debugln("client: ", c.conn.GetState(), "->closed")
		return c.conn.Close()
	}
	return nil
}

// Message message handler
func (c *grpcClient) Message() pb.MessageClient {
	return c.message
}
