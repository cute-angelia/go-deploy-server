package controller

import (
	"fmt"
	"github.com/cute-angelia/go-utils/http/api"
	"github.com/cute-angelia/go-utils/http/validation"
	"github.com/go-chi/chi/v5"
	"go-deploy/tmpl"
	"net/http"
)

type Index struct {
}

func (rs Index) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.index)
	return r
}

func (rs Index) index(w http.ResponseWriter, r *http.Request) {
	valid := validation.Validation{}
	var u = struct {
		Page int32
	}{
		Page: api.PostInt32(r, "page"),
	}

	if err := valid.Submit(u); err != nil {
		api.Error(w, r, nil, err.Error(), -1)
		return
	}

	fmt.Fprintf(w, "%s\n", tmpl.GetIndexTpl())
}
