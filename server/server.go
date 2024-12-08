package server

import (
	"net/http"

	"github.com/aakash-tyagi/kart-challenge/config"
	"github.com/aakash-tyagi/kart-challenge/db"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

type Server struct {
	DBClient *db.Db
	Logger   *logger.Logger
	Config   *config.Config
}

func New(
	dbClient *db.Db,
	logger *logger.Logger,
	config *config.Config,
) *Server {
	return &Server{
		DBClient: dbClient,
		Logger:   logger,
		Config:   config,
	}
}

func (s *Server) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/product", s.ListProducts).Methods("GET")
	r.HandleFunc("/api/v1/product", s.AddProduct).Methods("POST")
	r.HandleFunc("/api/v1/product/{productId}", s.GetProductById).Methods("GET")

	r.HandleFunc("/api/v1/order", s.CreateOrder).Methods("POST")
}

func (s *Server) Start() {
	r := mux.NewRouter()
	r.Use(AuthMiddleware)

	// Register routes
	s.RegisterRoutes(r)

	http.Handle("/", r)
	s.Logger.Info("Starting server on port: ", s.Config.ServerPort)
	if err := http.ListenAndServe(":"+s.Config.ServerPort, nil); err != nil {
		s.Logger.Fatal(err)
	}

	s.Logger.Info("Server stopped")
}
