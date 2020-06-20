package server_test

import (
	"context"
	"io"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/kaisawind/message/client"
	pb "github.com/kaisawind/message/pb"
	"github.com/kaisawind/message/server"
)

func tryHostAndPort() (host, port string, err error) {
	port = strconv.Itoa(6000 + rand.Intn(100))
	listener, err := net.Listen("tcp", net.JoinHostPort("localhost", port))
	if err != nil {
		return tryHostAndPort()
	}
	listener.Close()
	return
}

func TestServer_Publish(t *testing.T) {
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
	val1 := &pb.PubResp{Status: 1}
	req, err := ptypes.MarshalAny(val1)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	val2 := &pb.PubResp{Status: 2}
	resp, err := ptypes.MarshalAny(val2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	_, err = c.Message().Publish(context.TODO(), &pb.PubReq{
		Topic: "topic",
		Req:   req,
		Resp:  resp,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestServer_Subscribe(t *testing.T) {
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
	go func() {
		for i := int32(0); i < 11; i++ {
			val1 := &pb.PubResp{Status: i}
			req, err := ptypes.MarshalAny(val1)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			val2 := &pb.PubResp{Status: 2 * i}
			resp, err := ptypes.MarshalAny(val2)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			_, err = c.Message().Publish(context.TODO(), &pb.PubReq{
				Topic: "topic",
				Req:   req,
				Resp:  resp,
			})
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
		}
	}()

	subReq := &pb.SubReq{}
	stream, err := c.Message().Subscribe(context.TODO(), subReq)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for {
		pubReq, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Error(err)
			t.FailNow()
		}
		t.Log("stream recv", pubReq)
		val1 := &pb.PubResp{}
		err = ptypes.UnmarshalAny(pubReq.Req, val1)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if val1.Status == 10 {
			break
		}
	}
}
