package services

import (
	"crypto/rand"
	"encoding/base32"
	"go-fiber-template/src/domain/entities"
	"go-fiber-template/src/domain/repositories"
)

type usersService struct {
	UsersRepository repositories.IUsersRepository
}

type IUsersService interface {
	GetAllUser() (*[]entities.UserDataFormat, error)
	InsertNewAccount(data *entities.NewUserBody) error
	UpdateUser(userID string, data *entities.NewUserBody) error
	DeleteUser(userID string) error
	GetUser(userID string) (*entities.UserDataFormat, error)
}

func NewUsersService(repo0 repositories.IUsersRepository) IUsersService {
	return &usersService{
		UsersRepository: repo0,
	}
}

func (sv *usersService) InsertNewAccount(data *entities.NewUserBody) error {
	secret := make([]byte, 4)
	_, err := rand.Reader.Read(secret)
	if err != nil {
		return err
	}
	var b32NoPadding = base32.StdEncoding.WithPadding(base32.NoPadding)
	userID := b32NoPadding.EncodeToString(secret)
	newData := &entities.UserDataFormat{
		UserID:   userID,
		Username: data.Username,
		Email:    data.Email,
	}
	err = sv.UsersRepository.InsertNewUser(newData)
	if err != nil {
		return err
	}
	return nil
}

func (sv *usersService) GetAllUser() (*[]entities.UserDataFormat, error) {
	userData, err := sv.UsersRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return userData, nil

}

func (sv *usersService) UpdateUser(userID string, data *entities.NewUserBody) error {
	if err := sv.UsersRepository.UpdateUser(userID, data); err != nil {
		return err
	}
	return nil
}

func (sv *usersService) DeleteUser(userID string) error {
	if err := sv.UsersRepository.DeleteUser(userID); err != nil {
		return err
	}
	return nil
}

func (sv *usersService) GetUser(userID string) (*entities.UserDataFormat, error) {
	userData, err := sv.UsersRepository.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return userData, nil
}
