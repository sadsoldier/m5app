/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package loginController

import (

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "m5app/server/controller/tools"

    "m5app/server/config"
    "m5app/model/user"
    "m5app/model/token"
)

type Controller struct{
    config  *config.Config
    db      *sqlx.DB
}


type loginRequest struct {
    Username    string  `json:"username"`
    Password    string  `json:"password"`
}

type loginResult struct {
    Token       string  `json:"token"`
}

func (this *Controller) Login(context *gin.Context) {
    var request loginRequest
    var err error

    err = context.Bind(&request)
    if err != nil {
        controllerTools.SendError(context, err)
        return
    }

    user := userModel.User{
        Username: request.Username,
        Password: request.Password,
    }

    userModelEx := userModel.New(this.db)
    err = userModelEx.Check(&user)
    if err != nil {
        controllerTools.SendError(context, err)
        return
    }

    token, err := tokenModel.New(this.config.Auth.Issuer, request.Username,
                        this.config.Auth.Duration, this.config.Auth.Secret)
    if err != nil {
        controllerTools.SendError(context, err)
        return
    }

    result := loginResult{
        Token: string(token.Token),
    }
    controllerTools.SendResult(context, &result)
}

func New(config *config.Config, db *sqlx.DB) *Controller {
    return &Controller{
        config: config,
        db:     db,
    }
}
