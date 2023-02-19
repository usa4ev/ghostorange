package httpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"

	"github.com/usa4ev/ghostorange/internal/app/model"
	"github.com/usa4ev/ghostorange/internal/app/server"
)

type (
	Provider struct {
		client *http.Client
		cfg    config
		logger *zap.SugaredLogger
	}
	config interface {
		SrvAddr() string
	}
)

func New(cfg config, logger *zap.SugaredLogger) (*Provider, error) {
	cl := http.DefaultClient
	jar, err := cookiejar.New(
		&cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		},
	)
	if err != nil {
		return nil, err
	}

	cl.Jar = jar

	return &Provider{
			client: http.DefaultClient,
			cfg:    cfg},
		nil
}

func (prov Provider) Count(dataType int) (string, error) {
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://%v/v1/data/count?data_type=%v",
			prov.cfg.SrvAddr(), dataType),
		nil)
	if err != nil {
		return "", fmt.Errorf("failed to compose GetData request: %w", err)
	}

	res, err := prov.client.Do(req)

	if err != nil {
		return "", fmt.Errorf("GetData request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read server GetData response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	return string(message), nil
}

func (prov Provider) Register(item model.Credentials) error {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(item); err != nil {
		return fmt.Errorf("failed to encode Register message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v/v1/users/register",
			prov.cfg.SrvAddr()),
		buf)

	if err != nil {
		return fmt.Errorf("failed to compose Register request: %w", err)
	}

	req.Header.Set("Content-Type", server.CTJSON)

	res, err := prov.client.Do(req)

	if err != nil {
		return fmt.Errorf("Register request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read server Register response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	return nil
}

func (prov Provider) Login(item model.Credentials) error {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(item); err != nil {
		return fmt.Errorf("failed to encode Register message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v/v1/users/login",
			prov.cfg.SrvAddr()),
		buf)

	if err != nil {
		return fmt.Errorf("failed to compose Login request: %w", err)
	}

	req.Header.Set("Content-Type", server.CTJSON)

	res, err := prov.client.Do(req)

	if err != nil {
		return fmt.Errorf("Login request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read server Login response: %w", err)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("login or password must be wrong")
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	return nil
}

func (prov Provider) GetData(dataType int) (any, error) {
	req, err := http.NewRequest(http.MethodGet,
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
		return nil, fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	obj, err := model.DecodeItemsJSON(dataType, message)
	if err != nil {
		return nil, fmt.Errorf("failed to decode server message: %w", err)
	}

	return obj, nil
}

func (prov *Provider) AddData(dataType int, data any) error {
	msg, err := model.EncodeItemsJSON(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v/v1/data?data_type=%v",
			prov.cfg.SrvAddr(), dataType),
		bytes.NewBuffer(msg))
	if err != nil {
		return fmt.Errorf("failed to compose AddData request: %w", err)
	}

	req.Header.Set("Content-Type", server.CTJSON)

	res, err := prov.client.Do(req)

	if err != nil {
		return fmt.Errorf("AddData request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read server AddData response: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	return nil
}

func (prov *Provider) UpdateData(dataType int, data any) error {
	msg, err := model.EncodeItemsJSON(data)
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut,
		fmt.Sprintf("http://%v/v1/data?data_type=%v",
			prov.cfg.SrvAddr(), dataType),
		bytes.NewBuffer(msg))
	if err != nil {
		return fmt.Errorf("failed to compose UpdateData request: %w", err)
	}

	req.Header.Set("Content-Type", server.CTJSON)

	res, err := prov.client.Do(req)

	if err != nil {
		return fmt.Errorf("UpdateData request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read server UpdateData response: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	return nil
}

func (prov *Provider) GetCard(id, cvv string) (model.ItemCard, error) {
	var item model.ItemCard

	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("http://%v/v1/data/cards/%v",
			prov.cfg.SrvAddr(), id),
		bytes.NewBuffer([]byte(cvv)))
	if err != nil {
		return item, fmt.Errorf("failed to compose GetCard request: %w", err)
	}

	req.Header.Set("Content-Type", server.CTPlain)

	res, err := prov.client.Do(req)

	if err != nil {
		return item, fmt.Errorf("GetCard request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return item, fmt.Errorf("failed to read server GetCard response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return item, fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	val, err := model.DecodeItemJSON(model.KeyCards, message)
	if err != nil {
		return item, fmt.Errorf("failed to decode server message: %w", err)
	}

	var ok bool
	if item, ok = val.(model.ItemCard); !ok {
		return item, fmt.Errorf("bad type decoded, expected model.ItemCard")
	}

	return item, nil
}

func (prov *Provider) Lg() *zap.SugaredLogger {
	return prov.logger
}
