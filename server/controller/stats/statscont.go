/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package statsController

import (

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "m5app/server/config"
    "m5app/server/controller/tools"
    "m5app/model/stats"
)

type Controller struct{
    config  *config.Config
    dbx     *sqlx.DB
}

type statsRequest struct {
    MspName     string  `json:"mspName"     binding:"required"`
    Year        string  `json:"year"        binding:"required"`
}

func (this *Controller) GetStats(context *gin.Context) {
    var request statsRequest
    var err error

    err = context.Bind(&request)
    if err != nil {
        controllerTools.SendError(context, err)
        return
    }

    model := statsModel.New(this.dbx)

    result, err := model.GetStats(request.MspName, request.Year)
    if err != nil {
        controllerTools.SendError(context, err)
        return
    }
    controllerTools.SendResult(context, &result)
}

func New(config *config.Config, dbx *sqlx.DB) *Controller {
    return &Controller{
        config: config,
        dbx:    dbx,
    }
}
