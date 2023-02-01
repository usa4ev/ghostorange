package provider

import (
	"ghostorange/internal/app/model"
	"ghostorange/internal/app/provider/httpp"
)

type (
	config interface{
		SrvAddr() string
	}
	Provider interface {
		Login(model.Credentials) bool

		Count(dataType int) (int, error)

		GetData(dataType int) (any, error)
	}
)

func New(cfg config)(Provider, error){
	return httpp.New(cfg)
}