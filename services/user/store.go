package user

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

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)

	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserById(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?,?,?,?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUserAddressesByUserId(id int) ([]types.UserAddresses, error) {
	rows, err := s.db.Query("SELECT * FROM user_addresses WHERE userId = ?", id)
	if err != nil {
		return nil, err
	}

	addresses := []types.UserAddresses{}
	for rows.Next() {
		address, err := scanRowIntoUserAddresses(rows)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, *address)
	}

	return addresses, nil
}

func scanRowIntoUserAddresses(rows *sql.Rows) (*types.UserAddresses, error) {
	user_addresses := new(types.UserAddresses)
	err := rows.Scan(&user_addresses.Id, &user_addresses.UserId, &user_addresses.AddressType, &user_addresses.Address)
	if err != nil {
		return nil, err
	}

	return user_addresses, nil
}

func (s *Store) CreateUpdateAddress(addresses *types.UserAddresses) error {
	query := "INSERT INTO user_addresses (userId, address_type, address) VALUES (?,?,?) ON DUPLICATE KEY UPDATE address = VALUES(address)"
	_, err := s.db.Exec(query, addresses.UserId, addresses.AddressType, addresses.Address)
	if err != nil {
		return err
	}

	return nil
}
