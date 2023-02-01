package mock

import (
	"fmt"

	"ghostorange/internal/app/model"
)

type provider struct{}

func New() *provider {
	return &provider{}
}

func (p *provider) Login(model.Credentials) bool {
	return true
}

func (p *provider) Count(dataType int) (int, error) {
	return 0, nil
}

func (p *provider) GetData(dataType int) (any, error) {
	switch dataType {
	case model.KeyCredentials:
		return []model.ItemCredentials{
			{Credentials: model.Credentials{
					Login:"testlogin", 
					Password: "testpassword",
				},
				Comment: "confidential",
				Name: "My login",
			},
			{Credentials: model.Credentials{
				Login: "tanya", 
				Password: "dragon",
			},
				Comment: "lol",
				Name: "Tanya's login",
			},
		}, nil
	}
	return nil, fmt.Errorf("unknown data type")
}
