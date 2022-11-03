package server

import (
	"context"
	"errors"
	"github.com/aapelismith/kuntunnel/pkg/apiserver/config"
	"github.com/aapelismith/kuntunnel/pkg/log"
	"net/http"
)

type PeerServer struct {
	srv *http.Server
	cfg *config.Configuration
}

//Start the peer http server
func (s *PeerServer) Start(baseCtx context.Context) error {
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	l := log.FromContext(ctx).Sugar()
	l.Infof("http listen and server tls at %s", s.srv.Addr)

	if s.cfg.Peer.CertificateFile != "" {
		if err := s.srv.ListenAndServeTLS(
			s.cfg.Peer.CertificateFile,
			s.cfg.Peer.CertificateKeyFile,
		); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Errorln(err)
			return err
		}
		return nil
	}

	if err := s.srv.ListenAndServe(); err != nil {
		l.Errorln(err)
		return err
	}
	return nil
}

//GracefulStop graceful shutdown http server
func (s *PeerServer) GracefulStop(ctx context.Context) error {
	l := log.FromContext(ctx).Sugar()

	l.Infoln("stopping the http server gracefully..")
	if err := s.srv.Shutdown(ctx); err != nil {
		l.Errorln(err)
		return err
	}
	return nil
}
