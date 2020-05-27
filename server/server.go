/*
 * Copyright 2020 Oleg Borodin  <borodin@unix7.org>
 */

package server

import (
    "bytes"
    //"encoding/base64"
    "encoding/json"
    //"errors"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    //"net/http"
    //"html/template"
    "os"
    //"os/user"
    "path/filepath"
    //"strconv"
    "strings"
    //"syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/jessevdk/go-assets"

    "m5app/server/config"
    "m5app/server/daemon"
    "m5app/server/bundle"
    "m5app/server/controller"

)

type Server struct {
    config      *config.Config
    files       map[string]*assets.File
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
    if this.config.Debug{
        router.Use(RequestLogMiddleware())
        router.Use(ResponseLogMiddleware())
    }
    router.Use(gin.LoggerWithFormatter(logFormatter()))
    router.Use(gin.Recovery())


    controller := controller.New()
    router.GET("/hello", controller.Hello)
    router.POST("/hello", controller.Hello)

    /* start run loop */
    return router.Run(":" + fmt.Sprintf("%d", this.config.Port))
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

func RequestLogMiddleware() gin.HandlerFunc {
    return func(context *gin.Context) {

        var requestBody []byte
        if context.Request.Body != nil {
            requestBody, _ = ioutil.ReadAll(context.Request.Body)
        }

        contentType := context.GetHeader("Content-Type")
        contentType = strings.ToLower(contentType)

        buffer := bytes.NewBuffer(nil)
        json.Indent(buffer, requestBody, "", "    ")

        if strings.Contains(contentType, "application/json") {
            log.Print("request:\n", buffer.String())
        }

        context.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
        context.Next()
    }
}

func ResponseLogMiddleware() gin.HandlerFunc {

    return func(context *gin.Context) {
        contentType := context.GetHeader("Content-Type")
        contentType = strings.ToLower(contentType)

        writer := &LogWriter{
            body: bytes.NewBuffer(nil),
            ResponseWriter: context.Writer,
        }
        context.Writer = writer

        context.Next()

        buffer := bytes.NewBuffer(nil)
        json.Indent(buffer, writer.body.Bytes(), "", "    ")

        if strings.Contains(contentType, "application/json") {
            log.Print("response:\n", buffer.String())
        }
    }
}

type LogWriter struct {
    gin.ResponseWriter
    body *bytes.Buffer
}

func (this LogWriter) Write(data []byte) (int, error) {
    this.body.Write(data)
    return this.ResponseWriter.Write(data)
}

func (this LogWriter) WriteString(data string) (int, error) {
    this.body.WriteString(data)
    return this.ResponseWriter.WriteString(data)
}


func New() *Server {
    return &Server{
    }
}
