package app

import (
	"context"
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrparano1d/getregd/ent"
	"github.com/mrparano1d/getregd/pkg/adapters"
	"github.com/mrparano1d/getregd/pkg/app/auth"
	"github.com/mrparano1d/getregd/pkg/app/handler"
	"github.com/mrparano1d/getregd/pkg/core"
	"github.com/redis/go-redis/v9"

	"log"
	"net/http"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func EntClient() *ent.Client {

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	entClient, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, password))
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	if err := entClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return entClient
}

func ServeApp() error {

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	entClient := EntClient()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	authAdapter := adapters.NewAuthAdapter(entClient)
	packageAdapter := adapters.NewPackageAdapter()
	storeAdapter := adapters.NewFSStorageAdapter("./storage")
	sessionAdapter := adapters.NewSessionAdapter(redisClient)
	roleAdapter := adapters.NewRoleAdapter(entClient)
	userAdapter := adapters.NewUserAdapter(entClient)

	app := core.NewCoreApp(sessionAdapter, authAdapter, packageAdapter, storeAdapter, userAdapter, roleAdapter)
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	handler.AuthHandler(r, app)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(app))
		handler.PackageHandler(r, app)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})

	log.Println("server started at http://localhost:3000")
	return http.ListenAndServe(":3000", r)
}
