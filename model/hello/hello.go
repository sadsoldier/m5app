/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package helloModel

import (
    "log"
    "github.com/jmoiron/sqlx"
)

type Model struct {
    db *sqlx.DB
}

type Hello struct {
    Message string     `db:"message" json:"message"`
}

func (this *Model) Hello() (*[]Hello, error) {
    var request string
    var err error

    var hello []Hello
    request = `SELECT 'hello' AS message`

    err = this.db.Select(&hello, request)
    if err != nil {
        log.Println(err)
        return &hello, err
    }
    return &hello, nil
}

func New(db *sqlx.DB) *Model {
    model := Model{
        db: db,
    }
    return &model
}
