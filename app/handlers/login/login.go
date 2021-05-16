package login

import (
	"html/template"
	"log"
	"net/http"
)

type Handle struct {
	template *template.Template
}

func New(t *template.Template) *Handle {
	return &Handle{template: t}
}

func (h *Handle) Handler(w http.ResponseWriter, r *http.Request) {
	if err := h.template.Execute(w, "wss://"+r.Host+"/login"); err != nil {
		log.Println(err)
	}
}
