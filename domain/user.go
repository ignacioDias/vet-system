package domain

type User struct {
	ID             int64  `db:"id" json:"id"`
	DNI            string `db:"dni" json:"dni"`
	Email          string `db:"email" json:"email"`
	Password       string `db:"password" json:"-"`
	Name           string `db:"name" json:"name"`
	ProfilePicture string `db:"profile_picture" json:"profilePicture"`
}

func NewUser(dni string, email string, password string, name string, profilePicture string) *User {
	return &User{
		DNI:            dni,
		Email:          email,
		Password:       password,
		Name:           name,
		ProfilePicture: profilePicture,
	}
}
