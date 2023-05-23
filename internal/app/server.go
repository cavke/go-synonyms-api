package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"go-synonyms-api/internal/controller"
	"go-synonyms-api/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort    = 8080
	defaultTimeout = 15 * time.Second
)

func Run() {
	var port int
	flag.IntVar(&port, "port", defaultPort, "port to start service on")
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", defaultTimeout, "the duration for which the server gracefully wait for existing connections to finish")
	flag.Parse()

	graphSynonymer := service.NewGraphSynonymer()

	synonymController := controller.SynonymController{
		Synonymer: graphSynonymer,
	}

	router := loadRoutes(synonymController)

	addr := fmt.Sprintf(":%d", port)
	srv := &http.Server{
		Handler:      router,
		Addr:         addr,
		WriteTimeout: defaultTimeout,
		ReadTimeout:  defaultTimeout,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Server running on: %v", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("srv.Shutdown: %v", err)
	}
}

func loadRoutes(userHandler controller.SynonymController) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/synonym", userHandler.GetSynonym).Methods("GET")
	router.HandleFunc("/api/v1/synonym", userHandler.CreateSynonyms).Methods("POST")
	return router
}
