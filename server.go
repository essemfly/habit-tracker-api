package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/lessbutter/habit-tracker-api/config"
	"github.com/lessbutter/habit-tracker-api/graph"
	"github.com/lessbutter/habit-tracker-api/graph/generated"
)

const defaultPort = "8080"

func main() {

	conf := config.GetConfiguration()
	config.InitDB(conf)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", conf.PORT)
	log.Fatal(http.ListenAndServe(":"+conf.PORT, nil))
}
