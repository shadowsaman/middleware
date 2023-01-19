package models

type ProductPrimeryKey struct {
	Id string `json:"id"`
}

type CreateProduct struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	PhotoUrl    string  `json:"photo_url"`
	CategoryId  string  `json:"category_id"`
}

type Product struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	PhotoUrl    string   `json:"photo_url"`
	Category    Category `json:"category"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type Product1 struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	PhotoUrl    string  `json:"photo_url"`
	CategoryId  string  `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UpdateProduct struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CategoryId  string  `json:"category_id"`
	PhotoUrl    string  `json:"photo_url"`
}

type UpdateProductSwag struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	CategoryId  string  `json:"category_id"`
	PhotoUrl    string  `json:"photo_url"`
}

type GetListProductRequest struct {
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
	Search string `json:"search"`
}

type GetListProductResponse struct {
	Count    int64       `json:"count"`
	Products []*Product1 `json:"products"`
}
