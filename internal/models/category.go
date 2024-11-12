package models

type Category struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

type CategoryCreateRequest struct {
	Name   string `json:"name" binding:"required"`
	Color  string `json:"color"`
	Icon   string `json:"icon"`
	UserID int32  `json:"user_id"`
}

type CategoryUpdateRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
	Icon  *string `json:"icon"`
}
