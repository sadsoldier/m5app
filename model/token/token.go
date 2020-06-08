/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */


package tokenModel

import (
    "strings"
    "encoding/base64"
    "encoding/json"
    "bytes"
    "time"
    "errors"

    "github.com/lestrrat-go/jwx/jwa"
    "github.com/lestrrat-go/jwx/jwt"
)

type Header struct {
    Alg     string  `json:"alg"`
    Typ     string  `json:"typ"`
}

type Payload struct {
    Exp     int64   `json:"exp"`
    Iss     string  `json:"iss"`
    Sub     string  `json:"sub"`
}

type Model struct {
    Token       []byte
    Header      *Header
    Payload     *Payload
}

func (this *Model) GetSub() string {
    return this.Payload.Sub
}

func (this *Model) GetIss() string {
    return this.Payload.Iss
}

func (this *Model) GetToken() []byte {
    return this.Token
}


func (this *Model) VerifySign(secret []byte) error {
    alg := SigName2Type(this.Header.Alg)
    _, err := jwt.Parse(bytes.NewReader(this.Token), jwt.WithVerify(alg, secret))
    if err != nil {
        return err
    }
    return nil
}

func (this *Model) VerifyTime() error {
    now := time.Now().Unix()
    if this.Payload.Exp == 0 {
        return errors.New("token usage time is zero")
    }
    if now > this.Payload.Exp {
        return errors.New("token usage time expired")
    }
    return nil
}

func Parse(token []byte) (*Model, error) {
    var err error

    tokenArr := strings.SplitN(string(token), ".", 3)
    header, err := ParseHeader(tokenArr[0])
    if err != nil {
        return nil, err
    }

    payload, err := ParsePayload(tokenArr[1])
    if err != nil {
        return nil, err
    }

    model := &Model{
        Header:     header,
        Payload:    payload,
        Token:      token,
    }
    return model, nil
}

func New(subject, issuer string, duration int, secret string) (*Model, error) {
    var err error

    header := Header{
        Alg:    "HS512",
        Typ:    "JWT",
    }

    jwtToken := jwt.New()

    now := time.Now()
    expired := now.Add(time.Second * time.Duration(duration))

    jwtToken.Set(jwt.ExpirationKey, expired)
    jwtToken.Set(jwt.IssuedAtKey, now)
    jwtToken.Set(jwt.IssuerKey, issuer)
    jwtToken.Set(jwt.SubjectKey, subject)

    token, err := jwt.Sign(jwtToken, SigName2Type(header.Alg), []byte(secret))
    if err != nil {
        return nil, err
    }

    payload := Payload{
        Sub:    subject,
        Iss:    issuer,
        Exp:    expired.Unix(),
    }

    model := &Model{
        Token:      token,
        Header:     &header,
        Payload:    &payload,
    }
    return model, nil
}


func ParsePayload(dataB64 string) (*Payload, error) {
    var data Payload
    var err error
    dataJson, err := base64.RawStdEncoding.DecodeString(dataB64)
    if err != nil {
        return &data, err
    }
    err = json.Unmarshal([]byte(dataJson), &data)
    if err != nil {
        return &data, err
    }
    return &data, err
}

func ParseHeader(dataB64 string) (*Header, error) {
    var data Header
    var err error
    dataJson, err := base64.RawStdEncoding.DecodeString(dataB64)
    if err != nil {
        return &data, err
    }
    err = json.Unmarshal([]byte(dataJson), &data)
    if err != nil {
        return &data, err
    }
    return &data, err
}

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
