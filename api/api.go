package api

import (
	"log"
	"net/http"

	"github.com/chall-goflutter-api/api/handler"
	"github.com/chall-goflutter-api/internal/interaction"
	"github.com/chall-goflutter-api/internal/kermesse"
	"github.com/chall-goflutter-api/internal/stand"
	"github.com/chall-goflutter-api/internal/ticket"
	"github.com/chall-goflutter-api/internal/tombola"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
)

type APIServer struct {
	address string
	db      *sqlx.DB
}

func NewAPIServer(address string, db *sqlx.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	userStore := user.NewStore(s.db)
	userService := user.NewService(userStore)
	userHandler := handler.NewUserHandler(userService, userStore)
	userHandler.RegisterRoutes(router)

	standStore := stand.NewStore(s.db)
	standService := stand.NewService(standStore)
	standHandler := handler.NewStandHandler(standService, userStore)
	standHandler.RegisterRoutes(router)

	kermesseStore := kermesse.NewStore(s.db)
	kermesseService := kermesse.NewService(kermesseStore, userStore)
	kermesseHandler := handler.NewKermesseHandler(kermesseService, userStore)
	kermesseHandler.RegisterRoutes(router)

	interactionStore := interaction.NewStore(s.db)
	interactionService := interaction.NewService(interactionStore, standStore, userStore, kermesseStore)
	interactionHandler := handler.NewInteractionHandler(interactionService, userStore)
	interactionHandler.RegisterRoutes(router)

	tombolaStore := tombola.NewStore(s.db)
	tombolaService := tombola.NewService(tombolaStore, kermesseStore)
	tombolaHandler := handler.NewTombolaHandler(tombolaService, userStore)
	tombolaHandler.RegisterRoutes(router)

	ticketStore := ticket.NewStore(s.db)
	ticketService := ticket.NewService(ticketStore, tombolaStore, userStore)
	ticketHandler := handler.NewTicketHandler(ticketService, userStore)
	ticketHandler.RegisterRoutes(router)

	router.HandleFunc("/webhook", handler.HandleWebhook(userService)).Methods(http.MethodPost)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r := c.Handler(router)

	log.Printf("Starting server on %s", s.address)
	return http.ListenAndServe(s.address, r)
}
