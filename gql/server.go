package gql

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

// Server configures and starts the GQL server.
func Server(ctx client.Context) {
	if !viper.GetBool("gql-server") {
		return
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		Debug:          false,
	}).Handler)

	logFile := viper.GetString("log-file")

	port := viper.GetString("gql-port")

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{
		ctx:     ctx,
		logFile: logFile,
	}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/api"))

	if viper.GetBool("gql-playground") {
		apiBase := viper.GetString("gql-playground-api-base")

		router.Handle("/webui", playground.Handler("GraphQL playground", apiBase+"/api"))
		router.Handle("/console", playground.Handler("GraphQL playground", apiBase+"/graphql"))
	}

	router.Handle("/api", srv)
	router.Handle("/graphql", srv)

	log.Info("Connect to GraphQL playground", "url", fmt.Sprintf("http://localhost:%s", port))
	err := http.ListenAndServe(":"+port, router) //nolint: all
	if err != nil {
		panic(err)
	}
}
