package login

import (
	"github.com/emretiryaki/merkut/pkg/model"
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/util"
	"crypto/subtle"
)

var validatePassword= func (providedPassword string, userPassword string, userSalt string) error {

	passwordHashed := util.EncodePassword(providedPassword, userSalt)
	if subtle.ConstantTimeCompare([]byte(passwordHashed), []byte(userPassword)) != 1 {
		return ErrInvalidCredentials
	}

	return nil
}


var loginMerkutDB = func(query *model.LoginUserQuery) error{

	userQuery := model.GetUserByLoginQuery{LoginOrEmail:query.Username}


	if err := bus.Dispatch(&userQuery); err != nil {
		return err
	}

	user := userQuery.Result

	if err := validatePassword(query.Password, user.Password, user.Salt); err != nil {
		return err
	}

	query.User = user
	return nil
}