package server

import (
	"net/http"

	"github.com/aakash-tyagi/kart-challenge/db"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

type Server struct {
	DBClient *db.Db
	Logger   *logger.Logger
}

func New(
	dbClient *db.Db,
	logger *logger.Logger,
) *Server {
	return &Server{
		DBClient: dbClient,
		Logger:   logger,
	}
}

func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/product", s.ListProducts).Methods("GET")
	// Add more routes here as needed
}

func (s *Server) Start() {
	r := mux.NewRouter()

	// Register routes
	s.RegisterRoutes(r)

	http.Handle("/", r)
	s.Logger.Info("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		s.Logger.Fatal(err)
	}

	s.Logger.Info("Server stopped")
}
