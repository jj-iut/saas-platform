package restaurants

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *CreateRestaurantRequest) (*Restaurant, error) {
	restaurant := &Restaurant{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}

	if req.IsActive != nil {
		restaurant.IsActive = *req.IsActive
	}

	if err := s.repo.Create(restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (s *Service) GetByID(id int64) (*Restaurant, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetAll(page, pageSize int) ([]*Restaurant, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	return s.repo.GetAll(pageSize, offset)
}

func (s *Service) Update(id int64, req *UpdateRestaurantRequest) (*Restaurant, error) {
	restaurant, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		restaurant.Name = *req.Name
	}
	if req.Description != nil {
		restaurant.Description = req.Description
	}
	if req.Address != nil {
		restaurant.Address = req.Address
	}
	if req.Phone != nil {
		restaurant.Phone = req.Phone
	}
	if req.Email != nil {
		restaurant.Email = req.Email
	}
	if req.ImageURL != nil {
		restaurant.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		restaurant.IsActive = *req.IsActive
	}

	if err := s.repo.Update(id, restaurant); err != nil {
		return nil, err
	}

	return restaurant, nil
}

func (s *Service) Delete(id int64) error {
	return s.repo.Delete(id)
}
