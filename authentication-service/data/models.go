package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)



const dbTimeout = 3 * time.Second

var db *sql.DB

func NewDatabase(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}

type Models struct {
	User User
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



func (user *User) GetAllUsers() ([]*User, error) {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM users
		ORDER BY last_name`

	rows, err := db.QueryContext(context, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var fetchedUser User
		err := rows.Scan(
			&fetchedUser.ID,
			&fetchedUser.Email,
			&fetchedUser.FirstName,
			&fetchedUser.LastName,
			&fetchedUser.Password,
			&fetchedUser.Active,
			&fetchedUser.CreatedAt,
			&fetchedUser.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning user:", err)
			return nil, err
		}

		users = append(users, &fetchedUser)
	}

	return users, nil
}

func (user *User) GetUserByEmail(email string) (*User, error) {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM users
		WHERE email = $1`

	var fetchedUser User
	row := db.QueryRowContext(context, query, email)

	err := row.Scan(
		&fetchedUser.ID,
		&fetchedUser.Email,
		&fetchedUser.FirstName,
		&fetchedUser.LastName,
		&fetchedUser.Password,
		&fetchedUser.Active,
		&fetchedUser.CreatedAt,
		&fetchedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &fetchedUser, nil
}

func (user *User) GetUserByID(id int) (*User, error) {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM users
		WHERE id = $1`

	var fetchedUser User
	row := db.QueryRowContext(context, query, id)

	err := row.Scan(
		&fetchedUser.ID,
		&fetchedUser.Email,
		&fetchedUser.FirstName,
		&fetchedUser.LastName,
		&fetchedUser.Password,
		&fetchedUser.Active,
		&fetchedUser.CreatedAt,
		&fetchedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &fetchedUser, nil
}

func (user *User) UpdateUser() error {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		UPDATE users
		SET email = $1, first_name = $2, last_name = $3, user_active = $4, updated_at = $5
		WHERE id = $6`

	_, err := db.ExecContext(context, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) DeleteUser() error {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	_, err := db.ExecContext(context, query, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) DeleteUserByID(id int) error {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	_, err := db.ExecContext(context, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) InsertUser(newUser User) (int, error) {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	query := `
		INSERT INTO users (email, first_name, last_name, password, user_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = db.QueryRowContext(context, query,
		newUser.Email,
		newUser.FirstName,
		newUser.LastName,
		encryptedPassword,
		newUser.Active,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (user *User) ResetUserPassword(newPassword string) error {
	context, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err = db.ExecContext(context, query, encryptedPassword, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) IsPasswordMatching(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// Password does not match
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
