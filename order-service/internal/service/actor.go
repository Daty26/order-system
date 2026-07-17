package service

type Actor struct {
	UserID int
	Role   string
}

func NewActor(userId int, role string) Actor {
	return Actor{
		UserID: userId,
		Role:   role,
	}
}
func (a Actor) IsAdmin() bool {
	return a.Role == "ADMIN"
}
