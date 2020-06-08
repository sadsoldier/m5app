/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package middleware

import (
    "errors"
    "fmt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"

    "m5app/model/token"
    "m5app/server/config"
    "m5app/server/controller/tools"
)


func ParseAuthHeader(header string, bearerName string) (string, error) {
    auth := strings.SplitN(header, " ", 2)
    authType := strings.TrimSpace(auth[0])

    if authType != bearerName {
        return "", errors.New(fmt.Sprintf("authentification type is different from %s", bearerName))
    }
    token := strings.TrimSpace(auth[1])

    if len(token) == 0 {
        return "", errors.New("token lenght is zero")
    }
    return token, nil
}


func TokenAuthMiddleware(config *config.Config) gin.HandlerFunc {
    return func(context *gin.Context) {

        var signature string
        var err error
        /* get authorization header */
        authHeader := context.Request.Header.Get("Authorization")
        if len(authHeader) == 0 {
            controllerTools.AbortContext(context, http.StatusUnauthorized,
                    errors.New(fmt.Sprintf("not found authorization header")))
            return
        }

        /* check bearer into header */
        if !strings.Contains(authHeader, config.Auth.BearerName) {
            controllerTools.AbortContext(context, http.StatusUnauthorized,
                    errors.New(fmt.Sprintf("authorization header not contain bearer")))
            return
        }

        /* extract token string */
        signature, err = ParseAuthHeader(authHeader, config.Auth.BearerName)
        if err != nil {
            controllerTools.AbortContext(context, http.StatusUnauthorized,
                            errors.New(fmt.Sprintf("authorization header parsing error: %s", err)))
            return
        }

        /* create token object */
        token, err := tokenModel.Parse([]byte(signature))
        if err != nil {
            controllerTools.AbortContext(context, http.StatusUnauthorized,
                    errors.New(fmt.Sprintf("token parse error: %s", err)))
            return
        }

        /* check token signature */
        err = token.VerifySign([]byte(config.Auth.Secret))
        if err != nil {
            controllerTools.AbortContext(context, http.StatusUnauthorized,
                    errors.New(fmt.Sprintf("mismatch signatures: %s", err)))
            return
        }

        /* usage time checking */
        err = token.VerifyTime()
        if err != nil {
            controllerTools.AbortContext(context, http.StatusUnauthorized,
                errors.New(fmt.Sprintf("mismatch time of usage: %s", err)))
            return
        }

        context.Next()
    }
}
