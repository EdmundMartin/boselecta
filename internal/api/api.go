package api

import (
	"encoding/json"
	"errors"
	"github.com/EdmundMartin/boselecta/internal/flag"
	"github.com/EdmundMartin/boselecta/internal/storage"
	"io/ioutil"
	"net/http"
)

type ManagementAPI struct {
	storage storage.FlagStorage
}

func NewManagementAPI(engine storage.FlagStorage) *ManagementAPI {
	return &ManagementAPI{storage: engine}
}

type errorResp struct {
	Msg string
	Route string
}

type successResp struct {
	Msg string
}

func errorWithJson(w http.ResponseWriter, r *http.Request, err error, code int) {
	msg := err.Error()
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResp{Msg: msg, Route: r.RequestURI})
}

func (m *ManagementAPI) createFlag(w http.ResponseWriter, r *http.Request) {
	var newFlag flag.FeatureFlag

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		errorWithJson(w, r, err, http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &newFlag)
	if err != nil {
		errorWithJson(w, r, err, http.StatusBadRequest)
		return
	}

	err = m.storage.Create(newFlag.Namespace, &newFlag)
	if err != nil {
		errorWithJson(w, r, errors.New("unable to save flag in database"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResp{Msg: "sucessfully created"})
}