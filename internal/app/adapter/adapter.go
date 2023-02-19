package adapter

import (
	"go.uber.org/zap"

	"github.com/usa4ev/ghostorange/internal/app/adapter/httpp"
	"github.com/usa4ev/ghostorange/internal/app/model"
)

type (
	config interface {
		SrvAddr() string
	}
	Adapter interface {
		Login(model.Credentials) error
		Register(model.Credentials) error

		Count(dataType int) (string, error)

		GetData(dataType int) (any, error)
		AddData(dataType int, data any) error
		UpdateData(dataType int, data any) error
		GetCard(id, cvvHash string) (model.ItemCard, error)

		Lg() *zap.SugaredLogger
	}
)

func New(cfg config, logger *zap.SugaredLogger) (Adapter, error) {
	return httpp.New(cfg, logger)
}
