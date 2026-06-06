package transport_http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Daty26/order-system/user-service/internal/service"
	transport_http_dto "github.com/Daty26/order-system/user-service/internal/transport/http/dto"
	transport_http_response "github.com/Daty26/order-system/user-service/internal/transport/http/response"
)

type UserHandler struct {
	service *service.UserService
	logger  *slog.Logger
}

func NewUserHandler(service *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var request transport_http_dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.WarnContext(
			r.Context(),
			"invalid login request body",
			"error", err,
		)

		transport_http_response.ErrorJSON(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}
	input := service.LoginUserInput{
		Identifier: request.Identifier,
		Password:   request.Password,
	}
	userSummaryModel, token, err := h.service.LoginUser(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			transport_http_response.ErrorJSON(
				w,
				http.StatusUnauthorized,
				"invalid credentials",
			)
			return
		}
		h.logger.ErrorContext(
			r.Context(),
			"failed to login user",
			"error", err,
		)

		transport_http_response.ErrorJSON(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return
	}
	resp := transport_http_dto.LoginUserResponse{
		User:  userSummaryModel,
		Token: token,
	}
	transport_http_response.SuccessJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request transport_http_dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.WarnContext(r.Context(),
			"invalid user req body",
			"error", err,
		)
		transport_http_response.ErrorJSON(w, http.StatusBadRequest, "invalid req body")
		return
	}
	input := service.CreateUserInput{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
	userSummary, err := h.service.CreateUser(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidUserInput) {
			transport_http_response.ErrorJSON(w, http.StatusBadRequest, "invalid req body")
			return
		}
		if errors.Is(err, service.ErrUserAlreadyExists) {
			transport_http_response.ErrorJSON(w, http.StatusConflict, "user already exists")
			return
		}
		h.logger.ErrorContext(r.Context(), "failed to create user", "error", err)
		transport_http_response.ErrorJSON(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	transport_http_response.SuccessJSON(w, http.StatusCreated, userSummary)
}
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user_id").(int)
	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			transport_http_response.ErrorJSON(w, http.StatusNotFound, "user not found")
			return
		}
		h.logger.ErrorContext(r.Context(), "failed to get me", "error", err)
		transport_http_response.ErrorJSON(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	transport_http_response.SuccessJSON(w, http.StatusOK, user)
}
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, offset, ok := parsePagination(r)
	if !ok {
		transport_http_response.ErrorJSON(w, http.StatusBadRequest, "invalid pagination params")
		return
	}
	users, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to get all users", "error", err)
		transport_http_response.ErrorJSON(w, http.StatusBadRequest, "something went wrong")
		return
	}
	transport_http_response.SuccessJSON(w, http.StatusOK, users)
}
