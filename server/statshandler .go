package server

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/stats"
)

// StatsHandler Handler defines the interface for the related stats handling (e.g., RPCs, connections).
type StatsHandler struct{}

// TagRPC can attach some information to the given context.
// The context used for the rest lifetime of the RPC will be derived from
// the returned context.
func (h *StatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	logrus.Debugln("Tag RPC", info)
	return ctx
}

// HandleRPC processes the RPC stats.
func (h *StatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	// logrus.Debugln("RPCStats:", s)
	switch s.(type) {
	case *stats.Begin:
		begin := s.(*stats.Begin)
		logrus.Debugln("Begin:", begin.BeginTime.Format(time.RFC3339))
	case *stats.InPayload:
		inPayload := s.(*stats.InPayload)
		logrus.Debugln("InPayload:", inPayload.RecvTime.Format(time.RFC3339), "Length:", inPayload.Length)
	case *stats.InHeader:
		inHeader := s.(*stats.InHeader)
		logrus.Debugln("InHeader:", inHeader.FullMethod, "RemoteAddr:", inHeader.RemoteAddr, "LocalAddr:", inHeader.LocalAddr)
	case *stats.InTrailer:
		inTrailer := s.(*stats.InTrailer)
		logrus.Debugln("InTrailer:", inTrailer.WireLength)
	case *stats.OutPayload:
		outPayload := s.(*stats.OutPayload)
		logrus.Debugln("OutPayload:", outPayload.SentTime.Format(time.RFC3339))
	case *stats.OutHeader:
		outHeader := s.(*stats.OutHeader)
		logrus.Debugln("OutHeader:", outHeader.FullMethod, "RemoteAddr:", outHeader.RemoteAddr, "LocalAddr:", outHeader.LocalAddr)
	case *stats.OutTrailer:
		outTrailer := s.(*stats.OutTrailer)
		logrus.Debugln("OutTrailer:", outTrailer.Trailer)
	case *stats.End:
		end := s.(*stats.End)
		logrus.Debugln("End:", end.BeginTime.Format(time.RFC3339), end.EndTime.Format(time.RFC3339), "Error:", end.Error)
	default:
		logrus.Debugln("illegal RPCStats type")
	}
}

// TagConn can attach some information to the given context.
// The returned context will be used for stats handling.
// For conn stats handling, the context used in HandleConn for this
// connection will be derived from the context returned.
// For RPC stats handling,
//  - On server side, the context used in HandleRPC for all RPCs on this
// connection will be derived from the context returned.
//  - On client side, the context is not derived from the context returned.
func (h *StatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	logrus.Infoln("Tag Conn", info)
	return ctx
}

// HandleConn processes the Conn stats.
func (h *StatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	switch s.(type) {
	case *stats.ConnBegin:
		logrus.Infoln("Conn Begin")
	case *stats.ConnEnd:
		logrus.Infoln("Conn End ")
	default:
		logrus.Infoln("illegal ConnStats type")
	}
}
