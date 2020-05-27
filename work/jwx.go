/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "time"

  "github.com/lestrrat-go/jwx/jwa"
  "github.com/lestrrat-go/jwx/jwt"
)

func main() {

    token1 := jwt.New()
    token1.Set(jwt.IssuerKey, "Issuer")
    token1.Set(jwt.SubjectKey, "SubjectKey")
    token1.Set(jwt.AudienceKey, "AudienceKey")
    token1.Set(jwt.IssuedAtKey, time.Now())
    token1.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
    token1.Set("privateClaimKey", "Hello, World!")

    json1, err := json.MarshalIndent(token1, "", "    ")
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }

    fmt.Printf("%s\n", json1)

    sign, err := jwt.Sign(token1, jwa.HS256, []byte("secretword"))
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }

    fmt.Println(string(sign))

    token2, err := jwt.ParseBytes(sign)
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }
    json2, err := json.MarshalIndent(token2, "", "    ")
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }
    fmt.Printf("%s\n", json2)

    fmt.Printf("aud: %s\n", token2.Audience())
    fmt.Printf("sub: %s\n", token2.Subject())
    fmt.Printf("iat: %s\n", token2.IssuedAt().Format(time.RFC3339))
    fmt.Printf("exp: %s\n", token2.Expiration().Format(time.RFC3339))

    _, err = jwt.Parse(bytes.NewReader(sign), jwt.WithVerify(jwa.HS256, []byte("secretword1")))
    if err != nil {
        fmt.Printf("error: %s\n", err)
        return
    }
}
