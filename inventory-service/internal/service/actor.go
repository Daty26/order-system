package service

type Actor struct {
	UserID int
	Role   string
}

func NewActor(userID int, role string) Actor {
	return Actor{
		UserID: userID,
		Role:   role,
	}
}

func (a Actor) IsAdmin() bool {
	return a.Role == "ADMIN"
}
