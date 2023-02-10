package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"ghostorange/internal/app/auth"
	"ghostorange/internal/app/auth/session"
	"ghostorange/internal/app/model"
)

const ctJSON = "application/json"

// Register handler adds new user if one does not exist and opens a new session
func (srv *server) Register(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "" && ct != ctJSON {
		http.Error(w, fmt.Sprintf("unexpected content-type %v", ct), http.StatusBadRequest)

		return
	}

	defer r.Body.Close()

	cred := model.Credentials{}

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&cred); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode a message: %v", err), http.StatusBadRequest)

		return
	}

	userID, err := auth.RegisterUser(r.Context(), cred.Login, cred.Password, srv.usrStrg)
	if errors.Is(err, auth.ErrUserAlreadyExists) {
		http.Error(w, err.Error(), http.StatusConflict)

		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)

		return
	}

	token, expiresAt, err := session.Open(userID, srv.cfg.SessionLifetime())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to open new session: %v", err), http.StatusInternalServerError)

		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "Authorization",
			Value:   token,
			Expires: expiresAt,
		})
}

// Register handler adds a new user if one does not exist and opens a new session
func (srv *server) Login(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct != "" && strings.Compare(ct, ctJSON) == 0 {
		http.Error(w, fmt.Sprintf("unexpected content-type %v", ct), http.StatusBadRequest)

		return
	}

	defer r.Body.Close()

	cred := model.Credentials{}

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&cred); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode a message: %v", err), http.StatusBadRequest)

		return
	}

	userID, err := auth.Login(r.Context(), cred.Login, cred.Password, srv.usrStrg)

	if errors.Is(err, auth.ErrUnathorized) {
		http.Error(w, err.Error(), http.StatusUnauthorized)

		return
	} else if err != nil {
		http.Error(w, fmt.Sprintf("authentication failed: %v", err), http.StatusInternalServerError)

		return
	}

	token, expiresAt, err := session.Open(userID, srv.cfg.SessionLifetime())

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to open new session: %v", err), http.StatusInternalServerError)

		return
	}

	http.SetCookie(w, &http.Cookie{Name: "Authorization", Value: token, Expires: expiresAt})
}

// Register GetData responds with JSON encoded array of ojects, 
// type depending on data_type query parameter.
func (srv *server) GetData(w http.ResponseWriter, r *http.Request) {
	strDataType := r.URL.Query().Get("data_type")
	if strDataType == "" {
		http.Error(w, "data_type parameter is required", http.StatusBadRequest)
	}

	dataType, err := strconv.Atoi(strDataType)
	if err != nil || dataType >= model.KeyLimit {
		http.Error(w, "bad data_type parameter", http.StatusBadRequest)
	}

	data, err := srv.dataStrg.GetData(dataType)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to get data fropm storage: %v",
				err.Error()),
			http.StatusInternalServerError)
	}

	res, err := model.EncodeItemsJSON(data)
	if err != nil{
		http.Error(w,
			fmt.Sprintf("failed to encode data: %v",
				err.Error()),
			http.StatusInternalServerError)
	}

	w.Header().Set("content-type", ctJSON)
	w.Write(res)
}
