package server_test

import (
	"testing"
	"time"

	"github.com/kaisawind/message/server"
)

func TestNewServer(t *testing.T) {
	s := server.NewServer()
	err := s.Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestServer_Serve(t *testing.T) {
	s := server.NewServer()

	go func() {
		err := s.Serve()
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}()
	time.Sleep(100 * time.Millisecond)
	err := s.Close()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
