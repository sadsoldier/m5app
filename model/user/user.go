/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */


package userModel

import (
    "log"
    "github.com/jmoiron/sqlx"
    "m5app/tools"
)

const schema = `
    DROP TABLE IF EXISTS users;
    CREATE TABLE IF NOT EXISTS users (
        id          INTEGER PRIMARY KEY,
        username    VARCHAR(255) NOT NULL UNIQUE,
        password    VARCHAR(255) NOT NULL
    );`

type Model struct {
    db *sqlx.DB
}

type User struct {
    Id          int     `db:"id"        json:"id"`
    Username    string  `db:"username"  json:"username"`
    Password    string  `db:"password"  json:"password"`
}

func (this *Model) Migrate() error {
    _, err := this.db.Exec(schema)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) List() (*[]User, error) {
    var request string
    var err error

    var users []User
    request = `SELECT id, username, '' as password FROM users
                ORDER BY username`

    err = this.db.Select(&users, request)
    if err != nil {
        log.Println(err)
        return &users, err
    }
    return &users, nil
}

func (this *Model) Create(user User) error {
    user.Password, _ = tools.CreateHash(user.Password)
    request := `INSERT INTO users(username, password) VALUES ($1, $2)`
    _, err := this.db.Exec(request, user.Username, user.Password)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) Delete(user User) error {
    request := `DELETE FROM users WHERE id = $1`
    _, err := this.db.Exec(request, user.Id)
    if err != nil {
        log.Println(err)
        return err

    }
    return nil
}

func (this *Model) Find(user User) (User, error) {
    request := `SELECT id, username, '' as password, isadmin FROM users WHERE username = $1 LIMIT 1`
    var out User
    err := this.db.Get(&out, request, user.Username)
    if err != nil {
        log.Println(err)
        return out, err
    }
    return out, nil
}

func (this *Model) Update(user User) error {
    var err error
    if len(user.Password) > 0 {
        user.Password, _ = tools.CreateHash(user.Password)
        request := `UPDATE users SET username = $1, password = $2 WHERE id = $4`
        _, err = this.db.Exec(request, user.Username, user.Password, user.Id)
    } else {
        request := `UPDATE users SET username = $1 WHERE id = $3`
        _, err = this.db.Exec(request, user.Username, user.Id)
    }
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func (this *Model) Check(user *User) error {
    username := user.Username
    password := user.Password

    request := `SELECT * FROM users WHERE username = $1 LIMIT 1`
    err := this.db.Get(user, request, username)
    if err != nil {
        log.Println(err)
        return err
    }

    err = tools.CheckHash(user.Password, password)
    if err != nil {
        log.Println(err)
        return err
    }
    return nil
}

func New(db *sqlx.DB) *Model {
    model := Model{
        db: db,
    }
    return &model
}
