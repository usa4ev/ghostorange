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

	"ghostorange/internal/app/model"
	"ghostorange/internal/app/server"
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
	if err != nil{
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

func (prov Provider) Register(item model.Credentials) (error) {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(item); err != nil{
		return fmt.Errorf("failed to encode Register message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v/v1/users/register",
			prov.cfg.SrvAddr()),
			buf)

	if err != nil {
		return fmt.Errorf("failed to compose Register request: %w", err)
	}

	req.Header.Set("content-type", server.CTJSON)

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

	return  nil
}

func (prov Provider) Login(item model.Credentials) error {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(item); err != nil{
		return fmt.Errorf("failed to encode Register message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("http://%v/v1/users/login",
			prov.cfg.SrvAddr()),
		buf)

	if err != nil {
		return fmt.Errorf("failed to compose Login request: %w", err)
	}

	res, err := prov.client.Do(req)

	if err != nil {
		return fmt.Errorf("Login request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read server Login response: %w", err)
	}

	if res.StatusCode ==  http.StatusUnauthorized{
		return fmt.Errorf("login or password must be wrong")
	}else if res.StatusCode != http.StatusOK {
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
		return nil, fmt.Errorf("failed to decode server message: %v", dataType)
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

	res, err := prov.client.Do(req)

	if err != nil {
		return fmt.Errorf("AddData request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read server AppData response: %w", err)
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

	res, err := prov.client.Do(req)

	if err != nil {
		return fmt.Errorf("UpdateData request failed: %w", err)
	}

	defer res.Body.Close()

	message, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read server AppData response: %w", err)
	}

	if res.StatusCode != http.StatusAccepted {
		return fmt.Errorf(`server returned unexpected code: %v 
			response: %v`,
			res.StatusCode, string(message))
	}

	return nil
}

func (prov *Provider) Lg() *zap.SugaredLogger {
	return prov.logger
}
