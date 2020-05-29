/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package server

import (
    "errors"
    "flag"
    "fmt"
    "html/template"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    //"bytes"
    //"encoding/base64"
    //"encoding/json"
    //"os/user"
    //"strconv"
    //"syscall"

    "github.com/gin-gonic/gin"
    "github.com/jessevdk/go-assets"

    "m5app/server/config"
    "m5app/server/daemon"
    "m5app/server/bundle"

    "m5app/server/controller"
    "m5app/server/middleware"

    "m5app/server/login"
    "m5app/server/login/config"

    "m5app/tools"
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

    router.Use(login.CheckMiddleware())
    router.Use(gin.LoggerWithFormatter(logFormatter()))
    router.Use(gin.Recovery())

    /* read templates */
    if this.config.Devel {
        /* filesystem variant */
        router.LoadHTMLGlob(filepath.Join(this.config.LibDir, "public/index.html"))
    } else {
        /* embedded variant */
        data, err := ioutil.ReadAll(this.files["/public/index.html"])
        if err != nil {
            return err
        }
        tmpl, err := template.New("index.html").Parse(string(data))
        router.SetHTMLTemplate(tmpl)
    }

    /* set route handlers */
    router.GET("/", this.Index)


    loginConfigInst := loginConfig.New()

    loginInst := login.New(loginConfigInst)
    router.POST("/login", loginInst.Login)

    controllerInst := controller.New(this.config)
    router.GET("/hello", controllerInst.Hello)
    router.POST("/hello", controllerInst.Hello)

    routerGroup := router.Group("/api/v1")
    routerGroup.Use(login.JwtAuthMiddleware(loginConfigInst))
    routerGroup.GET("/hello", controllerInst.Hello)
    routerGroup.POST("/hello", controllerInst.Hello)

    /* noroute handler */
    router.NoRoute(this.NoRoute)

    /* start run loop */
    log.Printf("start listen on :%d", this.config.Port)
    return router.Run(":" + fmt.Sprintf("%d", this.config.Port))
}

func (this *Server) Index(context *gin.Context) {
    context.HTML(http.StatusOK, "index.html", nil)
}

func (this *Server) NoRoute(context *gin.Context) {

    requestPath := context.Request.URL.Path

    if this.config.Devel {
        /* filesystem assets */
        publicDir := filepath.Join(this.config.LibDir, "public")
        filePath := filepath.Clean(filepath.Join(publicDir, requestPath))
        if !strings.HasPrefix(filePath, publicDir) {
            err := errors.New(fmt.Sprintf("wrong file patch %s\n", filePath))
            log.Println(err)
            context.HTML(http.StatusOK, "index.html", nil)
            return
        }
        /* for frontend handle: If file not found will send index.html */
        if !tools.FileExists(filePath) {
            err := errors.New(fmt.Sprintf("file path not found %s\n", filePath))
            log.Println(err)
            context.HTML(http.StatusOK, "index.html", nil)
            return
        }
        context.File(filePath)
    } else {
        /* embedded assets variant */
        file := this.files[filepath.Join("/public", requestPath)] //io.Reader
        if file == nil {
            err := errors.New(fmt.Sprintf("file path not found %s, send index", requestPath))
            log.Println(err)
            context.HTML(http.StatusOK, "index.html", nil)
            return
        }
        http.ServeContent(context.Writer, context.Request, requestPath, file.ModTime(), file)
    }
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

func New() *Server {
    return &Server{
    }
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
