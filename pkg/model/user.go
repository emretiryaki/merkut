package model

type LoginUserQuery struct {
	Username   string
	Password   string
	User       *User
	IpAddress  string
}

type User struct {
	Id            int64
	Email         string
	Name          string
	Login         string
	Password      string
	Salt          string
	Rands         string

}