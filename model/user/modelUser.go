package user

import (
	"github.com/emersonary/go-authentication/model/base"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TUserDTO struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Password string    `json:"password"`
}

type TUserAuth struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type TUser struct {
	base.TBaseWithName
	Email    string
	Phone    string
	Password string
}

func (u *TUserDTO) TUserFrom() (*TUser, error) {

	return NewUser(u.Name, u.Email, u.Phone, u.Password)

}

func NewUserDTO(name string, email string, phone string, password string) *TUserDTO {

	return &TUserDTO{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password}

}

func NewUserDTOEmpty() *TUserDTO {

	return &TUserDTO{}

}

func NewUser(name string, email string, phone string, password string) (*TUser, error) {

	var encryptedPassword string

	if password != "" {

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		encryptedPassword = string(hash)

	} else {

		encryptedPassword = ""

	}

	return &TUser{
		TBaseWithName: base.TBaseWithName{
			TBase: base.TBase{
				Id: uuid.Nil,
			},
			Name: name,
		},
		Email:    email,
		Phone:    phone,
		Password: encryptedPassword}, nil
}

func (u *TUser) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func NewUserEmpty() *TUser {
	return &TUser{
		TBaseWithName: base.TBaseWithName{
			TBase: base.TBase{
				Id: uuid.Nil,
			},
		}}
}
