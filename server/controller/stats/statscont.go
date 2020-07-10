/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package statsController

import (

    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"

    "m5app/server/config"
    "m5app/server/controller/tools"
    //"m5app/model/stats"
)

type Controller struct{
    config  *config.Config
    dbx     *sqlx.DB
}

type statsRequest struct {
    InsuranceType   []string    `json:"insuranceType"`
    ClaimStatus     []string    `json:"claimStatus"`
    ClaimSubject    []string    `json:"claimSubject"`
    PolicyHolders   []string    `json:"policyHolders"`
    Insurers        []string    `json:"insurers"`
    Currency        string      `json:"currency"`
    AggregateTypes  []string    `json:"aggregateTypes"`
    Vip             bool        `json:"vip"`
}

func (this *Controller) GetStats(context *gin.Context) {
    var request statsRequest
    var err error

    err = context.Bind(&request)
    if err != nil {
        controllerTools.SendError(context, err)
        return
    }

    //model := statsModel.New(this.dbx)
    //result, err := model.GetStats()
    //if err != nil {
        //controllerTools.SendError(context, err)
        //return
    //}
    //controllerTools.SendResult(context, &result)

    controllerTools.SendResult(context, &request)
}

func New(config *config.Config, dbx *sqlx.DB) *Controller {
    return &Controller{
        config: config,
        dbx:    dbx,
    }
}
