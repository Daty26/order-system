package service

type LoginUserInput struct {
	Identifier string
	Password   string
}
type CreateUserInput struct{
	Username string
	Email string
	Password string 
}