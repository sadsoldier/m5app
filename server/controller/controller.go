/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package controller

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Controller struct{
}

type Response struct {
    Error       bool    `json:"error"`
    Data        string  `json:"data"`
}

func (this *Controller) Hello(context *gin.Context) {
    response := Response{
        Error:      false,
        Data:       "hello",
    }
    context.JSON(http.StatusOK, response)
}

func New() *Controller {
    return &Controller{
    }
}
