/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package server

import (
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "time"

    //"bytes"
    //"encoding/base64"
    //"encoding/json"
    //"errors"
    //"html/template"
    //"io/ioutil"
    //"net/http"
    //"os/user"
    //"strconv"
    //"strings"
    //"syscall"

    "github.com/gin-gonic/gin"
    "github.com/jessevdk/go-assets"
    "github.com/appleboy/gin-jwt/v2"

    "m5app/server/config"
    "m5app/server/daemon"
    "m5app/server/bundle"

    "m5app/server/controller"
    "m5app/server/middleware"

)

type Server struct {
    config      *config.Config
    files       map[string]*assets.File
}

func (this *Server) Run() error {
    var err error

    /* daemonize process */
    daemon := daemon.New(this.config)
    err = daemon.Daemonize()
    if err != nil {
        return err
    }
    /* set signal handlers */
    daemon.SetSignalHandler()

    /* init embedded assets */
    this.files = bundle.Assets.Files

    /* setup gin */
    this.setupGin()

    /* create and setup router */
    router := gin.New()

    if this.config.Debug {
        router.Use(middleware.RequestLogMiddleware())
        router.Use(middleware.ResponseLogMiddleware())
    }

    router.Use(gin.LoggerWithFormatter(logFormatter()))
    router.Use(gin.Recovery())

    router.POST("/login", authMiddleware.LoginHandler)

    controller := controller.New()
    router.GET("/hello", controller.Hello)
    router.POST("/hello", controller.Hello)

    /* start run loop */
    log.Printf("start listen on :%d", this.config.Port)
    return router.Run(":" + fmt.Sprintf("%d", this.config.Port))
}

func (this *Server) Configure() {

    /* read configuration file */
    this.config = config.New()
    this.config.Read()
    //this.config.Write()

    /* parse cli options */
    optForeground := flag.Bool("foreground", false, "run in foreground")
    flag.BoolVar(optForeground, "f", false, "run in foreground")

    optPort := flag.Int("port", this.config.Port, "listen port")
    flag.IntVar(optPort, "p", this.config.Port, "listen port")

    optDebug := flag.Bool("debug", this.config.Debug, "debug mode")
    flag.BoolVar(optDebug, "d", false, "debug mode")

    optDevel := flag.Bool("devel", this.config.Devel, "devel mode")
    flag.BoolVar(optDebug, "e", false, "devel mode")

    optWrite := flag.Bool("write", false, "write config")
    flag.BoolVar(optWrite, "w", false, "write config")

    exeName := filepath.Base(os.Args[0])

    flag.Usage = func() {
        fmt.Println("")
        fmt.Printf("usage: %s command [option]\n", exeName)
        fmt.Println("")
        flag.PrintDefaults()
        fmt.Println("")
    }
    flag.Parse()

    this.config.Port = *optPort
    this.config.Debug = *optDebug
    this.config.Devel = *optDevel
    this.config.Foreground = *optForeground
}

func (this *Server) setupGin() error {
    gin.DisableConsoleColor()
    if this.config.Debug{
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }
    accessLogFile, err := os.OpenFile(this.config.AccessLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
    if err != nil {
      return err
    }
    gin.DefaultWriter = io.MultiWriter(accessLogFile, os.Stdout)
    //gin.DefaultWriter = ioutil.Discard
    return nil
}

func logFormatter() func(param gin.LogFormatterParams) string {
    return func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s %s %s %s %s %d %d %s\n",
            param.TimeStamp.Format(time.RFC3339),
            param.ClientIP,
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.BodySize,
            param.Latency,
        )
    }
}
