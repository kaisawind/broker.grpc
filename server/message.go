package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/kaisawind/message/pb"
	"github.com/sirupsen/logrus"
)

// Publish publish topic message
func (s *Server) Publish(_ context.Context, req *pb.PubReq) (resp *pb.PubResp, err error) {
	select {
	case <-s.quit:
	case s.pubChan <- req:
		logrus.Infoln("Publish PubReq", req)
	default:
		logrus.Warningln("Lose PubReq", req)
		return &pb.PubResp{Status: 1}, nil
	}
	return &pb.PubResp{Status: 0}, nil
}

// Subscribe subscribe topic or topics
func (s *Server) Subscribe(req *pb.SubReq, stream pb.Message_SubscribeServer) (err error) {
Loop:
	for {
		select {
		case <-s.quit:
			logrus.Warningln("pub chan closed")
			break Loop
		case pubReq, ok := <-s.pubChan:
			if !ok {
				break Loop
			}
			logrus.Infoln("Subscribe PubReq", pubReq)
			switch val := req.Oneof.(type) {
			case *pb.SubReq_Topic:
				if val.Topic != pubReq.Topic {
					continue
				}
			case *pb.SubReq_Topics:
				if val.Topics == nil {
					continue
				}
				has := false
				for _, topic := range val.Topics.Topics {
					if topic == pubReq.Topic {
						has = true
					}
				}
				if !has {
					continue
				}
			default:
				continue
			}
			err = stream.Send(pubReq)
			if err != nil {
				logrus.WithError(err).Errorln("send pub error", pubReq)
				break Loop
			}
		}
	}
	return
}

// Ping test grpc server health
func (s *Server) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

// run create go func to monitor pubChan
func (s *Server) run() {
Loop:
	for {
		select {
		case <-s.quit:
			logrus.Warningln("pub chan closed")
			break Loop
		case pubReq, ok := <-s.pubChan:
			if !ok {
				break Loop
			}
			logrus.Infoln("Receive PubReq", pubReq)
		}
	}
}
