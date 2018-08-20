package login

import (
	"errors"
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/model"

	)

var(
	ErrInvalidCredentials =errors.New("Invalid Username or Password")
)

func Init(){
	bus.AddHandler("auth",AuthenticateUser)
}

func AuthenticateUser(query *model.LoginUserQuery) error{

	err := loginMerkutDB(query)
	if err == nil || (err != model.ErrUserNotFound && err != ErrInvalidCredentials) {
		return err
	}

	if err == model.ErrUserNotFound {
		return ErrInvalidCredentials
	}

	return err
}

