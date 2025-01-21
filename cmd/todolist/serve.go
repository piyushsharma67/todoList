package main

import (
	"net/http"
	"time"

	sqlitedb "go.altair.com/todolist/pkg/db"
	"go.altair.com/todolist/pkg/todolist"
	"go.altair.com/todolist/pkg/todolist/store"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the Todolist server",
	Long:  ``,
	RunE:  doServe,
}

var (
	bindAddress string
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&bindAddress, "bind", "b", "0.0.0.0:8080", "set the bind address for the server")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Allow certain methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Allow certain headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Allow credentials (if needed)
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// If the request method is OPTIONS, return immediately
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next middleware/handler
		next.ServeHTTP(w, r)
	})
}

func newRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(chimw.Recoverer)
	router.Use(chimw.Timeout(60 * time.Second))
	router.Use(corsMiddleware)
	return router
}

func doServe(cmd *cobra.Command, args []string) error {
    log.Info().Msg(description + " starting")

    tododb, err := sqlitedb.CreateDb()
    if err != nil {
        log.Error().Err(err).Msg("Failed to create SQLite database")
        return err
    }

    todostore := store.NewSqlStore(tododb)
    todoService := todolist.NewItemsService(todostore)

    handler := &todolist.ItemsHandlers{
        ItemsService: todoService,
    }

    router := newRouter()
    handler.ConfigureRoutes(router)

    log.Info().Str("bindAddress", bindAddress).Msg("Listening for HTTP requests")
    return http.ListenAndServe(bindAddress, router)
}

