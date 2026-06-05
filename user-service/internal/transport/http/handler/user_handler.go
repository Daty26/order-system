package transport_http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Daty26/order-system/user-service/internal/service"
	transport_http_dto "github.com/Daty26/order-system/user-service/internal/transport/http/dto"
	transport_http_response "github.com/Daty26/order-system/user-service/internal/transport/http/response"
)

type UserHandler struct{
	service *service.UserService
}
func NewUserHandler(service *service.UserService) *UserHandler{
	return &UserHandler{service: service}
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request){
	var request transport_http_dto.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil{
		transport_http_response.ErrorJSON(
			w, 
			http.StatusBadRequest, 
			fmt.Sprintf("invalid request body: %s", err.Error()),	
		)
		return
	}
	input := service.LoginUserInput{
		Identifier: request.Identifier,
		Password: request.Password,
	}
	userSummaryModel,token,  err := h.service.LoginUser(r.Context(), input)
	if err != nil{
		if errors.Is(err, service.ErrInvalidCredentials){
			transport_http_response.ErrorJSON(
				w,
				http.StatusUnauthorized,
				"invalid credentials",
			)
			return 
		}
		transport_http_response.ErrorJSON(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return 
	}
	resp := transport_http_dto.LoginUserResponse{
		User: userSummaryModel,
		Token: token,
	}
	transport_http_response.SuccessJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request){
	var request transport_http_dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil{
		transport_http_response.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return 
	}
	input := service.CreateUserInput{
		Username: request.Username,
		Email: request.Email,
		Password: request.Password,
	}
	userSummary, err :=h.service.CreateUser(r.Context(), input)
	if err != nil{
		if errors.Is(err, service.ErrInvalidUserInput){
			transport_http_response.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrUserAlreadyExists){
			transport_http_response.ErrorJSON(w, http.StatusConflict, err.Error())
			return	
		}
		transport_http_response.ErrorJSON(w, http.StatusInternalServerError, "something went wrong")
		return	
	}
	transport_http_response.SuccessJSON(w, http.StatusCreated, userSummary)
}
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request){
	id := r.Context().Value("user_id").(int)
	user, err :=h.service.GetByID(r.Context(), id)
	if err != nil{
		if errors.Is(err, service.ErrNotFound){
			transport_http_response.ErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}
		transport_http_response.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	transport_http_response.SuccessJSON(w, http.StatusOK, user)
}