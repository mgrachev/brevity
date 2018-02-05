package model

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

var allowedFields = []string{"url", "token"}

type Link struct {
	Url        string
	Token      string
	Conversion int
}

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) LinkRepository {
	return LinkRepository{db: db}
}

func (l Link) IsBlank() bool {
	return l.Url == "" && l.Token == ""
}

func (repo LinkRepository) Create(url, token string) error {
	_, err := repo.db.Exec("INSERT INTO links (url, token) VALUES ($1, $2)", url, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo LinkRepository) FindByUrl(url string) (Link, error) {
	return repo.findLink("url", url)
}

func (repo LinkRepository) FindByToken(token string) (Link, error) {
	return repo.findLink("token", token)
}

func (repo LinkRepository) IncreaseConversion(token string, conversion int) error {
	_, err := repo.db.Exec("UPDATE links SET conversion = $1 WHERE token = $2", conversion, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo LinkRepository) findLink(field, value string) (Link, error) {
	link := Link{}

	if !repo.allowedField(field) {
		return link, errors.New("Field not allowed")
	}

	query := fmt.Sprintf("SELECT url, token, conversion FROM links WHERE %s = $1", field)
	row := repo.db.QueryRow(query, value)
	err := row.Scan(&link.Url, &link.Token, &link.Conversion)
	if err == sql.ErrNoRows {
		return link, nil
	} else if err != nil {
		return link, err
	}
	return link, nil
}

func (repo LinkRepository) allowedField(field string) bool {
	for _, v := range allowedFields {
		if v == field {
			return true
		}
	}
	return false
}
