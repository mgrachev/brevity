package web

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mgrachev/brevity/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	urlNotFound = "Url not found"
)

type Handler struct {
	repo model.LinkRepository
}

func NewHandler(db *sql.DB) Handler {
	repo := model.NewLinkRepository(db)
	return Handler{repo: repo}
}

func (h Handler) RedirectToOrigin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]

	link, err := h.repo.FindByToken(token)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if link.IsBlank() {
		log.Error(fmt.Sprintf("%s, using token: %s", urlNotFound, token))
		http.Error(w, urlNotFound, http.StatusNotFound)
		return
	}

	// Increase conversions
	err = h.repo.IncreaseConversion(token, link.Conversion+1)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info(fmt.Sprintf("Redirect to %s, using token: %s", link.Url, link.Token))
	http.Redirect(w, r, link.Url, http.StatusMovedPermanently)
}
