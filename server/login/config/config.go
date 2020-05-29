/*
 * Copyright 2019 Oleg Borodin  <borodin@unix7.org>
 */

package loginConfig

import (
    "os"
)

type Config struct {
    Hostname            string  `yaml:"-"`
    CookieName          string  `yaml:"-"`
    SignAlg             string  `yaml:"-"`
    TokenName           string  `yaml:"-"`
    BearerName          string  `yaml:"-"`

    Secret              string  `yaml:"secret"`
    Issuer              string  `yaml:"issuer"`
    Subject             string  `yaml:"subject"`
    Duration            int     `yaml:"duration"`
}

func New() *Config {
    hostname, _ := os.Hostname()
    return &Config{
        Hostname:       hostname,
        CookieName:     "token",
        SignAlg:        "HS512",
        TokenName:      "token",
        BearerName:     "Bearer",

        Secret:         "secret",
        Issuer:         "m5app",
        Subject:        "m5app",
        Duration:       3600 * 1,
    }
}
