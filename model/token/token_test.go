/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package tokenModel

import (
    "testing"
    "fmt"
)

func TestAll(t *testing.T) {

    secret := "secretX"
    subject := "subjectX"
    issuer := "issuerX"

    tokenA, err := New(subject, issuer, 3600, secret)
    if err != nil {
        t.Error(err)
    }
    err = tokenA.VerifySign([]byte(secret))
    if err != nil {
        t.Error(err)
    }

    tokenData := tokenA.GetToken()
    fmt.Println(string(tokenData))

    tokenM, err := Parse(tokenData)
    if err != nil {
        t.Error(err)
    }

    if subject != tokenM.GetSub() {
        t.Error("subject should be equal")
    }
    if issuer != tokenM.GetIss() {
        t.Error("issuer should be equal")
    }

    err = tokenM.VerifySign([]byte(secret))
    if err != nil {
        t.Error(err)
    }

    if tokenA.GetSub() != tokenM.GetSub() {
        t.Error("subject should be equal")
    }
    if tokenA.GetIss() != tokenM.GetIss() {
        t.Error("issuer should be equal")
    }

}
