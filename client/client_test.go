package client_test

import (
	"context"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/kaisawind/broker.grpc/client"
	"github.com/kaisawind/broker.grpc/server"
)

func tryHostAndPort() (host, port string, err error) {
	port = strconv.Itoa(6000 + rand.Intn(100))
	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", port))
	if err != nil {
		return tryHostAndPort()
	}
	err = listener.Close()
	return
}

func TestNewGRPCClient(t *testing.T) {
	var err error
	s := server.NewServer()
	s.GRPCHost, s.GRPCPort, err = tryHostAndPort()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer s.Close()
	go func() {
		err := s.Serve()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}()

	c, err := client.NewGRPCClient(s.GRPCHost, s.GRPCPort)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if s.GRPCPort != c.Port() {
		t.Error("server port isn't equal client port")
		t.FailNow()
	}
	if s.GRPCHost != c.Host() {
		t.Error("server host isn't equal client host")
		t.FailNow()
	}
	err = c.Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestGrpcClient_Message(t *testing.T) {
	var err error
	s := server.NewServer()
	s.GRPCHost, s.GRPCPort, err = tryHostAndPort()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer s.Close()
	go func() {
		err := s.Serve()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}()

	time.Sleep(100 * time.Millisecond)
	c, err := client.NewGRPCClient(s.GRPCHost, s.GRPCPort)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer c.Close()
	_, err = c.Message().Ping(context.TODO(), &empty.Empty{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
