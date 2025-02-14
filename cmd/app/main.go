package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/notblinkyet/url_shortner/internal/app"
	"github.com/notblinkyet/url_shortner/internal/config"
)

func main() {
	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger.Println(cfg)
	logger.Println("Success read config and setup logger")
	app := app.New(logger, cfg)

	go app.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.GRPCServer.Stop()
	logger.Println("Gracefull stopped")
}
