package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"tstUser/internal/config"
	"tstUser/internal/http-server/handlers/products"
	"tstUser/internal/http-server/handlers/redirect"
	"tstUser/internal/http-server/handlers/user/check"
	"tstUser/internal/http-server/handlers/user/create"
	"tstUser/internal/http-server/handlers/user/delete"
	"tstUser/internal/http-server/middleware/logger"
	"tstUser/internal/lib/logger/handlers/slogpretty"
	"tstUser/internal/lib/logger/sl"
	"tstUser/internal/storage/sqlite"
)

const (
	envLocal = "local"
)

func main() {
	//читаем конфиг файл
	cfg := config.MustLoad()
	//инициализируем логгер в перменной окружения, работающей в данной ситуации
	log := setupLogger(cfg.Env)
	log.Info("starting tst-user", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	//инициализируем бд
	storageUser, err := sqlite.NewUserTable(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	storageProducts, err := sqlite.NewProductsTable(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	//инициализируем роутер
	router := chi.NewRouter()
	//инициализируем middleware для роутера
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/user", func(r chi.Router) {
		r.Use(middleware.BasicAuth("tst-user", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", create.New(log, storageUser))
		r.Delete("/", delete.New(log, storageUser))
		r.Get("/", check.New(log, storageUser))
	})

	router.Route("/product", func(r chi.Router) {
		r.Use(middleware.BasicAuth("tst-user", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", products.CreateProduct(log, storageProducts))
	})

	router.Get("/{mail}", redirect.New(log, storageUser))
	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

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
