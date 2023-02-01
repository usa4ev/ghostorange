// Package model contains entities that describe data model
// and provides helpful functions to develop consistent approach
// to different data objects.
package model

import (
	"bytes"
	"encoding/json"
)

const (
	// stored data types
	KeyCredentials = iota
	KeyText
	KeyBinary
	KeyCard
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
		ID      string `json:"id"`
		Data    string `json:"data"`
		Name    string `json:"name"`
		Comment string `json:"comment"`
	}

	ItemCard struct {
		ID                 string `json:"id"`
		Number             string `json:"number"`
		Exp                string `json:"expiration_date"`
		CardholderName     string `json:"holder_name"`
		CardholderSurename string `json:"holder_surename"`
		CVVHash            string `json:"cvv_hash"`
		Name               string `json:"name"`
		Comment            string `json:"comment"`
	}

	Item interface{
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
	case KeyCard:
		return "Card info"
	}

	return ""
}

func EncodeItemsJSON(data any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)

	if err := enc.Encode(data); err != nil{
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeItemsJSON[T Item](dataType int, data []byte) ([]T, error) {
	buf := bytes.NewBuffer(data)
	dec := json.NewDecoder(buf)

	res := make([]T,0)

	err := dec.Decode(&res)

	return res, err
}
