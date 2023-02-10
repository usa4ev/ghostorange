package storage

import "ghostorange/internal/app/storage/mock"

type(
	Storage interface{
		GetData(dataType int)(any,error)
	}
)

func New()Storage{
	return mock.New()
}