package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/lessbutter/habit-tracker-api/auth"
	"github.com/lessbutter/habit-tracker-api/config"
	"github.com/lessbutter/habit-tracker-api/graph"
	"github.com/lessbutter/habit-tracker-api/graph/generated"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {

	conf := config.GetConfiguration()
	config.InitDB(conf)

	router := chi.NewRouter()
	router.Use(auth.Middleware())

	router.Use(cors.AllowAll().Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", conf.PORT)
	log.Fatal(http.ListenAndServe(":"+conf.PORT, router))
}
