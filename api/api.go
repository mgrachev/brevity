package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mgrachev/brevity/config"
	"github.com/mgrachev/brevity/model"
	"github.com/mgrachev/brevity/token"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	cannotParseFormData   = "Cannot parse form data"
	urlNotFound           = "Url not found"
	cannotCheckUrl        = "Cannot check Url"
	cannotCreateShortlink = "Cannot create shortlink"
)

var identifier string

type Handler struct {
	repo model.LinkRepository
}

func NewHandler(db *sql.DB) Handler {
	repo := model.NewLinkRepository(db)
	return Handler{repo: repo}
}

func (h Handler) FindOrCreateShortLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		response := ErrorResponse{Error: cannotParseFormData}
		handleError(h, w, http.StatusUnprocessableEntity, response)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		response := ErrorResponse{Error: urlNotFound}
		handleError(h, w, http.StatusUnprocessableEntity, response)
		return
	}

	// Check if there is already such a URL in the database
	link, err := h.repo.FindByUrl(url)
	if err != nil {
		log.Error(err.Error())
		response := ErrorResponse{Error: cannotCheckUrl}
		handleError(h, w, http.StatusUnprocessableEntity, response)
		return
	}

	if link.IsBlank() {
		// If not, create a new record in the database
		identifier, err = h.createShortLink(url)
		if err != nil {
			log.Error(err.Error())
			response := ErrorResponse{Error: cannotCreateShortlink}
			handleError(h, w, http.StatusUnprocessableEntity, response)
			return
		}
	} else {
		identifier = link.Token
	}

	shortLink := fmt.Sprintf("%s/%s", config.AppDomain(), identifier)
	data := make(map[string]string)
	data["shortlink"] = shortLink

	log.Info(fmt.Sprintf("Successful returns a short link, short: %s, original: %s", shortLink, url))

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h Handler) createShortLink(url string) (string, error) {
	token, err := h.generateUniqueToken()
	if err != nil {
		return "", err
	}

	err = h.repo.Create(url, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h Handler) generateUniqueToken() (string, error) {
	identifier := token.GenerateToken(config.AppTokenLength())

	link, err := h.repo.FindByToken(identifier)
	if err != nil {
		return "", err
	}

	if link.IsBlank() {
		return identifier, nil
	}
	// Token is already exists. Try again...
	return h.generateUniqueToken()
}

func handleError(h Handler, w http.ResponseWriter, status int, response interface{}) {
	w.WriteHeader(status)

	errEncode := json.NewEncoder(w).Encode(response)
	if errEncode != nil {
		log.Error(errEncode.Error())
	}
}
