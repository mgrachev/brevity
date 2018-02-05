package web

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/humaniq/hmnqlog"
	"github.com/humaniq/url-shortener/config"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	db         *sql.DB
	err        error
	router     *mux.Router
	server     *httptest.Server
	logger     hmnqlog.Logger
	conversion int
)

func setup() {
	os.Setenv("APP_ENV", "test")
	//os.Setenv("TEST_DB_URL", "postgres://localhost/url-shortener_development?sslmode=disable")

	db, err = sql.Open("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	logger, _ = hmnqlog.NewZapLogger(hmnqlog.ZapOptions{
		AppName:     config.AppName,
		AppEnv:      config.AppEnv(),
		AppRevision: config.AppRevision,
	})

	router = mux.NewRouter()
	server = httptest.NewServer(router)
}

func teardown() {
	db.Close()
	server.Close()
}

func TestHandler_RedirectToOrigin(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)
	handler := NewHandler(db, logger)

	router.HandleFunc("/{token}", handler.RedirectToOrigin).Methods("GET")

	token := "aaBBcc"
	url := "http://google.com"
	// First test - url is not found
	req, _ := http.NewRequest("GET", "/"+token, nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(res.Code, http.StatusNotFound)

	// Second test - successful redirect
	db.Exec("INSERT INTO links (url, token) VALUES ($1, $2)", url, token)
	defer db.Exec("DELETE FROM links")

	req, _ = http.NewRequest("GET", "/"+token, nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	assert.Equal(res.Code, http.StatusMovedPermanently)
	assert.Equal(res.HeaderMap["Location"], []string{url})

	// Third test - increment conversion
	row := db.QueryRow("SELECT conversion FROM links")
	row.Scan(&conversion)

	assert.Equal(conversion, 1)

	req, _ = http.NewRequest("GET", "/"+token, nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	row = db.QueryRow("SELECT conversion FROM links")
	row.Scan(&conversion)

	assert.Equal(conversion, 2)
}
