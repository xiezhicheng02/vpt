package udp

import (
	"context"
	"net"
	"time"

	"github.com/vpt/common/log"
	"github.com/vpt/gateway/internal/proxy"
)

// Server is a minimal UDP relay that forwards packets to the tracker service.
// BT UDP tracker protocol (BEP 15) framing is parsed downstream by tracker.
type Server struct {
	addr   string
	router *proxy.Router
}

func NewServer(addr string, router *proxy.Router) *Server {
	return &Server{addr: addr, router: router}
}

func (s *Server) Run(ctx context.Context) {
	pc, err := net.ListenPacket("udp", s.addr)
	if err != nil {
		log.Error("udp listen", "err", err)
		return
	}
	defer pc.Close()
	log.Info("gateway udp listening", "addr", s.addr)

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
	upstream, err := s.router.PickUpstream(ctx, "tracker")
	if err != nil {
		log.Warn("udp upstream", "err", err)
		return
	}
	// upstream URL host carries host:port; tracker also exposes UDP on the same host but a configured UDP port.
	// For simplicity we expect tracker UDP listener at the same host on a fixed port (config-driven later).
	host, _, _ := net.SplitHostPort(upstream.Host)
	dst := host + ":8003"

	conn, err := net.Dial("udp", dst)
	if err != nil {
		log.Warn("udp dial", "err", err)
		return
	}
	defer conn.Close()
	if _, err := conn.Write(packet); err != nil {
		return
	}
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	resp := make([]byte, 2048)
	n, err := conn.Read(resp)
	if err != nil {
		return
	}
	_, _ = pc.WriteTo(resp[:n], src)
}
