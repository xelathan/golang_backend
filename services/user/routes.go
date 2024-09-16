package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/xelathan/golang_backend/config"
	"github.com/xelathan/golang_backend/services/auth"
	"github.com/xelathan/golang_backend/types"
	"github.com/xelathan/golang_backend/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/set_addresses", auth.WithJWTAuth(h.handleSetAddresses, h.store)).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	payload := types.LoginUserPayload{}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if user currently exists
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s doest not exist", payload.Email))
		return
	}

	if auth.CheckHashedPassword(payload.Password, user.Password) != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect password: %s", payload.Password))
		return
	}

	secret := []byte(config.Envs.JWTSecret)

	token, err := auth.CreateJWT(secret, user.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get json payload
	payload := types.RegisterUserPayload{}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if user currently exists
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil && err.Error() != "user not found" || user != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if does not exist create the new user
	err = h.store.CreateUser(
		types.User{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Password:  hashedPassword,
		},
	)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = utils.WriteJSON(w, http.StatusCreated, map[string]string{"registered": payload.Email})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) handleSetAddresses(w http.ResponseWriter, r *http.Request) {
	// read payload of setting addresses
	payload := types.SetAddressesPayload{}
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, err)
		return
	}

	// validate of the user has passed in valid JWT token with valid userId in claims
	userId := auth.GetUserIdFromContext(r.Context())

	if userId == -1 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid userId"))
		return
	}

	encryptedDefault, err := encryptAddress(payload.Default)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	encryptedSecondary, err := encryptAddress(payload.Secondary)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	encryptedTertiary, err := encryptAddress(payload.Tertiary)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userAddresses := types.UserAddresses{
		UserId:    userId,
		Default:   encryptedDefault,
		Secondary: encryptedSecondary,
		Tertiary:  encryptedTertiary,
	}

	// insert addresses into store
	if err := h.store.CreateUpdateAddress(&userAddresses); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"status": "success"})
}

func encryptAddress(address string) (string, error) {
	if address != "" {
		encrypted_address, err := EncryptAES(address, []byte(config.Envs.EncryptionKey))
		if err != nil {
			return "", err
		}

		return encrypted_address, nil
	}

	return "", nil
}
