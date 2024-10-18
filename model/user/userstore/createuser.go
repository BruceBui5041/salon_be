package userstore

import (
	"context"
	"errors"
	models "salon_be/model"
	user "salon_be/model/user/usermodel"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateNewUser(ctx context.Context, input *user.CreateUser) error {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		return errors.New("user with this email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// Create new user
	newUser := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		Salt:      input.Salt,
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
