package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"tstUser/data-logic/functional"
	"tstUser/data-logic/read"
	"tstUser/data-logic/write"
	"tstUser/internal/config"
	"tstUser/internal/http-server/handlers/redirect"
	"tstUser/internal/http-server/handlers/user/check"
	"tstUser/internal/http-server/handlers/user/create"
	"tstUser/internal/http-server/handlers/user/delete"
	"tstUser/internal/http-server/middleware/logger"
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
	storage, err := sqlite.NewUserTable(cfg.StoragePath)
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
		r.Post("/", create.New(log, storage))
		r.Delete("/", delete.New(log, storage))
		r.Get("/", check.New(log, storage))
	})

	router.Get("/{mail}", redirect.New(log, storage))
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

	go func() {
		for {
			whatToDo()
		}
	}()
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	<-sigC
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func whatToDo() {
	fmt.Println("Please input what you wanna do")
	fmt.Printf("Type |%d| to find user, |%d| to create new user, |%d| to sort users by name, |%d| to sort users by data,|%d| to delete user,|%d| to buy stuff\n", 1, 2, 3, 4, 5, 6)
	fmt.Printf("Type |%d| to send cash\n", 7)
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	i := read.ReadButton()
	switch i {
	case 1:
		usr, ok := functional.CheckUsr()
		if !ok {
			fmt.Println("there is no such user")
		} else {
			fmt.Println(usr)
		}
	case 2:
		b, check := functional.CreateUser()
		if check {
			write.WriteInFile(b)
		}
	case 3:
		functional.SortByName()
	case 4:
		functional.SortByData()
	case 5:
		functional.DeleteUser()
	case 6:
		functional.BuyThing()
	case 7:
		functional.SendCash(677985, 351057, 100)
	default:
		<-sigC
		fmt.Println("exit")
	}
}
