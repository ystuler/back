package service

import (
	"back/internal/exceptions"
	"back/internal/models"
	"back/internal/repository"
	"back/internal/schemas"
	"back/internal/util"
	"errors"
	"time"
)

type AuthenticationImpl struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthenticationImpl {
	return &AuthenticationImpl{repo: repo}
}

func (s *AuthenticationImpl) SignUp(userSchema *schemas.CreateUserReq) (*schemas.CreateUserResp, error) {
	existingUser, err := s.repo.GetUserByUsername(userSchema.Username)

	if existingUser != nil {
		return nil, errors.New(exceptions.ErrUserAlreadyExists)
	}

	hashedPassword, err := util.HashPassword(userSchema.Password)

	user := &models.User{
		Username:     userSchema.Username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Time{},
	}

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	generatedJWT, err := util.GenerateJWT(createdUser)
	if err != nil {
		return nil, err
	}

	userSchemaResp := schemas.CreateUserResp{
		Token:    generatedJWT,
		ID:       createdUser.ID,
		Username: createdUser.Username,
	}

	return &userSchemaResp, nil
}

func (s *AuthenticationImpl) SignIn(userSchema *schemas.SignInReq) (*schemas.SignInResp, error) {
	existingUser, err := s.repo.GetUserByUsername(userSchema.Username)

	if err != nil {
		return nil, err
	}

	err = util.CheckPassword(userSchema.Password, existingUser.PasswordHash)
	if err != nil {
		return nil, errors.New("password does not match")
	}

	generatedJWT, err := util.GenerateJWT(existingUser)
	if err != nil {
		return nil, err
	}

	userSchemaResp := schemas.SignInResp{
		Token:    generatedJWT,
		ID:       existingUser.ID,
		Username: existingUser.Username,
	}

	return &userSchemaResp, nil
}

func (s *AuthenticationImpl) GetProfile(userID int) (*schemas.GetProfileResp, error) {
	user, err := s.repo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	profile := schemas.Profile{
		ID:       user.ID,
		Username: user.Username,
	}

	return &schemas.GetProfileResp{Profile: profile}, nil
}

func (s *AuthenticationImpl) UpdateUsername(usernameSchema *schemas.UpdateUsernameReq) (*schemas.UpdateUsernameResp, error) {
	user, err := s.repo.GetUserById(usernameSchema.ID)
	if err != nil {
		return nil, err
	}

	user.Username = usernameSchema.Username

	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return &schemas.UpdateUsernameResp{ID: updatedUser.ID, Username: updatedUser.Username}, nil
}

func (s *AuthenticationImpl) UpdatePassword(passwordSchema *schemas.UpdatePasswordReq) error {
	user, err := s.repo.GetUserById(passwordSchema.ID)
	if err != nil {
		return err
	}

	if err := util.CheckPassword(passwordSchema.OldPassword, user.PasswordHash); err != nil {
		return errors.New("old password does not match")
	}

	hashedPassword, err := util.HashPassword(passwordSchema.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	_, err = s.repo.UpdateUser(user)
	return err
}
