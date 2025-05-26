package dto

import "github.com/MingPV/clean-go-template/internal/entities"

// From entity.User to UserResponse
func ToUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}

func ToUserResponseList(users []*entities.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponse(u)
	}
	return responses
}

// From RegisterRequest to entity.User (optional, if want to use in usecase)
func ToUserEntity(req *RegisterRequest) *entities.User {
	return &entities.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
}
