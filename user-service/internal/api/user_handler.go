package api

import (
	"encoding/json"
	"github.com/Daty26/order-system/user-service/internal/model"
	"github.com/Daty26/order-system/user-service/internal/service"
	"net/http"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{srv: service}
}
func (uh *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string      `json:"username"`
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Role     model.Roles `json:"role"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Wrong request: "+err.Error())
		return
	}
	user := model.User{Username: req.Username, Email: req.Email, Password: req.Password, Role: req.Role}

	createdUser, err := uh.srv.CreateUser(user)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Couldn't create user: "+err.Error())
		return
	}
	SuccessResponse(w, http.StatusCreated, createdUser)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	userLoggedIn, token, err := uh.srv.Login(req.Identifier, req.Password)
	var resp = struct {
		User  model.User `json:"user"`
		Token string     `json:"token"`
	}{}
	resp.User = userLoggedIn
	resp.Token = token

	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SuccessResponse(w, http.StatusOK, resp)

}
