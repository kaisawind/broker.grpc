package main

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/kaisawind/message/server"
	"github.com/sirupsen/logrus"
	"os"
)

// inject by go build
var (
	Version   = "0.0.0"
	BuildTime = "2020-01-13-0802 UTC"
)

func init() {
	logrus.Infoln("Version:", Version)
	logrus.Infoln("BuildTime:", BuildTime)
}

func main() {
	s := server.NewServer()

	parser := flags.NewParser(s, flags.Default)
	parser.ShortDescription = "grpc broker"
	parser.LongDescription = "This is a grpc broker with pub and sub."

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		logrus.WithError(err).Errorln("flags Parse Error")
		os.Exit(code)
	}

	logrus.Infoln("grpc broker service is starting ...")
	if err := s.Serve(); err != nil {
		logrus.Fatalln(err)
	}
}
