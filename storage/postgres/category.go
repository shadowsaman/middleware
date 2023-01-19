package postgres

import (
	"app/models"
	"app/pkg/helper"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (r *CategoryRepo) Insert(ctx context.Context, category *models.CreateCategory) (string, error) {

	var (
		id = uuid.New().String()
	)

	query := `
		INSERT INTO categories (
			id,
			name,
			parent_id,
			photo_url,
			updated_at
		) VALUES ($1, $2, $3, $4, now())
	`
	_, err := r.db.Exec(ctx, query,
		id,
		category.Name,
		helper.NullString(category.ParentId),
		category.PhotoUrl,
	)

	if err != nil {
		return "", err
	}

	return id, nil

}

func (r *CategoryRepo) GetByID(ctx context.Context, category *models.CategoryPrimeryKey) (*models.Category1, error) {

	var (
		id        sql.NullString
		name      sql.NullString
		parentID  sql.NullString
		photoUrl  sql.NullString
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query := `
		SELECT 
			id,
			name,
			parent_id,
			photo_url,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1	
	`
	err := r.db.QueryRow(ctx, query, category.Id).Scan(
		&id,
		&name,
		&parentID,
		&photoUrl,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	resp := &models.Category1{
		Id:        id.String,
		Name:      name.String,
		ParentId:  parentID.String,
		PhotoUrl:  photoUrl.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}

	queryChild := `
		SELECT
			id,
			name,
			parent_id,
			photo_url,
			created_at,
			updated_at
		FROM categories
		WHERE parent_id = $1
	`
	rows, err := r.db.Query(ctx, queryChild, category.Id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&id,
			&name,
			&parentID,
			&photoUrl,
			&createdAt,
			&updatedAt,
		)

		resp.Childs = append(resp.Childs, &models.Category{
			Id:        id.String,
			Name:      name.String,
			ParentId:  parentID.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	return resp, err
}

func (r *CategoryRepo) GetList(ctx context.Context, req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error) {

	var (
		offset = "OFFSET 0"
		limit  = "LIMIT 10"
		resp   = &models.GetListCategoryResponse{}
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf("OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf("LIMIT %d", req.Limit)
	}

	query := `
		SELECT 
			COUNT(*) OVER(),
			id,
			name,
			parent_id,
			photo_url,
			created_at,
			updated_at
		FROM categories
		WHERE parent_id IS NULL
	`

	query += offset + limit

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var (
			id        sql.NullString
			name      sql.NullString
			parentID  sql.NullString
			photoUrl  sql.NullString
			createdAt sql.NullString
			updatedAt sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&parentID,
			&photoUrl,
			&createdAt,
			&updatedAt,
		)

		resp.Categories = append(resp.Categories, &models.Category1{
			Id:        id.String,
			Name:      name.String,
			ParentId:  parentID.String,
			PhotoUrl:  photoUrl.String,
			CreatedAt: createdAt.String,
			UpdatedAt: updatedAt.String,
		})
	}

	for _, category := range resp.Categories {

		queryChild := `
			SELECT
				id,
				name,
				parent_id,
				photo_url,
				created_at,
				updated_at
			FROM categories
			WHERE parent_id = $1
		`
		rows, err := r.db.Query(ctx, queryChild, category.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return resp, nil
			}

			return nil, err
		}

		for rows.Next() {

			var (
				id        sql.NullString
				name      sql.NullString
				parentID  sql.NullString
				photoUrl  sql.NullString
				createdAt sql.NullString
				updatedAt sql.NullString
			)

			err = rows.Scan(
				&id,
				&name,
				&parentID,
				&photoUrl,
				&createdAt,
				&updatedAt,
			)

			category.Childs = append(category.Childs, &models.Category{
				Id:        id.String,
				Name:      name.String,
				ParentId:  parentID.String,
				PhotoUrl:  photoUrl.String,
				CreatedAt: createdAt.String,
				UpdatedAt: updatedAt.String,
			})
		}
	}

	return resp, err
}

func (r *CategoryRepo) Update(ctx context.Context, category *models.UpdateCategory) error {
	query := `
		UPDATE
			categories
		SET
			name = $2,
			parent_id = $3,
			photo_url = $4,
			updated_at = now()
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		category.Id,
		category.Name,
		category.ParentId,
		category.PhotoUrl,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepo) Delete(ctx context.Context, req *models.CategoryPrimeryKey) error {

	_, err := r.db.Exec(ctx, "delete from categories where id = $1", req.Id)
	if err != nil {
		return err
	}

	return nil
}
