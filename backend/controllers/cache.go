package controllers

import (
	"context"
	"fmt"
	"nexus/dao"
	"nexus/utils"

	"github.com/gin-gonic/gin"
)

// 永久 key
var (
	ProblemStatusSolvedBit    = "user:%d:problem_status_solved_bit"
	ProblemStatusAttemptedBit = "user:%d:problem_status_attempted_bit"
)

// 半永久key
var (
	ContestPassword       = "contest:%d:password"
	ContestParticipants   = "contest:%s:participants" // set
	ContestRanking        = "contest:%s:ranking"
	ContestRankDetail     = "contest:%s:ranking:detail"
	ContestInfo           = "contest:%s:info"
	ContestParticipantBit = "contest:%s:participants:bit"  // bitmap: key=contestID
	ContestProblemStatus  = "contest:%s:problem_status:%d" // hash: contestID:userID
)
var ctx = context.Background()

type RedisCache struct{}

func (RedisCache) RefreshRedisCache(c *gin.Context) {
	var userIds []int
	var problemIds []int
	var contestIds []int
	err := dao.MysqlClient.Raw("SELECT id FROM user").Scan(&userIds).Error
	if err != nil {
		utils.ReturnError(c, 500, "所有用户 id 拉取失败")
	}
	err = dao.MysqlClient.Raw("SELECT id FROM problem").Scan(&problemIds).Error
	if err != nil {
		utils.ReturnError(c, 500, "所有题目 id 拉取失败")
	}
	err = dao.MysqlClient.Raw("SELECT id FROM contest").Scan(&contestIds).Error
	if err != nil {
		utils.ReturnError(c, 500, "所有比赛 id 拉取失败")
	}
	RefreshUserProblemStatus(userIds)
	RefreshContestPassword()
	utils.ReturnSuccess(c, 200, "完成 Redis 刷新", nil)
}

func RefreshUserProblemStatus(userIds []int) {
	for _, userId := range userIds {
		var results []struct {
			ProblemId int    `json:"problem_id"`
			Verdict   string `json:"verdict"`
		}
		dao.MysqlClient.Raw("SELECT problem_id, verdict FROM record WHERE user_id = ?", userId).Scan(&results)
		Map := make(map[int]string)
		for _, result := range results {
			if Map[result.ProblemId] == "" {
				Map[result.ProblemId] = result.Verdict
			}
			if result.Verdict == "Accepted" {
				Map[result.ProblemId] = "Accepted"
			}
		}
		for problemId, verdict := range Map {
			problemIndex := int64(problemId) - 1000
			if verdict == "Accepted" {
				dao.RedisClient.SetBit(ctx, fmt.Sprintf(ProblemStatusSolvedBit, userId), problemIndex, 1)
			}
			dao.RedisClient.SetBit(ctx, fmt.Sprintf(ProblemStatusAttemptedBit, userId), problemIndex, 1)
		}
	}
}

func RefreshContestPassword() {
	var results []struct {
		Id       int    `json:"id"`
		Password string `json:"password"`
	}
	dao.MysqlClient.Raw("SELECT id, password FROM contest").Scan(&results)
	for _, result := range results {
		dao.RedisClient.Set(ctx, fmt.Sprintf(ContestPassword, result.Id), result.Password, 0)
	}
}

var globalMapper = dao.NewMappedUserIDstrategy("global")
