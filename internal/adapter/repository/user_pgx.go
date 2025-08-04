package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/JacobD36/appfe_frontpage_api/internal/domain"
	ui "github.com/JacobD36/appfe_frontpage_api/internal/domain/interfaces"
	"github.com/JacobD36/appfe_frontpage_api/internal/usecase/interfaces"
	"github.com/jackc/pgx/v5"
)

const (
	pgxUserTableCreate = `
	CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL UNIQUE,
        password TEXT,
        img TEXT,
        role VARCHAR(50) DEFAULT 'CLIENT_ROLE',
        status BOOLEAN DEFAULT TRUE,
        email_validated BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ
    );`
	pgxUserCreate = `
	INSERT INTO users (name, email, password, img, role, status, email_validated, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    RETURNING id;
	`
	pgxUserGetAll          = `SELECT id, name, email, password, img, role, status, email_validated, created_at, updated_at FROM users ORDER BY created_at DESC;`
	pgxUserGetAllPaginated = `SELECT id, name, email, password, img, role, status, email_validated, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2;`
	pgxUserGetAllWithSearch = `SELECT id, name, email, password, img, role, status, email_validated, created_at, updated_at 
		FROM users 
		WHERE (LOWER(name) LIKE LOWER($1) OR LOWER(email) LIKE LOWER($1) OR LOWER(CONCAT(name, ' ', email)) LIKE LOWER($1))
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3;`
	pgxUserCount           = `SELECT COUNT(*) FROM users;`
	pgxUserCountWithSearch = `SELECT COUNT(*) FROM users 
		WHERE (LOWER(name) LIKE LOWER($1) OR LOWER(email) LIKE LOWER($1) OR LOWER(CONCAT(name, ' ', email)) LIKE LOWER($1));`
	pgxUserGetByID = `SELECT id, name, email, password, img, role, status, email_validated, created_at, updated_at
    FROM users
    WHERE id = $1;`
	pgxUserFindByEmail = `SELECT id, name, email, password, img, role, status, email_validated, created_at, updated_at
    FROM users
    WHERE email = $1;`
	pgxUserUpdate = `UPDATE users SET %s WHERE id = $%d;`
	pgxUserDetele = `UPDATE users
		SET status = false,
		    updated_at = $1
		WHERE id = $2;`
)

type pgxUserRepository struct {
	db pgx.Tx
}

func NewPgxUser(db pgx.Tx) ui.UserRepository {
	return &pgxUserRepository{db}
}

func (r *pgxUserRepository) Migrate(ctx context.Context) error {
	_, err := r.db.Exec(context.Background(), pgxUserTableCreate)
	return err
}

func (r *pgxUserRepository) Create(ctx context.Context, u *domain.User) error {
	err := r.db.QueryRow(ctx, pgxUserCreate,
		u.Name,
		u.Email,
		u.Password,
		u.Img,
		u.Role,
		u.Status,
		u.EmailValidated,
		u.CreatedAt,
	).Scan(&u.ID)

	return err
}

func (r *pgxUserRepository) GetAll(ctx context.Context, pagination *domain.Pagination) ([]*domain.User, int64, error) {
	var total int64
	var rows pgx.Rows
	var err error

	hasSearch := pagination != nil && pagination.Search != ""

	if hasSearch {
		searchTerm := "%" + pagination.Search + "%"
		err = r.db.QueryRow(ctx, pgxUserCountWithSearch, searchTerm).Scan(&total)
	} else {
		err = r.db.QueryRow(ctx, pgxUserCount).Scan(&total)
	}

	if err != nil {
		return nil, 0, err
	}

	if pagination == nil {
		rows, err = r.db.Query(ctx, pgxUserGetAll)
	} else if hasSearch {
		searchTerm := "%" + pagination.Search + "%"
		rows, err = r.db.Query(ctx, pgxUserGetAllWithSearch, searchTerm, pagination.Limit, pagination.Offset)
	} else {
		rows, err = r.db.Query(ctx, pgxUserGetAllPaginated, pagination.Limit, pagination.Offset)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *pgxUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, pgxUserGetByID, id)
	return scanUser(row)
}

func (r *pgxUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, pgxUserFindByEmail, email)
	return scanUser(row)
}

func (r *pgxUserRepository) UpdateByID(ctx context.Context, input ui.UpdateUserInput) error {
	fields := input.FieldsToUpdate()

	if len(fields) == 0 {
		return fmt.Errorf("no hay campos para actualizar para usuario con ID %s", input.GetID())
	}

	set := []string{}
	args := []interface{}{}
	i := 1

	for col, val := range fields {
		set = append(set, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	set = append(set, fmt.Sprintf("updated_at = $%d", i))
	args = append(args, time.Now())
	i++

	args = append(args, input.GetID())
	query := fmt.Sprintf(pgxUserUpdate, strings.Join(set, ", "), i)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *pgxUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, pgxUserDetele, time.Now(), id)
	return err
}

func scanUser(s interfaces.Scanner) (*domain.User, error) {
	var (
		password  *string
		img       *string
		updatedAt *time.Time
	)

	u := &domain.User{}

	err := s.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&password,
		&img,
		&u.Role,
		&u.Status,
		&u.EmailValidated,
		&u.CreatedAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	u.Password = password
	u.Img = img
	u.UpdatedAt = updatedAt

	return u, nil
}
