package httpp

import (
	"fmt"
	"io"
	"net/http"

	"ghostorange/internal/app/model"
)

type (
	Provider struct {
		client *http.Client
		cfg    config
	}
	config interface {
		SrvAddr() string
	}
)

func New(cfg config) (*Provider, error) {
	return &Provider{
		client: http.DefaultClient,
		cfg: cfg},
		nil
}

func (prov Provider) Count(dataType int) (int, error){
 return 0, nil
}

func (prov Provider) Login(model.Credentials) bool{
	return false
}

func (prov Provider) GetData(dataType int) (any, error) {
	req,err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://%v/v1/data?data_type=%v",
			prov.cfg.SrvAddr(), dataType),
		nil)	
	if err != nil {
		return nil, fmt.Errorf("failed to compose GetData request: %w", err)
	}

	res, err := prov.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("GetData request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read server GetData response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`server returned unexpected code: %v\n
			response: %v`,
			res.StatusCode, message)
	}

	switch dataType{
	case model.KeyCredentials:
		return model.DecodeItemsJSON[model.ItemCredentials](dataType, message)
	case model.KeyText:
		return model.DecodeItemsJSON[model.ItemText](dataType, message)
	case model.KeyBinary:
		return model.DecodeItemsJSON[model.ItemBinary](dataType, message)
	case model.KeyCard:
		return model.DecodeItemsJSON[model.ItemCard](dataType, message)
	}

	return nil, fmt.Errorf("unknown data-type: %v", dataType)
}
