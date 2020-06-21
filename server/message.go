package server

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/kaisawind/message/pb"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

// Publish publish topic message
func (s *Server) Publish(_ context.Context, req *pb.PubReq) (resp *pb.PubResp, err error) {
	s.subscribers.Range(func(key, value interface{}) bool {
		pubChan, ok := value.(chan *pb.PubReq)
		if !ok {
			return true
		}
		select {
		case <-s.quit:
			logrus.Warningln("receive quit signal from server")
			return false
		case pubChan <- req:
			logrus.Infoln("receive pub req", req)
		default:
			logrus.Warningln("lose pub req")
		}
		return true
	})
	return &pb.PubResp{Status: 0}, nil
}

// Subscribe subscribe topic or topics
func (s *Server) Subscribe(req *pb.SubReq, stream pb.Message_SubscribeServer) (err error) {
	id := xid.New().String()
	pubChan := make(chan *pb.PubReq)
	defer close(pubChan)
	s.subscribers.Store(id, pubChan)

Loop:
	for {
		select {
		case <-s.quit:
			logrus.Warningln("receive quit signal from server")
			s.subscribers.Delete(id)
			break Loop
		case <-stream.Context().Done():
			logrus.Warningln("receive done signal from client")
			s.subscribers.Delete(id)
			break Loop
		case pubReq, ok := <-pubChan:
			if !ok {
				break Loop
			}
			switch val := req.Oneof.(type) {
			case *pb.SubReq_Topic:
				if val.Topic == pubReq.Topic {
					err := stream.Send(pubReq)
					if err != nil {
						logrus.WithError(err).Errorln("stream send pub req error")
						continue Loop
					}
				}
			case *pb.SubReq_Topics:
				for _, topic := range val.Topics.Topics {
					if topic == pubReq.Topic {
						err := stream.Send(pubReq)
						if err != nil {
							logrus.WithError(err).Errorln("stream send pub req error")
							continue Loop
						}
						break
					}
				}
			default:
				logrus.Warningln("unsupported req type", fmt.Sprintf("%T", req.Oneof))
			}
		}
	}
	return
}

// Ping test grpc server health
func (s *Server) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
