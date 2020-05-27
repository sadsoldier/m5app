/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package jwt

import (
    "time"
    "net/http"
    "strings"
    "bytes"
    "errors"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/lestrrat-go/jwx/jwa"
    "github.com/lestrrat-go/jwx/jwt"

    "m5app/server/config"
    "m5app/server/controller"
)

func CheckMiddleware() gin.HandlerFunc {
    return func(context *gin.Context) {
        context.Next()
    }
}

type Controller struct{
    config  *config.Config
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

    token := jwt.New()
    token.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
    token.Set(jwt.IssuedAtKey, time.Now())
    token.Set(jwt.IssuerKey, "Issuer")
    token.Set(jwt.SubjectKey, "Subject")
    token.Set(jwt.AudienceKey, "Audience")
    token.Set("username", request.Username)

    signature, err := jwt.Sign(token, jwa.HS256, []byte("secretword"))
    if err != nil {
        controller.SendError(context, err)
        return
    }

    result := loginResult{
        Token: string(signature),
    }
    controller.SendResult(context, &result)
}

func New(config *config.Config) *Controller {
    return &Controller{
        config: config,
    }
}

func JwtAuthMiddleware(context *gin.Context) {
    authHeader := context.Request.Header.Get("Authorization")
    signature, err := parseAuthHeader(authHeader)
    if err != nil {
        response := controller.Response{
            Error: true,
            Message: fmt.Sprintf("parse auth header error: %s", err),
            Result: "",
        }
        context.JSON(http.StatusUnauthorized, response)
        context.Abort()
        return
    }

    _, err = jwt.Parse(bytes.NewReader([]byte(signature)), jwt.WithVerify(jwa.HS256, []byte("secretword")))

    if err != nil {
        response := controller.Response{
            Error: true,
            Message: fmt.Sprintf("parse auth header error: %s", err),
            Result: "",
        }
        context.JSON(http.StatusUnauthorized, response)
        context.Abort()
        return
    }
    context.Next()
}

func parseAuthHeader(header string) (string, error) {
    auth := strings.SplitN(header, " ", 2)
    authType := strings.TrimSpace(auth[0])

    if authType != "Bearer" {
        return "", errors.New("authentification type is different from bearer")
    }
    token := strings.TrimSpace(auth[1])

    if len(token) == 0 {
        return "", errors.New("token lenght is zero")
    }
    return token, nil
}
