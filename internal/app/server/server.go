package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"ghostorange/internal/app/auth"
	"ghostorange/internal/app/router"
	"ghostorange/internal/app/server/middleware"
	"ghostorange/internal/app/storage"
)

type (
	Server struct {
		httpsrv  *http.Server
		cfg      config
		usrStrg  auth.UsrStorage
		dataStrg storage.Storage
	}

	config interface {
		SrvAddr() string
		DBDSN() string
		SessionLifetime() time.Duration
	}
)

func New(c config, s storage.Storage) *Server {
	srv := Server{cfg: c,
		usrStrg:  s,
		dataStrg: s}
	r := router.NewRouter(&srv)
	srv.httpsrv = &http.Server{Addr: c.SrvAddr(), Handler: r}

	return &srv
}

func (srv *Server) Handlers() []router.HandlerDesc {
	return []router.HandlerDesc{
		// POST: /users/register
		{Method: "POST",
			Path:    "/v1/users/register",
			Handler: http.HandlerFunc(srv.Register),
			Middlewares: chi.Middlewares{
				middleware.GzipMW},
		},

		// POST: /users/login
		{Method: "POST",
			Path:    "/v1/users/login",
			Handler: http.HandlerFunc(srv.Login),
			Middlewares: chi.Middlewares{
				middleware.GzipMW},
		},

		// GET: /data?data_type={data_type}
		{Method: "GET",
			Path:    "/v1/data",
			Handler: http.HandlerFunc(srv.GetData),
			Middlewares: chi.Middlewares{
				middleware.GzipMW,
				middleware.AuthorisationMW},
		},

		// POST: /data?data_type={data_type}
		{Method: "POST",
			Path:    "/v1/data",
			Handler: http.HandlerFunc(srv.AddData),
			Middlewares: chi.Middlewares{
				middleware.GzipMW,
				middleware.AuthorisationMW},
		},

		// PUT: /data?data_type={data_type}
		{Method: "PUT",
			Path:    "/v1/data",
			Handler: http.HandlerFunc(srv.UpdateData),
			Middlewares: chi.Middlewares{
				middleware.GzipMW,
				middleware.AuthorisationMW},
		},

		// GET: /data/count?data_type={data_type}
		{Method: "GET",
			Path:    "/v1/data/count",
			Handler: http.HandlerFunc(srv.Count),
			Middlewares: chi.Middlewares{
				middleware.GzipMW,
				middleware.AuthorisationMW},
		},
	}
}

func (srv *Server) Run() error {
	// Run the server
	// if srv.cfg.UseTLS() {
	// 	return srv.httpsrv.ListenAndServeTLS(
	// 		filepath.Join(srv.cfg.SslPath(), "example.crt"),
	// 		filepath.Join(srv.cfg.SslPath(), "example.key"))
	// } else {
	return srv.httpsrv.ListenAndServe()
	// }
}

func (srv *Server) Shutdown(ctx context.Context) error {
	srv.httpsrv.Shutdown(ctx)
	return nil
}
