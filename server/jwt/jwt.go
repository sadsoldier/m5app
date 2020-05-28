/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package jwt

import (
    "time"
    "net/http"
    "strings"
    //"bytes"
    "errors"
    "fmt"
    "encoding/json"
    "encoding/base64"

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

    context.SetCookie("token", string(signature), int(time.Duration(time.Hour).Seconds()), "/", "unix7.org", false, true)

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

    var signature string
    var err error

    authHeader := context.Request.Header.Get("Authorization")


    if strings.Contains(authHeader, "Bearer") {
        signature, err = parseAuthHeader(authHeader)
        if err != nil {
            controller.AbortContext(context, http.StatusUnauthorized,
                            errors.New(fmt.Sprintf("authorization error: %s", err)))
            return
        }
    } else {
        signature, err = context.Cookie("token")
        if err != nil {
            controller.AbortContext(context, http.StatusUnauthorized,
                            errors.New(fmt.Sprintf("authorization error: %s", err)))
            return
        }
    }

    alg, _ := getJwsAlg(signature)
    _, err = jwt.Parse(strings.NewReader(signature), jwt.WithVerify(alg, []byte("secretword")))

    /* there you will be write right control */

    if err != nil {
        controller.AbortContext(context, http.StatusUnauthorized,
                            errors.New(fmt.Sprintf("authorization error: %s", err)))
        return
    }
    context.Next()
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

    var alg jwa.SignatureAlgorithm
    switch algStr := strings.ToUpper(header.Alg); algStr {
    case "HS256":
            alg = jwa.HS256
    case "HS384":
            alg = jwa.HS384
    case "HS512":
            alg = jwa.HS512
    case "RS256":
            alg = jwa.RS256
    case "RS384":
            alg = jwa.RS384
    case "RS512":
            alg = jwa.RS512
    case "PS256":
            alg = jwa.PS256
    case "PS384":
            alg = jwa.PS384
    default:
            alg = jwa.HS512
    }
    return alg, nil
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
