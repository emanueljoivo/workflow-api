package api

import (
	"errors"
	"net/http"
	"os"
)

type Version struct {
	Tag  string `json:"Tag"`
	Name string `json:"Name"`
}

var (
	EncodeResErr = errors.New("error while trying encode response")
)

func (a *HttpApi) GetVersion(w http.ResponseWriter, r *http.Request) {
	Write(w, http.StatusOK, Version{Tag: os.Getenv("VERSION_TAG"), Name: os.Getenv("VERSION_NAME")})
}
