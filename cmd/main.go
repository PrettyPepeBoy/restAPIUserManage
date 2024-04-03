package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tstUser/internal/config"
	"tstUser/internal/http-server/handlers/operations"
	"tstUser/internal/http-server/handlers/products"
	"tstUser/internal/http-server/handlers/user"
	"tstUser/internal/http-server/middleware/logger"
	"tstUser/internal/lib/logger/handlers/slogpretty"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage/storages"
)

const (
	envLocal = "local"
)

var Router *chi.Mux
var log *slog.Logger

func main() {
	cfg := config.MustLoad()
	log = setupLogger(cfg.Env)
	Router = setupRouter(log, cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      Router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()
	<-ctx.Done()
	log.Info("stopping server")
	time.Sleep(time.Second * 5)
}

func setupRouter(log *slog.Logger, cfg *config.Config) *chi.Mux {
	log.Info("starting tst-user", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storageUser, err := storages.NewUserTable(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	storageProducts, err := storages.NewProductsTable(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/user", func(r chi.Router) {
		r.Use(middleware.BasicAuth("tst-user", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", user.CreateUser(log, storageUser))
		r.Delete("/", user.DeleteUser(log, storageUser))
		r.Get("/", user.FindUser(log, storageUser))
		r.Put("/", user.UpdateUser(log, storageUser))
	})

	router.Route("/product", func(r chi.Router) {
		r.Use(middleware.BasicAuth("tst-user", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", products.CreateProduct(log, storageProducts))
		r.Get("/", products.GetProduct(log, storageProducts))
		r.Put("/", products.UpdateProduct(log, storageProducts))
		r.Delete("/", products.DeleteProduct(log, storageProducts))
	})

	router.Route("/operations", func(r chi.Router) {
		r.Use(middleware.BasicAuth("tst-user", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Put("/buy/{userID}&{productID}", operations.BuyProduct(log, storageProducts, storageUser))
	})

	return router
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlersOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
