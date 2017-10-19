package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/istyle-inc/multimissile"
	"github.com/istyle-inc/multimissile/config"
	"github.com/istyle-inc/multimissile/server"
	"github.com/istyle-inc/multimissile/wlog"
)

func main() {
	versionPrinted := flag.Bool("v", false, "print multimissile version")
	port := flag.String("p", "", "listening port number or socket path")
	configPath := flag.String("c", "", "configuration file path")
	flag.Parse()

	if *versionPrinted {
		msl.PrintVersion()
		return
	}

	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// overwrite if port is specified by flags
	if *port != "" {
		conf.Port = *port
	}

	// set global configuration
	msl.Config = conf
	msl.AL = wlog.AccessLogger(conf.LogLevel)
	msl.EL = wlog.ErrorLogger(conf.LogLevel)

	// Setup server
	mux := http.NewServeMux()
	server.RegisterHandlers(mux)
	server.SetupClient(&msl.Config)

	srv := &http.Server{
		Handler: mux,
	}

	go func() {
		msl.EL.Out(wlog.Debug, "Start running server")
		if err := server.Run(srv, &msl.Config); err != nil {
			msl.EL.Out(wlog.Error, "Failed to run server: %s", err)
		}
	}()

	// Watch SIGTERM signal and then gracefully shutdown
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	msl.EL.Out(wlog.Debug, "Start to shutdown server")
	timeout := time.Duration(conf.ShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		msl.EL.Out(wlog.Error, "Failed to shutdown server: %s", err)
		return
	}

	msl.EL.Out(wlog.Debug, "Successfully shutdown server")
}
