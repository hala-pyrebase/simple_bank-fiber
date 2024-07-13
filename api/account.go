package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	db "github.com/saadhasan/simplebank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=USD CAD EUR"`
}

func (server *Server) createAccount(ctx *fiber.Ctx) error {
	var req createAccountRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("Enter proper parameters.")
	}

	if errs := server.validator.Validate(req); len(errs) > 0 && errs[0].HasError {
		errorMessages := make([]string, 0)

		for _, err2 := range errs {
			errorMessages = append(errorMessages, fmt.Sprintf("%s field has failed. Validation is: %s", err2.Field, err2.Tag))
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(strings.Join(errorMessages, " and that "))
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  int64(0),
		Currency: req.Currency,
	}

	account, err3 := server.store.CreateAccount(ctx.Context(), arg)
	if err3 != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err3.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Account is created successfully!", "account": account})
}

type getAccountRequest struct {
	ID int64 `params:"id" validate:"required,min=1"`
}

func (server *Server) getAccount(ctx *fiber.Ctx) error {
	var req getAccountRequest
	err := ctx.ParamsParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON("Enter proper parameters.")
	}

	if errs := server.validator.Validate(req); len(errs) > 0 && errs[0].HasError {
		errorMessages := make([]string, 0)

		for _, err2 := range errs {
			errorMessages = append(errorMessages, fmt.Sprintf("%s field has failed. Validation is: %s", err2.Field, err2.Tag))
		}

		return ctx.Status(fiber.StatusBadRequest).JSON(strings.Join(errorMessages, " and that "))
	}

	account, err3 := server.store.GetAccount(ctx.Context(), req.ID)
	if err3 != nil {
		if err3 == sql.ErrNoRows {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err3.Error()})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err3.Error()})
	}

	return ctx.Status(fiber.StatusFound).JSON(fiber.Map{"message": "Account is found!", "account": account})
}
