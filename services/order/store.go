package order

import (
	"database/sql"
	"fmt"

	"github.com/xelathan/golang_backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?,?,?,?)", order.UserId, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?,?,?,?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateOrder(order types.Order) error {
	_, err := s.db.Exec("UPDATE orders SET total = ?, status = ?, address = ? WHERE id = ?", order.Total, order.Status, order.Address, order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetOrderHistoryByUserId(userId int) ([]types.OrderHistory, error) {
	rows, err := s.db.Query("SELECT o.id, o.total, o.status, o.address, o.createdAt, oi.productId, oi.quantity, oi.price FROM orders o JOIN order_items oi ON o.id = oi.orderId WHERE o.userId = ? ORDER BY o.createdAt DESC", userId)
	if err != nil {
		return nil, err
	}

	orders := []types.OrderHistory{}

	for rows.Next() {
		order, err := scanRowsIntoOrderHistory(rows)
		if err != nil {
			return nil, err
		}

		orders = append(orders, *order)
	}

	return orders, nil
}

func scanRowsIntoOrderHistory(rows *sql.Rows) (*types.OrderHistory, error) {
	orderHistoryRow := new(types.OrderHistory)

	err := rows.Scan(&orderHistoryRow.OrderId, &orderHistoryRow.Total, &orderHistoryRow.Status, &orderHistoryRow.Address, &orderHistoryRow.CreatedAt, &orderHistoryRow.ProductId, &orderHistoryRow.Quantity, &orderHistoryRow.Price)
	if err != nil {
		return nil, err
	}

	return orderHistoryRow, nil
}

func (s *Store) GetOrderById(orderId int) (*types.Order, error) {
	rows, err := s.db.Query("SELECT * FROM orders WHERE orderId = ?", orderId)
	if err != nil {
		return nil, err
	}

	order := new(types.Order)
	for rows.Next() {
		order, err = scanRowsIntoOrder(rows)
		if err != nil {
			return nil, err
		}
	}

	if order.ID == 0 {
		return nil, fmt.Errorf("order does not exist")
	}

	return order, nil
}

func scanRowsIntoOrder(rows *sql.Rows) (*types.Order, error) {
	order := new(types.Order)
	err := rows.Scan(&order.ID, &order.UserId, &order.Total, &order.Status, &order.Address, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	return order, nil
}
