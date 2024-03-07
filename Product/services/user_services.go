package services

import (
	"fmt"

	"github.com/Ferriem/go-web/Product/datamodels"
	"github.com/Ferriem/go-web/Product/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(username string, pwd string) (*datamodels.User, bool)
	AddUser(user *datamodels.User) (int64, error)
}

type UserService struct {
	userRepository repositories.IUser
}

func NewUserService(repository repositories.IUser) IUserService {
	return &UserService{userRepository: repository}
}

func ValidatePassword(userPassword string, hashed string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}

func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

func (u *UserService) IsPwdSuccess(username string, pwd string) (*datamodels.User, bool) {
	user, err := u.userRepository.LogDefinition(username)
	fmt.Println(err)
	if err != nil {
		return &datamodels.User{}, false
	}
	isOk, _ := ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		err = u.userRepository.FailPlusOne(username)
		if err != nil {
			return &datamodels.User{}, false
		}
		return &datamodels.User{}, false
	}
	return user, true
}

func (u *UserService) AddUser(user *datamodels.User) (int64, error) {
	pwdByte, err := GeneratePassword(user.HashPassword)
	if err != nil {
		return 0, err
	}
	user.HashPassword = string(pwdByte)
	return u.userRepository.Insert(user)
}
