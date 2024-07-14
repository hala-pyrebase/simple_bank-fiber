package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	db "github.com/saadhasan/simplebank/db/sqlc"
)

// ValidationError provides formated struct for passing error between context.
type ValidationError struct {
	HasError bool
	Field    string
	Tag      string
	Value    interface{}
}

type CustomValidatior struct {
	validator *validator.Validate
}

// Server serves HTTP requests for our banking service.
type Server struct {
	store     *db.Store
	router    *fiber.App
	validator *CustomValidatior
}

var validate = validator.New()

func (cv CustomValidatior) Validate(data interface{}) []ValidationError {
	var validationErrors []ValidationError

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var ve ValidationError

			ve.Value = err.Value()
			ve.Field = err.Field()
			ve.Tag = err.Tag()
			ve.HasError = true

			validationErrors = append(validationErrors, ve)
		}
	}

	return validationErrors
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := fiber.New()

	customValidator := &CustomValidatior{
		validator: validate,
	}

	server.validator = customValidator

	router.Post("/accounts", server.createAccount)
	router.Get("/accounts/:id", server.getAccount)
	router.Get("/accounts", server.listAccounts)
	router.Delete("/accounts/:id", server.deleteAccount)
	router.Put("/accounts", server.updateAccount)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Listen(address)
}
