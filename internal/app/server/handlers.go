package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	"github.com/usa4ev/ghostorange/internal/app/auth"
	"github.com/usa4ev/ghostorange/internal/app/auth/session"
	"github.com/usa4ev/ghostorange/internal/app/model"
	"github.com/usa4ev/ghostorange/internal/pkg/argon2hash"
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
	ct := r.Header.Get("Content-Type")
	if ct == "" || strings.Compare(ct, CTJSON) != 0 {
		http.Error(w, fmt.Sprintf("unexpected Content-Type %v", ct), http.StatusBadRequest)

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

// Login handler opens a new session after verifying username and password
func (srv *Server) Login(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if ct == "" || strings.Compare(ct, CTJSON) != 0 {
		http.Error(w, fmt.Sprintf("unexpected Content-Type %v", ct), http.StatusBadRequest)

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

// GetData responds with JSON encoded array of objects,
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

	w.Header().Set("Content-Type", CTJSON)
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
}

// GetData responds with JSON encoded model.ItemCard object
// after verifying CVV code
func (srv *Server) CardData(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "item id if missing in request URL", http.StatusBadRequest)

		return
	}

	userID, ok := r.Context().Value(session.CtxKeyUserID).(string)
	if !ok {
		http.Error(w, "context is missing user ID", http.StatusInternalServerError)

		return
	}

	defer r.Body.Close()

	message, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to read request body: %v",
				err.Error()),
			http.StatusInternalServerError)

		return
	}

	cvv := string(message)
	data, err := srv.dataStrg.GetCardInfo(r.Context(), id, userID)
	if err != nil {
		http.Error(w,
			fmt.Sprintf("failed to get data from storage: %v",
				err.Error()),
			http.StatusInternalServerError)

		return
	}

	if ok, err := argon2hash.ComparePasswordAndHash(cvv, data.CVVHash); err != nil {
		http.Error(w,
			fmt.Sprintf("failed to validate CVV code: %v",
				err.Error()),
			http.StatusInternalServerError)

		return
	} else if !ok {
		http.Error(w,
			"passed CVV code is not valid",
			http.StatusUnauthorized)

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

	w.Header().Set("Content-Type", CTJSON)
	w.Write(res)
}
