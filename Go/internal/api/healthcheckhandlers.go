package api

import (
	"net/http"

	"github.com/r3d5un/rosetta/Go/internal/rest"
)

type HealthCheckMessage struct {
	Status string `json:"status"`
}

func (api *API) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	rest.RespondWithJSON(
		w, r, http.StatusOK, HealthCheckMessage{Status: "available"}, nil)
}
