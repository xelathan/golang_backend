package types

import "time"

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=16"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=16"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(User) error
}

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image" validate:"required"`
	Price       float64 `json:"price" validate:"required,decimal2"`
	Quantity    int     `json:"quantity"`
}

type Product struct {
	ID          int       `json:"id"`
	Quantity    int       `json:"quantity"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProductStore interface {
	GetProducts() ([]Product, error)
	GetProductsByID(productIDs []int) ([]Product, error)
	CreateProduct(Product) error
	UpdateProduct(Product) error
}

type Order struct {
	ID        int       `json:"id"`
	UserId    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}

type CartItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}
