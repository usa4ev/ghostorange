package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"

	"github.com/usa4ev/ghostorange/internal/app/auth"
	"github.com/usa4ev/ghostorange/internal/app/router"
	"github.com/usa4ev/ghostorange/internal/app/server/middleware"
	"github.com/usa4ev/ghostorange/internal/app/storage"
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
			Path:        "/v1/users/register",
			Handler:     http.HandlerFunc(srv.Register),
			Middlewares: nil,
		},

		// POST: /users/login
		{Method: "POST",
			Path:        "/v1/users/login",
			Handler:     http.HandlerFunc(srv.Login),
			Middlewares: nil,
		},

		// GET: /data?data_type={data_type}
		{Method: "GET",
			Path:    "/v1/data",
			Handler: http.HandlerFunc(srv.GetData),
			Middlewares: chi.Middlewares{
				chimw.Compress(5, CTJSON),
				middleware.AuthorisationMW},
		},

		// POST: /data?data_type={data_type}
		{Method: "POST",
			Path:    "/v1/data",
			Handler: http.HandlerFunc(srv.AddData),
			Middlewares: chi.Middlewares{
				chimw.Compress(5, CTJSON),
				middleware.AuthorisationMW},
		},

		// PUT: /data?data_type={data_type}
		{Method: "PUT",
			Path:    "/v1/data",
			Handler: http.HandlerFunc(srv.AddData),
			Middlewares: chi.Middlewares{
				chimw.Compress(5, CTJSON),
				middleware.AuthorisationMW},
		},

		// GET: /data/count?data_type={data_type}
		{Method: "GET",
			Path:    "/v1/data/count",
			Handler: http.HandlerFunc(srv.Count),
			Middlewares: chi.Middlewares{
				chimw.Compress(5, CTJSON),
				middleware.AuthorisationMW},
		},

		// GET: /v1/data/cards/{id}
		{Method: "GET",
			Path:    "/v1/data/cards/{id}",
			Handler: http.HandlerFunc(srv.CardData),
			Middlewares: chi.Middlewares{
				chimw.Compress(5, CTJSON),
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
