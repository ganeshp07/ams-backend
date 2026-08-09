package repositories

import (
	"ams-backend/models"
	"database/sql"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository { return &UserRepository{} }

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := DB.QueryRow(
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	var u models.User
	err := DB.QueryRow(
		`SELECT id, email, password_hash, role, is_active, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) CreateWithTx(tx *sql.Tx, email, hash, role string) (int64, error) {
	var id int64
	err := tx.QueryRow(
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id`,
		email, hash, role,
	).Scan(&id)
	return id, err
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&count)
	return count > 0, err
}
