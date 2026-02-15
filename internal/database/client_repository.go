package database

import (
	"database/sql"
	"errors"
	"vetsys/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ClientRepository struct {
	DB *sqlx.DB
}

var ErrClientNotFound = errors.New("Client not found")

func (clientRepository *ClientRepository) CreateClient(client *domain.Client) error {
	query := `INSERT INTO clients (dni, name, phone_number) VALUES (:dni, :name, :phone_number) RETURNING id`
	stmt, err := clientRepository.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(&client.ID, client)
}

func (clientRepository *ClientRepository) GetClientByID(id int64) (*domain.Client, error) {
	var client domain.Client
	err := clientRepository.DB.Get(&client, "SELECT id, dni, name, phone_number FROM clients WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (clientRepository *ClientRepository) GetClientByDNI(dni string) (*domain.Client, error) {
	var client domain.Client
	err := clientRepository.DB.Get(&client, "SELECT id, dni, name, phone_number FROM clients WHERE dni = $1", dni)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (clientRepository *ClientRepository) UpdateClient(client *domain.Client) error {
	query := `UPDATE clients SET dni = $1, name = $2, phone_number = $3 WHERE id = $4`
	result, err := clientRepository.DB.Exec(query, client.DNI, client.Name, client.PhoneNumber, client.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrClientNotFound
	}
	return nil
}

func (clientRepository *ClientRepository) DeleteClientByID(id int64) error {
	result, err := clientRepository.DB.Exec("DELETE FROM clients WHERE id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrClientNotFound
	}
	return nil
}
