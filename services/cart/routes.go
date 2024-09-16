package cart

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xelathan/golang_backend/services/auth"
	"github.com/xelathan/golang_backend/types"
	"github.com/xelathan/golang_backend/utils"
)

type Handler struct {
	orderStore   types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
	db           *sql.DB
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore, userStore types.UserStore, db *sql.DB) *Handler {
	return &Handler{
		orderStore:   orderStore,
		productStore: productStore,
		userStore:    userStore,
		db:           db,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userId := auth.GetUserIdFromContext(r.Context())
	if userId == -1 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid userId"))
		return
	}

	cart_payload := types.CartCheckoutPayload{}

	if err := utils.ParseJSON(r, &cart_payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart_payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ids, err := getCartItemsIDs(cart_payload.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// get products
	products, err := h.productStore.GetProductsByID(ids)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	orderId, totalPrice, err := h.createOrder(products, cart_payload.Items, userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"orderId": orderId, "totalPrice": totalPrice})
}
