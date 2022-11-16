package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/alexm24/golang/internal/config"
	"github.com/alexm24/golang/internal/handler"
	"github.com/alexm24/golang/internal/models"
	"github.com/alexm24/golang/internal/server"
	"github.com/alexm24/golang/internal/service"
	"github.com/alexm24/golang/internal/transport"
	"github.com/alexm24/golang/internal/transport/postgres"
	"github.com/alexm24/golang/internal/transport/redis"
)

func App(configPath string) {
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		log.Panicf("error read config: %s", err.Error())
	}
	log.Println("config read successfully")

	db, err := postgres.NewPostgresDB(cfg.DBConfig)
	if err != nil {
		log.Panicf("failed to initialize postgres db: %s", err.Error())
	}

	rp, err := redis.NewRedisPool(cfg.RedisConfig)
	if err != nil {
		log.Panicf("failed to initialize redis db: %s", err.Error())
	}

	transports := transport.NewTransport(db, rp, cfg.CentrifugoConfig)
	services := service.NewService(transports)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)

	go func() {
		cfgHTP := models.HTTPServerConfig{
			Port: cfg.HTTPServerConfig.Port,
			Path: cfg.HTTPServerConfig.Path,
		}
		log.Printf("start http server on port %s", cfgHTP.Port)
		err = srv.Run(cfgHTP.Port, handlers.InitRoutes(cfgHTP.Path))
		if err != http.ErrServerClosed {
			log.Panicf("error occurred while running http server: %s", err.Error())
		}
	}()

	signalLisner := make(chan os.Signal, 1)
	signal.Notify(signalLisner,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	stop := <-signalLisner
	log.Printf("Shutting Down app: %s", stop)

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occurred on server shutting down: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		log.Printf("error on db connection close: %s", err.Error())
	}
}
