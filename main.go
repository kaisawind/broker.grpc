package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/kaisawind/broker.grpc/server"
	"github.com/sirupsen/logrus"
)

// inject by go build
var (
	Version   = "0.0.0"
	BuildTime = "2020-01-13-0802 UTC"
)

func main() {
	s := server.NewServer()

	parser := flags.NewParser(s, flags.Default)
	parser.ShortDescription = "grpc broker"
	parser.LongDescription = "This is a grpc broker with pub and sub."

	buildInfo := initHelp()
	_, err := parser.AddGroup("Build Info", "Show Build Info", &buildInfo)
	if err != nil {
		logrus.WithError(err).Errorln("parser add group error")
	}
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		} else {
			logrus.WithError(err).Errorln("parse flags error")
		}
		os.Exit(code)
	}

	logrus.Infoln("grpc broker service is starting ...")
	if err := s.Serve(); err != nil {
		logrus.Fatalln(err)
	}
}

// buildInfo 编译信息
type buildInfo struct {
	ShowVersion   func() error `short:"v" long:"version" description:"Show this version"`
	ShowBuildTime func() error `short:"b" long:"build" description:"Show this build time"`
}

func initHelp() buildInfo {
	return buildInfo{
		ShowVersion: func() error {
			fmt.Print(Version)
			return &flags.Error{
				Type: flags.ErrHelp,
			}
		},
		ShowBuildTime: func() error {
			fmt.Print(BuildTime)
			return &flags.Error{
				Type: flags.ErrHelp,
			}
		},
	}
}
