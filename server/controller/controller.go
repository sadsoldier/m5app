/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package controller

import (
    "net/http"
    "errors"
    "fmt"
    "log"

    "github.com/gin-gonic/gin"

    "m5app/server/config"
)

type Controller struct{
    config  *config.Config
}

type Response struct {
    Error       bool        `json:"error"`
    Message     string      `json:"message,omitempty"`
    Result      interface{} `json:"result,omitempty"`
}

func SendError(context *gin.Context, err error) {
    if err == nil {
        err = errors.New("undefined")
    }
    log.Printf("%s\n", err)
    response := Response{
        Error: true,
        Message: fmt.Sprintf("%s", err),
        Result: nil,
    }
    context.JSON(http.StatusOK, response)
}

func SendOk(context *gin.Context) {
    response := Response{
        Error: false,
        Message: "",
        Result: nil,
    }
    context.JSON(http.StatusOK, response)
}

func SendMessage(context *gin.Context, message string) {
    log.Printf("%s\n", message)
    response := Response{
        Error: false,
        Message: fmt.Sprintf("%s", message),
        Result: nil,
    }
    context.JSON(http.StatusOK, response)
}

func SendResult(context *gin.Context, result interface{}) {
    response := Response{
        Error: false,
        Result: result,
    }
    context.JSON(http.StatusOK, &response)
}

func (this *Controller) Hello(context *gin.Context) {
    SendMessage(context, "hello")
}

func New(config *config.Config) *Controller {
    return &Controller{
        config: config,
    }
}
