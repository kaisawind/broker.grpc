package client

import (
	pb "github.com/kaisawind/message/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
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

func (c *grpcClient) Host() string {
	return c.host
}

func (c *grpcClient) Port() string {
	return c.port
}

func (c *grpcClient) Close() error {
	if c.conn != nil {
		logrus.Debugln("client: ", c.conn.GetState(), "->closed")
		return c.conn.Close()
	}
	return nil
}

func (c *grpcClient) Message() pb.MessageClient {
	return c.message
}
