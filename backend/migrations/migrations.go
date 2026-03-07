package migrations

import (
	"nexus/dao"
	"nexus/models"
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
		&models.FriendShipRequest{},
	)
}
