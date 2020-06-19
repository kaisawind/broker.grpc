package server

import (
	"context"
	pb "github.com/kaisawind/message/pb"
	"github.com/sirupsen/logrus"
)

// Publish publish topic message
func (s *Server) Publish(_ context.Context, req *pb.PubReq) (resp *pb.PubResp, err error) {
	select {
	case <-s.quit:
	case s.pubChan <- req:
	default:
		logrus.Warningln("lose pub req", req)
	}
	return
}

// Run create go func to monitor pubChan
func (s *Server) Run() {
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
			logrus.Infoln("receive pub req", pubReq)
		}
	}
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
