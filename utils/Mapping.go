package utils

import (
	"my-gin-app/dto"
	"my-gin-app/models"
)

func MapUserToUserDTO(user models.User) dto.UserResponse {
	return dto.UserResponse{
		UserID:      user.UserID,
		UserName:    user.Username,
		Email:       user.UserEmail,
		PhoneNumber: user.UserPhoneNumber,
		CreatedAt:   user.UserCreatedAt,
		UpdatedAt:   user.UserUpdatedAt,
		DeletedAt:   user.UserDeletedAt,
	}
}
