package database

import (
	"errors"

	"vetsys/domain"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	DB *sqlx.DB
}

var ErrUserNotFound = errors.New("user not found")

func (userRepository *UserRepository) CreateUser(user *domain.User) error {
	query := `
	INSERT INTO users (dni, email, password, name, profile_picture)
	VALUES (:dni, :email, :password, :name, :profile_picture)
	RETURNING id
	`
	stmt, err := userRepository.DB.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.Get(&user.ID, user)
}

func (userRepository *UserRepository) GetUserByID(id int64) (*domain.User, error) {
	query := `SELECT id, dni, email, password, name, profile_picture FROM users WHERE id = $1`
	user := domain.User{}
	err := userRepository.DB.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *UserRepository) GetUserByDNI(dni string) (*domain.User, error) {
	query := `SELECT id, dni, email, password, name, profile_picture FROM users WHERE dni = $1`
	user := domain.User{}
	err := userRepository.DB.Get(&user, query, dni)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, dni, email, password, name, profile_picture FROM users WHERE email = $1`
	user := domain.User{}
	err := userRepository.DB.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *UserRepository) DeleteUserByID(id int64) error {
	result, err := userRepository.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (userRepository *UserRepository) UpdateUser(user *domain.User) error {
	query := "UPDATE users SET dni = $1, email = $2, password = $3, name = $4, profile_picture = $5 WHERE id = $6"
	result, err := userRepository.DB.Exec(query, user.DNI, user.Email, user.Password, user.Name, user.ProfilePicture, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
