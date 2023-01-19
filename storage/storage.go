package storage

import (
	"app/models"
	"context"
)

type StorageI interface {
	CloseDB()
	Category() CategoryRepoI
	Product() ProductRepoI
	User() UserRepoI
}

type CategoryRepoI interface {
	Insert(context.Context, *models.CreateCategory) (string, error)
	GetByID(context.Context, *models.CategoryPrimeryKey) (*models.Category1, error)
	GetList(context.Context, *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error)
	Update(context.Context, *models.UpdateCategory) error
	Delete(context.Context, *models.CategoryPrimeryKey) error
}

type ProductRepoI interface {
	Insert(context.Context, *models.CreateProduct) (string, error)
	GetByID(context.Context, *models.ProductPrimeryKey) (*models.Product, error)
	GetList(context.Context, *models.GetListProductRequest) (*models.GetListProductResponse, error)
	Update(context.Context, *models.UpdateProduct) error
	Delete(context.Context, *models.ProductPrimeryKey) error
}

type UserRepoI interface {
	Insert(context.Context, *models.CreateUser) (string, error)
	GetByID(context.Context, *models.UserPrimaryKey) (*models.User, error)
	GetList(context.Context, *models.GetListUserRequest) (*models.GetListUserResponse, error)
	Update(context.Context, *models.UpdateUser) error
	Delete(context.Context, *models.UserPrimaryKey) error
}
