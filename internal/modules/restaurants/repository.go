package restaurants

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(restaurant *Restaurant) error {
	query := `
		INSERT INTO restaurants (name, description, address, phone, email, image_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	var isActive bool = true
	if restaurant.IsActive == false {
		isActive = false
	}

	err := r.db.QueryRow(
		query,
		restaurant.Name,
		restaurant.Description,
		restaurant.Address,
		restaurant.Phone,
		restaurant.Email,
		restaurant.ImageURL,
		isActive,
	).Scan(&restaurant.ID, &restaurant.CreatedAt, &restaurant.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create restaurant: %w", err)
	}

	restaurant.IsActive = isActive
	return nil
}

func (r *Repository) GetByID(id int64) (*Restaurant, error) {
	query := `
		SELECT id, name, description, address, phone, email, image_url, is_active, created_at, updated_at
		FROM restaurants
		WHERE id = $1
	`

	restaurant := &Restaurant{}
	var description, address, phone, email, imageURL sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&description,
		&address,
		&phone,
		&email,
		&imageURL,
		&restaurant.IsActive,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("restaurant not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant: %w", err)
	}

	if description.Valid {
		restaurant.Description = &description.String
	}
	if address.Valid {
		restaurant.Address = &address.String
	}
	if phone.Valid {
		restaurant.Phone = &phone.String
	}
	if email.Valid {
		restaurant.Email = &email.String
	}
	if imageURL.Valid {
		restaurant.ImageURL = &imageURL.String
	}

	return restaurant, nil
}

func (r *Repository) GetAll(limit, offset int) ([]*Restaurant, int, error) {
	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM restaurants").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count restaurants: %w", err)
	}

	// Get restaurants
	query := `
		SELECT id, name, description, address, phone, email, image_url, is_active, created_at, updated_at
		FROM restaurants
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get restaurants: %w", err)
	}
	defer rows.Close()

	var restaurants []*Restaurant
	for rows.Next() {
		restaurant := &Restaurant{}
		var description, address, phone, email, imageURL sql.NullString

		err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&description,
			&address,
			&phone,
			&email,
			&imageURL,
			&restaurant.IsActive,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan restaurant: %w", err)
		}

		if description.Valid {
			restaurant.Description = &description.String
		}
		if address.Valid {
			restaurant.Address = &address.String
		}
		if phone.Valid {
			restaurant.Phone = &phone.String
		}
		if email.Valid {
			restaurant.Email = &email.String
		}
		if imageURL.Valid {
			restaurant.ImageURL = &imageURL.String
		}

		restaurants = append(restaurants, restaurant)
	}

	return restaurants, total, nil
}

func (r *Repository) Update(id int64, restaurant *Restaurant) error {
	query := `
		UPDATE restaurants
		SET name = $1,
			description = $2,
			address = $3,
			phone = $4,
			email = $5,
			image_url = $6,
			is_active = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		restaurant.Name,
		restaurant.Description,
		restaurant.Address,
		restaurant.Phone,
		restaurant.Email,
		restaurant.ImageURL,
		restaurant.IsActive,
		id,
	).Scan(&restaurant.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("restaurant not found")
	}
	if err != nil {
		return fmt.Errorf("failed to update restaurant: %w", err)
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM restaurants WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete restaurant: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("restaurant not found")
	}

	return nil
}
