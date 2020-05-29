/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package login

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/lestrrat-go/jwx/jwa"
    "github.com/lestrrat-go/jwx/jwt"

    "m5app/server/controller"

    "m5app/server/login/tools"
    "m5app/server/login/config"
)

func CheckMiddleware() gin.HandlerFunc {
    return func(context *gin.Context) {
        context.Next()
    }
}

func JwtAuthMiddleware(config *loginConfig.Config) gin.HandlerFunc {
    return func(context *gin.Context) {

        var signature string
        var err error

        authHeader := context.Request.Header.Get("Authorization")

        if strings.Contains(authHeader, config.BearerName) {
            signature, err = parseAuthHeader(authHeader, config.BearerName)
            if err != nil {
                controller.AbortContext(context, http.StatusUnauthorized,
                                errors.New(fmt.Sprintf("authorization error: %s", err)))
                return
            }
        } else {
            signature, err = context.Cookie(config.TokenName)
            if err != nil {
                controller.AbortContext(context, http.StatusUnauthorized,
                                errors.New(fmt.Sprintf("authorization error: %s", err)))
                return
            }
        }

        alg, err := getJwsAlg(signature)
        if err != nil {
            controller.AbortContext(context, http.StatusUnauthorized,
                                errors.New(fmt.Sprintf("authorization error: %s", err)))
            return
        }

        _, err = jwt.Parse(strings.NewReader(signature), jwt.WithVerify(alg, []byte(config.Secret)))

        /* there you will be write right control */

        if err != nil {
            controller.AbortContext(context, http.StatusUnauthorized,
                                errors.New(fmt.Sprintf("authorization error: %s", err)))
            return
        }
        context.Next()
    }
}

func getJwsAlg(signature string) (jwa.SignatureAlgorithm, error) {
    signatureArr := strings.SplitN(signature, ".", 2)

    headerStr, err := base64.StdEncoding.DecodeString(signatureArr[0])
    if err != nil {
        return jwa.HS256, err
    }

    type Header struct {
        Alg     string     `json:"alg"`
        Typ     string     `json:"typ"`
    }

    var header Header
    err = json.Unmarshal([]byte(headerStr), &header)
    if err != nil {
        return jwa.HS256, err
    }

    return loginTools.SigName2Type(header.Alg), nil
}

func parseAuthHeader(header string, bearerName string) (string, error) {
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
