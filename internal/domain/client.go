package domain

type Client struct {
	ID          int64  `json:"id" db:"id"`
	DNI         string `json:"dni" db:"dni"`
	Name        string `json:"name" db:"name"`
	PhoneNumber string `json:"phoneNumber" db:"phone_number"`
}

func NewClient(dni string, name string, phoneNumber string) *Client {
	return &Client{
		DNI:         dni,
		Name:        name,
		PhoneNumber: phoneNumber,
	}
}
