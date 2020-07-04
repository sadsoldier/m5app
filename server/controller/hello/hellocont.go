/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package helloController

import (
    //"net/http"
    //"errors"
    //"fmt"
    //"log"

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "m5app/server/config"
    "m5app/server/controller/tools"
    "m5app/model/hello"

)

type Controller struct{
    config  *config.Config
    db      *sqlx.DB
}

func (this *Controller) Hello(context *gin.Context) {
    model := helloModel.New(this.db)
    result, err := model.Hello()
    if err != nil {
        controllerTools.SendError(context, err)
    }
    controllerTools.SendResult(context, &result)
    //controller.SendMessage(context, "hello")
}

func New(config *config.Config, db *sqlx.DB) *Controller {
    return &Controller{
        config: config,
        db: db,
    }
}
