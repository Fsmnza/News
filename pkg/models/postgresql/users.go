package mysql

import (
	"alexedwards.net/snippetbox/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `
        INSERT INTO users (name, email, hashed_password, created, role) VALUES($1, $2, $3,CURRENT_TIMESTAMP, $4)`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword), models.RoleUser)
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) {
			if pgError.Code == "23505" && strings.Contains(pgError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = $1 "
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}
func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := "Select id, email, hashed_password, created, role from users where id =$1"
	row := m.DB.QueryRow(stmt, id)
	d := &models.User{}
	err := row.Scan(&d.ID, &d.Email, &d.HashedPassword, &d.Created, &d.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return d, nil
}

func (m *UserModel) AdminGet() ([]*models.User, error) {
	stmt := "Select * from users"
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userList := []*models.User{}
	for rows.Next() {
		n := &models.User{}
		err = rows.Scan(&n.ID, &n.Name, &n.Email, &n.HashedPassword, &n.Created, &n.Role)
		if err != nil {
			return nil, err
		}
		userList = append(userList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userList, nil
}
func (m *UserModel) Role(userID, newRole string) error {
	if userID == "" {
		return errors.New("userID is empty")
	}
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("unable to convert userID to integer: %w", err)
	}
	stmt := `UPDATE users SET role = $1 WHERE id = $2`
	_, err = m.DB.Exec(stmt, newRole, userIDInt)
	if err != nil {
		return fmt.Errorf("unable to update user role: %w", err)
	}
	return nil
}
