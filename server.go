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

	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8080",
			"http://127.0.0.1:8080",
			"http://127.0.0.1:3000",
			"http://localhost:3000",
		},
		AllowCredentials: true,
		AllowedHeaders:   []string{"content-type"},
		Debug:            true,
	}).Handler)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	// srv.AddTransport(&transport.Websocket{
	// 	Upgrader: websocket.Upgrader{
	// 		CheckOrigin: func(r *http.Request) bool {
	// 			// Check against your desired domains here
	// 			return r.Host == "localhost:8080"
	// 		},
	// 		ReadBufferSize:  1024,
	// 		WriteBufferSize: 1024,
	// 	},
	// })

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", conf.PORT)
	log.Fatal(http.ListenAndServe(":"+conf.PORT, nil))
}
