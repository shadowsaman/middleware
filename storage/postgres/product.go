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

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (r *ProductRepo) Insert(ctx context.Context, req *models.CreateProduct) (string, error) {

	var (
		id = uuid.New().String()
	)

	query := `
	INSERT INTO product (
			id,
			name,
			price,
			description,
			category_id,
			photo_url,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, now())
	`

	_, err := r.db.Exec(ctx, query,
		id,
		req.Name,
		req.Price,
		req.Description,
		helper.NullString(req.CategoryId),
		req.PhotoUrl,
	)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, req *models.ProductPrimeryKey) (*models.Product, error) {

	query := `
		SELECT 
			id,
			name,
			price,
			description,
			photo_url,
			created_at,
			updated_at
		FROM product 
		WHERE id = $1
	`

	queryCategory := `
		SELECT 
			c.id,
			c.name,
			c.photo_url,
		FROM product AS p
		JOIN categories AS c ON c.id = p.category_id
		WHERE p.id = $1
	`

	var (
		id          sql.NullString
		name        sql.NullString
		price       sql.NullFloat64
		description sql.NullString
		photoUrl    sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	err := r.db.QueryRow(ctx, query, req.Id).
		Scan(
			&id,
			&name,
			&price,
			&description,
			&photoUrl,
			&createdAt,
			&updatedAt,
		)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		Id:          id.String,
		Name:        name.String,
		Price:       price.Float64,
		Description: description.String,
		PhotoUrl:    photoUrl.String,
		CreatedAt:   createdAt.String,
		UpdatedAt:   updatedAt.String,
	}
	var category models.Category
	err = r.db.QueryRow(ctx, queryCategory, req.Id).Scan(
		&category.Id,
		&category.Name,
		&category.PhotoUrl,
	)

	product.Category = category

	return product, nil
}

func (r *ProductRepo) GetList(ctx context.Context, req *models.GetListProductRequest) (*models.GetListProductResponse, error) {

	var (
		resp   models.GetListProductResponse
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		search = req.Search
		f      = "%"
	)

	query := `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			price,
			description,
			category_id,
			photo_url,
			created_at,
			updated_at 
		FROM product
	`
	if search != "" {
		search = fmt.Sprintf("WHERE name LIKE  '%s%s' ", req.Search, f)
		query += search
	}
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
		return &models.GetListProductResponse{}, err
	}

	var (
		id          sql.NullString
		name        sql.NullString
		price       sql.NullFloat64
		description sql.NullString
		categoryId  sql.NullString
		photoUrl    sql.NullString
		createdAt   sql.NullString
		updatedAt   sql.NullString
	)

	for rows.Next() {

		err = rows.Scan(
			&resp.Count,
			&id,
			&name,
			&price,
			&description,
			&categoryId,
			&photoUrl,
			&createdAt,
			&updatedAt,
		)
		product := models.Product1{
			Id:          id.String,
			Name:        name.String,
			Price:       price.Float64,
			Description: description.String,
			CategoryId:  categoryId.String,
			PhotoUrl:    photoUrl.String,
			CreatedAt:   createdAt.String,
			UpdatedAt:   updatedAt.String,
		}
		if err != nil {
			return &models.GetListProductResponse{}, err
		}

		resp.Products = append(resp.Products, &product)

	}
	return &resp, nil
}

func (r *ProductRepo) Update(ctx context.Context, product *models.UpdateProduct) error {
	query := `
		UPDATE 
			product
		SET 
			name = $2,
			price = $3,
			description = $4,
			category_id = $5
			photoUrl = $6,
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		product.Id,
		product.Name,
		product.Price,
		product.Description,
		product.CategoryId,
		product.PhotoUrl,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, req *models.ProductPrimeryKey) error {

	_, err := r.db.Exec(ctx, "DELETE FROM product WHERE id = $1", req.Id)

	if err != nil {
		return err
	}

	return nil
}
