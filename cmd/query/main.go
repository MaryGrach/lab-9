package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "ps1"
	password = "1103"
	dbname   = "lr8"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

func (h *Handlers) GetQuery(c echo.Context) error {
	name := c.QueryParam("msg")

	if name == "" {
		return c.String(http.StatusBadRequest, "The parameter is not entered")
	}

	test, err := h.dbProvider.SelectQuery(name)
	if !test && err == nil {
		return c.String(http.StatusBadRequest, "The note has not been added to DB")
	} else if !test && err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Hello, "+name+"!")
}

func (h *Handlers) PostQuery(c echo.Context) error {
	name := c.QueryParam("msg")

	if name == "" {
		return c.String(http.StatusBadRequest, "The parameter is not entered")
	}

	test, err := h.dbProvider.SelectQuery(name)
	if test && err == nil {
		return c.String(http.StatusBadRequest, "The note has already been added to DB")
	}

	err = h.dbProvider.InsertQuery(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "Note added")
}

func (dp *DatabaseProvider) SelectQuery(msg string) (bool, error) {
	var rec string

	row := dp.db.QueryRow("SELECT name_query FROM query WHERE name_query = ($1)", msg)
	err := row.Scan(&rec)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (dp *DatabaseProvider) InsertQuery(msg string) error {
	_, err := dp.db.Exec("INSERT INTO query (name_query) VALUES ($1)", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	address := flag.String("address", "127.0.0.1:8083", "server startup address")
	flag.Parse()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dp := DatabaseProvider{db: db}
	h := Handlers{dbProvider: dp}

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/query", h.GetQuery)
	e.POST("/query", h.PostQuery)

	e.Logger.Fatal(e.Start(*address))
}
