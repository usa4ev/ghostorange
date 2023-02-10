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
	server struct {
		httpsrv *http.Server
		cfg     config
		usrStrg auth.UsrStorage
		dataStrg storage.Storage
	}

	config interface {
		SrvAddr() string
		// UseTLS() bool
		// SslPath() string
		SessionLifetime() time.Duration
	}
)

func New(c config) *server {
	strg := storage.New()

	srv := server{cfg: c,
		usrStrg: nil, // ToDo: implement methods
		dataStrg: strg,}
	r := router.NewRouter(&srv)
	srv.httpsrv = &http.Server{Addr: c.SrvAddr(), Handler: r}

	return &srv
}

func (srv *server) Handlers() []router.HandlerDesc {
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
				middleware.GzipMW},
		},

		// {Method: "POST",
		// Path: "/",
		// Handler: http.HandlerFunc(srv.makeShort),
		// Middlewares: chi.Middlewares{
		// 	middleware.GzipMW,
		// 	middleware.AuthMW(sm)}},
	}
}

func (srv *server) Run() error {
	// Run the server
	// if srv.cfg.UseTLS() {
	// 	return srv.httpsrv.ListenAndServeTLS(
	// 		filepath.Join(srv.cfg.SslPath(), "example.crt"),
	// 		filepath.Join(srv.cfg.SslPath(), "example.key"))
	// } else {
		return srv.httpsrv.ListenAndServe()
	// }
}

func (srv *server) Shutdown(ctx context.Context) error {
	srv.httpsrv.Shutdown(ctx)
	return nil
}
