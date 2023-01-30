package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"jwt/internal/repo"
	"jwt/internal/routers"
	"jwt/internal/services"
	"jwt/pkg/helpers/pg"

	"github.com/gorilla/mux"
	"github.com/pressly/goose/v3"

	_ "jwt/internal/migrations"

	_ "github.com/lib/pq"
)

func main() {
	services.LoadEnv()
	connectDb()

	fmt.Println("Application mode:", os.Getenv("APP_MODE"))
	mux := mux.NewRouter()
	routers.HandleRequest(mux)

	http.Handle("/", mux)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func connectDb() {
	cfg := &pg.Config{}
	cfg.DbName = os.Getenv("DB_USER")
	cfg.Host = "host.docker.internal"
	cfg.Port = os.Getenv("DB_PORT")
	cfg.Username = os.Getenv("DB_USER")
	cfg.Password = os.Getenv("DB_PWD")
	cfg.Timeout = 5

	poolConfig, err := pg.NewPoolConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pool config error: %v\n", err)
		os.Exit(1)
	}

	poolConfig.MaxConns = 5

	mdb, _ := sql.Open("postgres", poolConfig.ConnString())
	err = mdb.Ping()
	if err != nil {
		panic(err)
	}
	err = goose.Up(mdb, "/")
	if err != nil {
		panic(err)
	}

	connection, err := pg.NewConnection(poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect to database failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("DB connection is OK!")

	_, err = connection.Exec(context.Background(), ";")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("DB ping is OK!")

	repo.Init(connection)
}
