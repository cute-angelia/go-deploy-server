package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-deploy/config"
	"net/http"
)

type List struct {
}

func (rs List) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", rs.list)
	return r
}

func (rs List) list(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(config.C.Apps)
	fmt.Fprintf(w, "%s\n", str)
}
