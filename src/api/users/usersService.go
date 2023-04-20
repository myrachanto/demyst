package users

import (
	httperrors "github.com/myrachanto/erroring"
)

var (
	UserService UserServiceInterface = &userService{}
)

type UserServiceInterface interface {
	Create(user *User) (*User, httperrors.HttpErr)
	Login(user *LoginUser) (*Auth, httperrors.HttpErr)
	Logout(token string) (string, httperrors.HttpErr)
	GetOne(code string) (user *User, errors httperrors.HttpErr)
	GetAll(search string) ([]*User, httperrors.HttpErr)
	Update(code string, user *User) (*User, httperrors.HttpErr)
	PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr)
	Delete(code string) (string, httperrors.HttpErr)
	Forgot(token string) (string, httperrors.HttpErr)
	PasswordReset(email, password string) (string, httperrors.HttpErr)
}
type userService struct {
	repo UserrepoInterface
}

func NewUserService(repository UserrepoInterface) UserServiceInterface {
	return &userService{
		repository,
	}
}
func (service *userService) Create(user *User) (*User, httperrors.HttpErr) {
	return service.repo.Create(user)
}

func (service *userService) GetAll(search string) ([]*User, httperrors.HttpErr) {
	return service.repo.GetAll(search)
}

func (service *userService) Login(auser *LoginUser) (*Auth, httperrors.HttpErr) {
	return service.repo.Login(auser)
}

func (service *userService) Logout(token string) (string, httperrors.HttpErr) {
	return service.repo.Logout(token)
}
func (service *userService) Forgot(email string) (string, httperrors.HttpErr) {
	_, _, err := service.repo.Forgot(email)
	// go emails.ToEmailPassword(email, password)
	return "succesfully sent you an email", err
}
func (service *userService) PasswordReset(email, password string) (string, httperrors.HttpErr) {
	email, err := service.repo.PasswordReset(email, password)
	// go emails.ToEmailPassword(email, password)
	return "succesfully sent you an email", err
}
func (service *userService) GetOne(code string) (*User, httperrors.HttpErr) {
	return service.repo.GetOne(code)
}
func (service *userService) PasswordUpdate(oldpassword, email, newpassword string) (string, string, httperrors.HttpErr) {
	return service.repo.PasswordUpdate(oldpassword, email, newpassword)
}

func (service *userService) Update(code string, user *User) (*User, httperrors.HttpErr) {
	return service.repo.Update(code, user)
}

func (service *userService) Delete(id string) (string, httperrors.HttpErr) {
	return service.repo.Delete(id)
}
