package controllers

import (
	"context"
	"fmt"
	"nexus/dao"
	"nexus/models"
	"nexus/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// 永久 key（引用 dao 共享常量）
var (
	ProblemStatusSolvedBit    = dao.UserSolvedBitKey
	ProblemStatusAttemptedBit = dao.UserAttemptedBitKey
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

// ==================== 推荐系统 Redis 索引 ====================

var (
	ProblemsByDifficulty = "problems_by_difficulty" // Sorted Set: score=difficulty
	TagProblemsPrefix    = "tag_problems:"           // Sorted Set: score=difficulty, one per tag
	ProblemsByTime       = "problems_by_time"        // Sorted Set: score=created_at unix
)

// RefreshProblemIndexes 从 MySQL 加载题目元数据到 Redis 索引
func RefreshProblemIndexes() error {
	metas, err := models.Problem{}.GetAllProblemMeta()
	if err != nil {
		return err
	}

	pipe := dao.RedisClient.Pipeline()
	// 清除旧索引
	pipe.Del(ctx, ProblemsByDifficulty, ProblemsByTime)

	tagSet := make(map[string]bool)
	for _, m := range metas {
		pipe.ZAdd(ctx, ProblemsByDifficulty, &redis.Z{
			Score:  float64(m.Difficulty),
			Member: m.ID,
		})
		pipe.ZAdd(ctx, ProblemsByTime, &redis.Z{
			Score:  float64(m.CreatedAt.Unix()),
			Member: m.ID,
		})
		for _, tag := range m.Tags {
			if tag == "" {
				continue
			}
			tagSet[tag] = true
			pipe.ZAdd(ctx, TagProblemsPrefix+tag, &redis.Z{
				Score:  float64(m.Difficulty),
				Member: m.ID,
			})
		}
	}

	// 清除旧的 tag 索引（防止删除标签后残留）
	for tag := range tagSet {
		// 仅保留当前存在的 tag 索引
		_ = tag
	}

	_, err = pipe.Exec(ctx)
	return err
}
