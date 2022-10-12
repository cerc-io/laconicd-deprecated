package gql

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ethereum/go-ethereum/log"
	"github.com/spf13/viper"
)

// Server configures and starts the GQL server.
func Server(ctx client.Context) {
	if !viper.GetBool("gql-server") {
		return
	}
	logFile := viper.GetString("log-file")

	port := viper.GetString("gql-port")

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{
		ctx:     ctx,
		logFile: logFile,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/api"))

	if viper.GetBool("gql-playground") {
		apiBase := viper.GetString("gql-playground-api-base")

		http.Handle("/webui", playground.Handler("GraphQL playground", apiBase+"/api"))
		http.Handle("/console", playground.Handler("GraphQL playground", apiBase+"/graphql"))
	}

	http.Handle("/api", srv)
	http.Handle("/graphql", srv)

	log.Info("Connect to GraphQL playground", "url", fmt.Sprintf("http://localhost:%s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
