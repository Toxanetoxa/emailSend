package handler

import (
	resp "email-sendler/internal/libs/api"
	"email-sendler/internal/logger"
	"email-sendler/internal/redis"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

type Message struct {
	To      string `json:"to" validate:"required,email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Request struct {
	Messages []Message `json:"messages"`
}

type Response struct {
	resp.Response
	Emails []string `json:"emails"`
}

func New(logger *logger.File, que *redis.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			responseErr(w, r, http.StatusBadRequest, "Invalid request")
			return
		}

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			logger.Error(op, err)
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		//responseOK(w, r, validEmails)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, emails []string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Emails:   emails,
	})
}

func responseErr(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, resp.Error(message))
}
