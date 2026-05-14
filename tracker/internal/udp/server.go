package udp

import (
	"context"
	"net"
	"time"

	"github.com/vpt/common/log"
	"github.com/vpt/tracker/internal/service"
)

// Server implements BEP 15 UDP tracker protocol.
// TODO: parse connect / announce / scrape packets and respond per spec.
type Server struct {
	addr    string
	tracker *service.Tracker
}

func NewServer(addr string, t *service.Tracker) *Server {
	return &Server{addr: addr, tracker: t}
}

func (s *Server) Run(ctx context.Context) {
	pc, err := net.ListenPacket("udp", s.addr)
	if err != nil {
		log.Error("udp listen", "err", err)
		return
	}
	defer pc.Close()
	log.Info("tracker udp listening", "addr", s.addr)

	buf := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		_ = pc.SetReadDeadline(time.Now().Add(time.Second))
		n, src, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go s.handle(ctx, pc, src, append([]byte(nil), buf[:n]...))
	}
}

func (s *Server) handle(ctx context.Context, pc net.PacketConn, src net.Addr, packet []byte) {
	// TODO: implement BEP 15:
	//   - connect (16 bytes): respond with connection_id
	//   - announce (98+ bytes): build model.AnnounceRequest, call s.tracker.Announce
	//   - scrape: call s.tracker.Scrape
	_ = ctx
	_ = pc
	_ = src
	_ = packet
}
