package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/xelathan/golang_backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]types.Product, 0)

	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Price, &product.Quantity, &product.CreatedAt)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?,?,?,?,?)", product.Name, product.Description, product.Image, product.Price, product.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetProductsByID(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, price = ?, image = ?, description = ?, quantity = ? WHERE id = ?", product.Name, product.Price, product.Image, product.Description, product.Quantity, product.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateProductBatch(products map[int]types.Product) error {
	if len(products) == 0 {
		return nil
	}

	query := "UPDATE products SET quantity = CASE id"
	ids := []int{}

	for _, product := range products {
		query += fmt.Sprintf(" WHEN %d THEN %d", product.ID, product.Quantity)
		ids = append(ids, product.ID)
	}
	query += " END WHERE id IN ("
	for i := range len(products) {
		if i > 0 {
			query += ","
		}
		query += fmt.Sprintf("%d", ids[i])
	}
	query += ");"

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
