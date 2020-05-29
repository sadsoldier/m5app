/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package login

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/lestrrat-go/jwx/jwt"
    "github.com/jmoiron/sqlx"

    "m5app/server/controller"

    "m5app/server/login/tools"
    "m5app/server/login/config"

    "m5app/model/user"
)

type Controller struct{
    config  *loginConfig.Config
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
        controller.SendError(context, err)
        return
    }

    user := userModel.User{
        Username: request.Username,
        Password: request.Password,
    }

    userModelImpl := userModel.New(this.db)
    err = userModelImpl.Check(&user)
    if err != nil {
        controller.SendError(context, err)
        return
    }

    token := jwt.New()
    token.Set(jwt.ExpirationKey, time.Now().Add(time.Minute * time.Duration(this.config.Duration)))
    token.Set(jwt.IssuedAtKey, time.Now())
    token.Set(jwt.IssuerKey, this.config.Issuer)
    token.Set(jwt.SubjectKey, this.config.Subject)
    token.Set(jwt.AudienceKey, this.config.Subject)
    token.Set("username", request.Username)


    signature, err := jwt.Sign(token, loginTools.SigName2Type(this.config.SignAlg), []byte(this.config.Secret))
    if err != nil {
        controller.SendError(context, err)
        return
    }

    context.SetCookie(this.config.CookieName, string(signature),
        int(time.Duration(time.Minute).Seconds()) * this.config.Duration,
        "/", this.config.Hostname, false, true)

    result := loginResult{
        Token: string(signature),
    }
    controller.SendResult(context, &result)
}

func New(config *loginConfig.Config, db *sqlx.DB) *Controller {
    return &Controller{
        config: config,
        db:     db,
    }
}
