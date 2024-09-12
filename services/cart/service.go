package cart

import (
	"fmt"

	"github.com/xelathan/golang_backend/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product")
		}

		productIDs[i] = item.ProductID
	}

	return productIDs, nil
}

func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	// check if all products in stock
	// if in stock calculate total price
	// calculate total price
	// reduce quantity of product
	// create the order
	// create the order items

	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	totalPrice := calculateTotalPrice(items, productMap)

	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	orderId, err := h.orderStore.CreateOrder(types.Order{
		UserId: userID,
		Total:  totalPrice,
		Status: "pending",
		// TODO: create user address table
		Address: "some address",
	})

	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderId,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderId, totalPrice, nil
}

func checkIfCartIsInStock(items []types.CartItem, productsMap map[int]types.Product) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range items {
		product, ok := productsMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d is not available for the requested quantity", item.ProductID)
		}
	}

	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	totalPrice := 0.0
	for _, item := range items {
		totalPrice += productMap[item.ProductID].Price * float64(item.Quantity)
	}

	return totalPrice
}
