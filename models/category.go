package models

type CategoryPrimeryKey struct {
	Id string `json:"id"`
}

type CreateCategory struct {
	Name     string `json:"name"`
	ParentId string `json:"parent_id"`
	PhotoUrl string `json:"photo_url"`
}

type Category struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ParentId  string `json:"parent_id"`
	PhotoUrl  string `json:"photo_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Category1 struct {
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	ParentId  string      `json:"parent_id"`
	PhotoUrl  string      `json:"photo_url"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Childs    []*Category `json:"childs"`
}

type UpdateCategory struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ParentId string `json:"parent_id"`
	PhotoUrl string `json:"photo_url"`
}

type UpdateCategorySwag struct {
	Name     string `json:"name"`
	ParentId string `json:"parent_id"`
	PhotoUrl string `json:"photo_url"`
}

type GetListCategoryRequest struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

type GetListCategoryResponse struct {
	Count      int64        `json:"count"`
	Categories []*Category1 `json:"categories"`
}

type Empty struct{}
