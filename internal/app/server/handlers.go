package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"ghostorange/internal/app/auth"
	"ghostorange/internal/app/auth/session"
	"ghostorange/internal/app/model"
)

const (
	CTJSON  = "application/json"
	CTPlain = "plain/text"
)

// Register handler adds new user if one does not exist and opens a new session
func (srv *Server) Count(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(session.CtxKeyUserID).(string)
	if !ok {
		http.Error(w, "request context is missing user ID", http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()

	strDataType := r.URL.Query().Get("data_type")
	if strDataType == "" {
		http.Error(w, "data_type parameter is required", http.StatusBadRequest)
	}

	dataType, err := strconv.Atoi(strDataType)
	if err != nil || dataType >= model.KeyLimit {
		http.Error(w, "bad data_type parameter", http.StatusBadRequest)
	}

	res, err := srv.dataStrg.Count(r.Context(), dataType, userID)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to get data from storage: %v",
				err.Error()),
			http.StatusInternalServerError)
	}

	w.Header().Set("content-type", CTJSON)

	w.Write([]byte(strconv.Itoa(res)))
}

// Register handler adds new user if one does not exist and opens a new session
func (srv *Server) Register(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	if ct != "" && ct != CTJSON {
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
			Path:    "/v1",
		})

	http.SetCookie(w,
		&http.Cookie{
			Name:    "UserID",
			Value:   userID,
			Expires: expiresAt,
			Path:    "/v1",
		})

}

// Register handler adds a new user if one does not exist and opens a new session
func (srv *Server) Login(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("content-type")
	if ct != "" && strings.Compare(ct, CTJSON) == 0 {
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

	http.SetCookie(w,
		&http.Cookie{
			Name:    "Authorization",
			Value:   token,
			Expires: expiresAt,
			Path:    "/v1",
		})

	http.SetCookie(w,
		&http.Cookie{
			Name:    "UserID",
			Value:   userID,
			Expires: expiresAt,
			Path:    "/v1",
		})
}

// GetData responds with JSON encoded array of ojects,
// type depending on data_type query parameter.
func (srv *Server) GetData(w http.ResponseWriter, r *http.Request) {
	strDataType := r.URL.Query().Get("data_type")
	if strDataType == "" {
		http.Error(w, "data_type parameter is required", http.StatusBadRequest)

		return
	}

	dataType, err := strconv.Atoi(strDataType)
	if err != nil || dataType >= model.KeyLimit {
		http.Error(w, "bad data_type parameter", http.StatusBadRequest)

		return
	}

	data, err := srv.dataStrg.GetData(r.Context(), dataType)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to get data from storage: %v",
				err.Error()),
			http.StatusInternalServerError)

		return
	}

	res, err := model.EncodeItemsJSON(data)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to encode data: %v",
				err.Error()),
			http.StatusInternalServerError)

		return
	}

	w.Header().Set("content-type", CTJSON)
	w.Write(res)
}

// AddData adds new object to storage.
func (srv *Server) AddData(w http.ResponseWriter, r *http.Request) {
	strDataType := r.URL.Query().Get("data_type")
	if strDataType == "" {
		http.Error(w, "data_type parameter is required", http.StatusBadRequest)

		return
	}

	dataType, err := strconv.Atoi(strDataType)
	if err != nil || dataType >= model.KeyLimit {
		http.Error(w, "bad data_type parameter", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(session.CtxKeyUserID).(string)
	if !ok {
		http.Error(w, "context is missing user ID", http.StatusInternalServerError)
		
		return
	}

	defer r.Body.Close()
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to read message: %v", err.Error()),
			http.StatusBadRequest)

		return
	}

	obj, err := model.DecodeItemJSON(dataType, msg)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to decode JSON: %v", err.Error()),
			http.StatusInternalServerError)

		return
	}

	err = srv.dataStrg.AddData(r.Context(), dataType, userID, obj)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to store data: %v",
				err.Error()),
			http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("content-type", CTJSON)
}

// UpdateData updates new object to storage.
func (srv *Server) UpdateData(w http.ResponseWriter, r *http.Request) {
	strDataType := r.URL.Query().Get("data_type")
	if strDataType == "" {
		http.Error(w, "data_type parameter is required", http.StatusBadRequest)

		return
	}

	dataType, err := strconv.Atoi(strDataType)
	if err != nil || dataType >= model.KeyLimit {
		http.Error(w, "bad data_type parameter", http.StatusBadRequest)

		return
	}

	userID, ok := r.Context().Value(session.CtxKeyUserID).(string)
	if !ok {
		http.Error(w, "context is missing user ID", http.StatusInternalServerError)
		
		return
	}

	defer r.Body.Close()
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to read message: %v", err.Error()),
			http.StatusBadRequest)

		return
	}

	obj, err := model.DecodeItemJSON(dataType, msg)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to decode JSON: %v", err.Error()),
			http.StatusInternalServerError)

		return
	}

	err = srv.dataStrg.UpdateData(r.Context(), dataType, userID, obj)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to update data: %v",
				err.Error()),
			http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("content-type", CTJSON)
}
