package server

import (
	"invoices/config"
	"invoices/controllers/v1"
	"invoices/handlers"

	"github.com/gin-gonic/gin"
)

// Server is the http layer for role and user resource
type Server struct {
	config            *config.Config
	invoiceController *controllers.InvoiceController
	handlers          *handlers.HttpHandlers
}

// NewServer is the Server constructor
func NewServer(cf *config.Config,
	pc *controllers.InvoiceController,
	hand *handlers.HttpHandlers) *Server {

	return &Server{
		config:            cf,
		invoiceController: pc,
		handlers:          hand,
	}
}

// Run loads server with its routes and starts the server
func (s *Server) Run() {
	// Instantiate a new router
	r := gin.Default()

	// generic routes
	r.HandleMethodNotAllowed = false
	r.NoRoute(s.handlers.NotFound)

	// Invoice resource
	invoiceApi := r.Group("/api/v1/invoice")
	{
		// Create a new invoice
		invoiceApi.POST("", s.invoiceController.CreateAction)

		// List invoices with filtering and pagination
		invoiceApi.GET("", s.invoiceController.ListAction)
	}

	// Fire up the server
	r.Run(s.config.Host)
}
