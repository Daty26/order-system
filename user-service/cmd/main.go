package cmd

import (
	"github.com/Daty26/order-system/user-service/internal/db"
	"github.com/Daty26/order-system/user-service/internal/repository"
)

func main() {
	db.DBInit()
	repository.NewPostgresRepository(db.DBConn)
}
