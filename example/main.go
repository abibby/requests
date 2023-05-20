package main

import (
	"context"
	"net/http"

	"github.com/abibby/bob"
	"github.com/abibby/bob/dialects"
	"github.com/abibby/bob/dialects/sqlite"
	"github.com/abibby/bob/migrate"
	"github.com/abibby/bob/models"
	"github.com/abibby/requests"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type Foo struct {
	models.BaseModel
	ID   int    `db:"id,autoincrement"`
	Name string `db:"name"`
}

type ListRequest struct {
	Request *http.Request
}

var list = requests.Handler(func(r *ListRequest) ([]*Foo, error) {
	tx := requests.UseTx(r.Request)

	foos, err := bob.From[*Foo]().Get(tx)
	if err != nil {
		return nil, err
	}
	return foos, nil
})

type AddRequest struct {
	Request *http.Request
	Name    string `query:"name"`
}

var add = requests.Handler(func(r *AddRequest) (*Foo, error) {
	tx := requests.UseTx(r.Request)
	foo := &Foo{Name: r.Name}
	err := bob.Save(tx, foo)
	if err != nil {
		return nil, err
	}
	return foo, nil
})

type GetRequest struct {
	Foo *Foo `di:"foo"`
}

var get = requests.Handler(func(r *GetRequest) (*Foo, error) {
	return r.Foo, nil
})

func main() {
	dialects.DefaultDialect = &sqlite.SQLite{}
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()

	c, err := migrate.CreateFromModel(&Foo{})
	if err != nil {
		panic(err)
	}

	err = c.Run(context.Background(), db)
	if err != nil {
		panic(err)
	}

	r.Use(requests.WithDB(db))

	r.HandleFunc("/foo", list)
	r.HandleFunc("/foo/create", add)
	r.HandleFunc("/foo/{foo}", get)

	http.ListenAndServe(":8087", r)
}
