package postgres

import (
	"app/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Insert(ctx context.Context, req *models.CreateUser) (string, error) {

	var (
		id = uuid.New().String()
	)

	query := `
	INSERT INTO users (
			id,
			name,
			login,
			password,
			updated_at
		) VALUES ($1, $2, $3, $4, now())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Login,
		req.Password,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *UserRepo) GetByID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error) {

	query := `
		SELECT 
			id,
			name,
			login,
			password,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var (
		id       sql.NullString
		name     sql.NullString
		login    sql.NullString
		password sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).
		Scan(
			&id,
			&name,
			&login,
			&password,
		)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Id:       id.String,
		Name:     name.String,
		Login:    login.String,
		Password: password.String,
	}

	return user, nil
}

func (r *UserRepo) GetList(ctx context.Context, req *models.GetListUserRequest) (*models.GetListUserResponse, error) {

	var (
		resp   models.GetListUserResponse
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			login,
			password,
			created_at,
			updated_at 
		FROM users
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += offset + limit

	rows, err := r.db.Query(ctx, query)

	defer rows.Close()

	if err != nil {
		return &models.GetListUserResponse{}, err
	}

	var (
		id       sql.NullString
		name     sql.NullString
		login    sql.NullString
		password sql.NullString
	)

	for rows.Next() {

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&login,
			&password,
		)
		user := models.User{
			Id:       id.String,
			Name:     name.String,
			Login:    login.String,
			Password: password.String,
		}
		if err != nil {
			return &models.GetListUserResponse{}, err
		}

		resp.Users = append(resp.Users, &user)

	}
	return &resp, nil
}

func (r *UserRepo) Update(ctx context.Context, user *models.UpdateUser) error {
	query := `
		UPDATE 
			users
		SET 
			name = $2,
			login = $3,
			password = $3
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		user.Id,
		user.Name,
		user.Login,
		user.Password,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, req *models.UserPrimaryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", req.Id)

	if err != nil {
		return err
	}

	return nil
}
