package userstore

import (
	"context"
	"errors"
	models "salon_be/model"
	"salon_be/model/auth/authconst"
	user "salon_be/model/user/usermodel"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateNewUser(ctx context.Context, input *user.CreateUser) error {
	// Check if user already exists
	var existingUser models.User
	if input.AuthType == authconst.AuthTypePassword || input.AuthType == authconst.AuthTypeEmail {
		if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			return gorm.ErrDuplicatedKey
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
	} else if input.AuthType == authconst.AuthTypePhone {
		if err := s.db.Where("phone_number = ?", input.PhoneNumber).First(&existingUser).Error; err == nil {
			return gorm.ErrDuplicatedKey
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	// Create new user
	newUser := &models.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    input.Password,
		Salt:        input.Salt,
	}

	if input.SQLModel != nil && input.Status != "" {
		newUser.Status = input.Status
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return err
	}

	// Create user authentication entry
	auth := models.UserAuth{
		UserID:            newUser.Id,
		AuthType:          input.AuthType,
		AuthProviderID:    input.AuthProviderID,
		AuthProviderToken: input.AuthProviderToken,
	}

	// If it's a local auth type, we need to hash the password
	if input.AuthType == "password" {
		if input.Password == "" {
			return errors.New("password is required for local authentication")
		}
		// TODO: Implement password hashing
		// hashedPassword, err := hashPassword(input.Password)
		// if err != nil {
		//     return err
		// }
		// auth.AuthProviderToken = hashedPassword
	}

	if err := s.db.Create(&auth).Error; err != nil {
		return err
	}

	// Assign default role (assuming 'user' role exists with ID 1)
	if err := s.db.Exec("INSERT INTO user_role (user_id, role_id) VALUES (?, ?)", newUser.Id, 1).Error; err != nil {
		return err
	}

	return nil
}
