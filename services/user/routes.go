package user

import (
	"fmt"
	"net/http"

	"github.com/evinicius16/ecommerce-api.git/services/auth"
	"github.com/evinicius16/ecommerce-api.git/types"
	"github.com/evinicius16/ecommerce-api.git/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/user/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(response http.ResponseWriter, request *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(request, &payload); err != nil {
		utils.WriteError(response, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(response, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(response, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(response, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(response, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(response, http.StatusCreated, nil)
}
