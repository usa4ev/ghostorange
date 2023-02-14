// Package model contains entities that describe data model
// and provides helpful functions to develop consistent approach
// to different data objects.
package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

const (
	// stored data types
	KeyCredentials = iota
	KeyText
	KeyBinary
	KeyCards
	KeyLimit
)

type (
	Credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	UserInfo struct {
		Credentials Credentials `json:"credentials"`
		Email       string      `json:"email"`
	}

	// stored data items
	ItemCredentials struct {
		ID          string      `json:"id"`
		Credentials Credentials `json:"credentials"`
		Name        string      `json:"name"`
		Comment     string      `json:"comment"`
	}

	ItemText struct {
		ID      string `json:"id"`
		Text    string `json:"text"`
		Name    string `json:"name"`
		Comment string `json:"comment"`
	}

	ItemBinary struct {
		ID        string `json:"id"`
		Size      int    `json:"size"`
		Extention string `json:"extention"`
		Data      string `json:"data"`
		Name      string `json:"name"`
		Comment   string `json:"comment"`
	}

	ItemCard struct {
		ID                 string    `json:"id"`
		Number             string    `json:"number"`
		Exp                time.Time `json:"expiration_date"`
		CardholderName     string    `json:"holder_name"`
		CardholderSurename string    `json:"holder_surename"`
		CVVHash            string    `json:"cvv_hash"`
		Name               string    `json:"name"`
		Comment            string    `json:"comment"`
	}

	Item interface {
		ItemCredentials |
			ItemText |
			ItemBinary |
			ItemCard
	}
)

func GetItemTitle(dataType int) string {
	switch dataType {
	case KeyCredentials:
		return "Credentials"
	case KeyText:
		return "Text data"
	case KeyBinary:
		return "Binary data"
	case KeyCards:
		return "Card info"
	}

	return ""
}

func EncodeItemsJSON(data any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)

	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeItemsJSON(dataType int, message []byte) (any, error) {
	switch dataType {
	case KeyCredentials:
		return decodeItemsJSON[ItemCredentials](dataType, message)
	case KeyText:
		return decodeItemsJSON[ItemText](dataType, message)
	case KeyBinary:
		return decodeItemsJSON[ItemBinary](dataType, message)
	case KeyCards:
		return decodeItemsJSON[ItemCard](dataType, message)
	}

	return nil, fmt.Errorf("unsupported data type")
}

func decodeItemsJSON[T Item](dataType int, data []byte) ([]T, error) {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)

	res := make([]T, 0)

	err := dec.Decode(&res)

	return res, err
}

func DecodeItemJSON(dataType int, message []byte) (any, error) {
	switch dataType {
	case KeyCredentials:
		return decodeItemJSON[ItemCredentials](message)
	case KeyText:
		return decodeItemJSON[ItemText](message)
	case KeyBinary:
		return decodeItemJSON[ItemBinary](message)
	case KeyCards:
		return decodeItemJSON[ItemCard](message)
	}

	return nil, fmt.Errorf("unsupported data type")
}

func decodeItemJSON[T Item](data []byte) (T, error) {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)

	var res T

	err := dec.Decode(&res)

	return res, err
}
