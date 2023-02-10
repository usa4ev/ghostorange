package mock

import (
	"fmt"

	"ghostorange/internal/app/model"
)

type (
	Storage struct {
	}
)

func New()*Storage{
	return &Storage{}
}

func (s *Storage) GetData(dataType int) (any, error) {
	switch dataType {
	case model.KeyCredentials:
		return []model.ItemCredentials{
			{Credentials: model.Credentials{
				Login:    "testlogin",
				Password: "testpassword",
			},
				Comment: "confidential",
				Name:    "My login",
			},
			{Credentials: model.Credentials{
				Login:    "tanya",
				Password: "dragon",
			},
				Comment: "lol",
				Name:    "Tanya's login",
			},
		}, nil
	}

	return nil, fmt.Errorf("wrong data-type %v", dataType)
}
