package restaurants

import "time"

type Restaurant struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Address     *string   `json:"address,omitempty"`
	Phone       *string   `json:"phone,omitempty"`
	Email       *string   `json:"email,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRestaurantRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email" binding:"omitempty,email"`
	ImageURL    *string `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateRestaurantRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email" binding:"omitempty,email"`
	ImageURL    *string `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}
