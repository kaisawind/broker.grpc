package server

import (
	"net"
	"sync"

	pb "github.com/kaisawind/message/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server server
type Server struct {
	pb.UnimplementedMessageServer
	grpcServer *grpc.Server
	pubChan    chan *pb.PubReq
	once       sync.Once
	quit       chan bool

	GRPCHost string `long:"grpc-host" description:"the IP to listen on for grpc" default:"0.0.0.0" env:"GRPC_HOST"`
	GRPCPort string `long:"grpc-port" description:"the port to listen on for grpc's insecure connections" default:"6653" env:"GRPC_PORT"`
}

// NewServer ...
func NewServer() *Server {
	server := &Server{
		pubChan: make(chan *pb.PubReq, 100),
		quit:    make(chan bool),
	}
	return server
}

// Close clean all data if need
func (s *Server) Close() (err error) {
	s.once.Do(func() {
		s.quit <- true
		close(s.pubChan)
	})
	return
}

// Serve ...
func (s *Server) Serve() (err error) {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.grpcServe()
	}()
	wg.Wait()

	logrus.Infoln("Server Exited ...")
	return
}

func (s *Server) grpcServe() error {
	listener, err := net.Listen("tcp", net.JoinHostPort(s.GRPCHost, s.GRPCPort))
	if err != nil {
		logrus.Fatalf("failed to grpc listen: %v", err)
		return err
	}
	opts := []grpc.ServerOption{
		// StatsHandler returns a ServerOption that sets the stats handler for the server.
		grpc.StatsHandler(&StatsHandler{}),
	}
	s.grpcServer = grpc.NewServer(opts...)
	pb.RegisterMessageServer(s.grpcServer, s)

	// grpc cli
	reflection.Register(s.grpcServer)

	logrus.Infoln("grpc service is started ...  addr:", listener.Addr().String())
	err = s.grpcServer.Serve(listener)
	if err != nil {
		logrus.Fatalf("failed to monitor serve: %v", err)
		return err
	}
	logrus.Infoln("grpc server existed ...")
	return nil
}