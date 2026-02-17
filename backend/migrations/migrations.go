package migrations

import (
	"pkg/dao"
	"pkg/models"
)

func Migrate() {
	dao.MysqlClient.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.Record{},
		&models.Blog{},
		&models.Training{},
		&models.Contest{},
		&models.FriendShips{},
	)
}
