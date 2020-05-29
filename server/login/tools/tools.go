/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package loginTools

import (
    "strings"
    "github.com/lestrrat-go/jwx/jwa"
)

func SigName2Type(sigName string) jwa.SignatureAlgorithm {
    var algType jwa.SignatureAlgorithm
    switch algStr := strings.ToUpper(sigName); algStr {
    case "HS256":
            algType = jwa.HS256
    case "HS384":
            algType = jwa.HS384
    case "HS512":
            algType = jwa.HS512
    case "RS256":
            algType = jwa.RS256
    case "RS384":
            algType = jwa.RS384
    case "RS512":
            algType = jwa.RS512
    case "PS256":
            algType = jwa.PS256
    case "PS384":
            algType = jwa.PS384
    default:
            algType = jwa.HS512
    }
    return algType
}
