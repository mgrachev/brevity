package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/mgrachev/brevity/config"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

var (
	db     *sql.DB
	err    error
	router *mux.Router
	server *httptest.Server
)

func setup() {
	os.Setenv("APP_ENV", "test")
	//os.Setenv("TEST_DB_URL", "postgres://localhost/url-shortener_development?sslmode=disable")

	db, err = sql.Open("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	router = mux.NewRouter()
	server = httptest.NewServer(router)
}

func teardown() {
	db.Close()
	server.Close()
}

func TestHandler_FindOrCreateShortLink(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)
	handler := NewHandler(db)

	route := "/api/v1/shortlink"
	router.HandleFunc(route, handler.FindOrCreateShortLink).Methods("POST")

	// First test - request without url
	req, _ := http.NewRequest("POST", route, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	data, err := retrieveJSON(res.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(res.Code, http.StatusUnprocessableEntity)
	assert.Equal(data["error"], cannotParseFormData)

	// Second test - url is empty
	values := url.Values{}
	values.Add("url", "")
	valStr := values.Encode()

	req, _ = http.NewRequest("POST", route, strings.NewReader(valStr))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	data, err = retrieveJSON(res.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(res.Code, http.StatusUnprocessableEntity)
	assert.Equal(data["error"], urlNotFound)

	// Third test - create a new url
	values = url.Values{}
	values.Add("url", "http://google.com")
	valStr = values.Encode()

	req, _ = http.NewRequest("POST", route, strings.NewReader(valStr))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	data, err = retrieveJSON(res.Body)
	if err != nil {
		t.Error(err)
	}

	row := db.QueryRow("SELECT token FROM links")
	row.Scan(&identifier)
	defer db.Exec("DELETE FROM links")

	shortLink := fmt.Sprintf("%s/%s", config.AppDomain(), identifier)

	assert.Equal(res.Code, http.StatusOK)
	assert.Equal(data["shortlink"], shortLink)

	// Fourth - return already created url
	req, _ = http.NewRequest("POST", route, strings.NewReader(valStr))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	data, err = retrieveJSON(res.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(res.Code, http.StatusOK)
	assert.Equal(data["shortlink"], shortLink)
}

func retrieveJSON(respBody *bytes.Buffer) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err = json.Unmarshal(respBody.Bytes(), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
