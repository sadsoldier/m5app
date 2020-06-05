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

    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"

    "m5app/server/config"
    "m5app/server/daemon"
    "m5app/server/bundle"

    "m5app/server/controller/login"
    "m5app/server/controller/hello"
    "m5app/server/middleware"

    "m5app/tools"
)

type Server struct {
    Config      *config.Config
    Dbx         *sqlx.DB
    files       map[string]*assets.File
}

func (this *Server) Run() error {
    var err error

    /* daemonize process */
    daemon := daemon.New(this.Config)
    err = daemon.Daemonize()
    if err != nil {
        return err
    }
    /* set signal handlers */
    daemon.SetSignalHandler()

    /* init embedded assets */
    this.files = bundle.Assets.Files


    /* set DB handler */
    dbUrl := fmt.Sprintf("%s", this.Config.DbPath)
    this.Dbx, err = sqlx.Open("sqlite3", dbUrl)
    if err != nil {
        return err
    }
    /* check DB connection */
    err = this.Dbx.Ping()
    if err != nil {
        return err
    }

    /* setup gin */
    this.setupGin()

    /* create and setup router */
    router := gin.New()

    if this.Config.Debug {
        router.Use(middleware.RequestLogMiddleware())
        router.Use(middleware.ResponseLogMiddleware())
    }

    router.Use(gin.LoggerWithFormatter(logFormatter()))
    router.Use(gin.Recovery())

    /* read templates */
    if this.Config.Devel {
        /* filesystem variant */
        router.LoadHTMLGlob(filepath.Join(this.Config.LibDir, "public/index.html"))
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


    loginEx := loginController.New(this.Config, this.Dbx)
    router.POST("/api/v1/login", loginEx.Login)

    helloControllerIm := helloController.New(this.Config, this.Dbx)


    routerGroup := router.Group("/api/v1")
    routerGroup.Use(middleware.TokenAuthMiddleware(this.Config))
    routerGroup.GET("/hello", helloControllerIm.Hello)

    /* noroute handler */
    router.NoRoute(this.NoRoute)

    /* start run loop */
    log.Printf("start listen on :%d", this.Config.Port)
    return router.Run(":" + fmt.Sprintf("%d", this.Config.Port))
}

func (this *Server) Index(context *gin.Context) {
    context.HTML(http.StatusOK, "index.html", nil)
}

func (this *Server) NoRoute(context *gin.Context) {

    requestPath := context.Request.URL.Path

    if this.Config.Devel {
        /* filesystem assets */
        publicDir := filepath.Join(this.Config.LibDir, "public")
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
    this.Config = config.New()
    this.Config.Read()
    //this.Config.Write()

    /* parse cli options */
    optForeground := flag.Bool("foreground", false, "run in foreground")
    flag.BoolVar(optForeground, "f", false, "run in foreground")

    optPort := flag.Int("port", this.Config.Port, "listen port")
    flag.IntVar(optPort, "p", this.Config.Port, "listen port")

    optDebug := flag.Bool("debug", this.Config.Debug, "debug mode")
    flag.BoolVar(optDebug, "d", false, "debug mode")

    optDevel := flag.Bool("devel", this.Config.Devel, "devel mode")
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

    this.Config.Port = *optPort
    this.Config.Debug = *optDebug
    this.Config.Devel = *optDevel
    this.Config.Foreground = *optForeground
}

func (this *Server) setupGin() error {
    gin.DisableConsoleColor()
    if this.Config.Debug{
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }
    accessLogFile, err := os.OpenFile(this.Config.AccessLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
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
