package order

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
	userStore    types.UserStore
	orderStore   types.OrderStore
	productStore types.ProductStore
	db           *sql.DB
}

func NewHandler(userStore types.UserStore, orderStore types.OrderStore, productStore types.ProductStore, db *sql.DB) *Handler {
	return &Handler{userStore: userStore, orderStore: orderStore, productStore: productStore, db: db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cancel_order", auth.WithJWTAuth(h.handleCancelOrder, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
	// retrieve user id from context from token claims
	userId := auth.GetUserIdFromContext(r.Context())
	if userId == -1 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid userId in token claims"))
		return
	}

	// retrieve request body: needs orderId to cancel order
	cancelOrderPayload := types.CancelOrderPayload{}
	if err := utils.ParseJSON(r, &cancelOrderPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	tx, err := h.db.Begin()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	// get history of orders from user
	orderHistory, err := h.orderStore.GetOrderHistoryByUserId(userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	if len(orderHistory) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("no orders to cancel"))
		tx.Rollback()
		return
	}

	// create map of orderId: orderHistoryItem
	orderHistoryMap := map[int][]types.OrderHistory{}
	for _, or := range orderHistory {
		orderHistoryMap[or.OrderId] = append(orderHistoryMap[or.OrderId], or)
	}

	// get orderItems of the orderId given in payload
	orderItems := orderHistoryMap[cancelOrderPayload.OrderId]
	if len(orderItems) == 0 {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid order id"))
		tx.Rollback()
		return
	}

	// if order is not pending then we cancel
	if orderItems[0].Status != "pending" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("cannot cancel an order that is not pending"))
		tx.Rollback()
		return
	}

	// create list of productIds to query products and map of new quantities
	newItemQuantities := map[int]int{}
	productIds := []int{}
	for _, item := range orderItems {
		newItemQuantities[item.ProductId] = item.Quantity
		productIds = append(productIds, item.ProductId)
	}

	products, err := h.productStore.GetProductsByID(productIds)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	// update products based on new quantities: restock since order cancelled
	productsMap := map[int]types.Product{}
	for _, product := range products {
		product.Quantity += newItemQuantities[product.ID]
		productsMap[product.ID] = product
	}

	if err := h.productStore.UpdateProductBatch(productsMap); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	// set order status to cancelled
	if err := h.orderStore.UpdateOrder(types.Order{
		ID:        orderItems[0].OrderId,
		UserId:    userId,
		Total:     orderItems[0].Total,
		Status:    types.Cancelled,
		Address:   orderItems[0].Address,
		CreatedAt: orderItems[0].CreatedAt,
	}); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, err)
		tx.Rollback()
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "lets go baby"})
}
