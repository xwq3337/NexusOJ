package migrations

import (
	"fmt"
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
		&models.ContestProblem{},
		&models.ContestParticipant{},
		&models.ContestRankItem{},
		&models.ContestRecord{},
		&models.FriendShips{},
		&models.FriendShipRequest{},
	)

	// 回填已有 contest_problem 的题目数据（从 problem 表拷贝）
	backfillContestProblems()

	// 确保 title 列为 varchar(255)，以便创建全文索引
	titleColumns := []struct {
		table  string
		column string
	}{
		{"contest", "title"},
		{"problem", "title"},
		{"blog", "title"},
		{"training", "title"},
	}
	for _, col := range titleColumns {
		dao.MysqlClient.Exec(fmt.Sprintf(
			"ALTER TABLE %s MODIFY COLUMN %s VARCHAR(255)", col.table, col.column,
		))
	}

	// 创建全文索引（如果不存在）
	type fulltextIndex struct {
		table  string
		column string
		name   string
	}
	indexes := []fulltextIndex{
		{"contest", "title", "ft_contest_title"},
		{"problem", "title", "ft_problem_title"},
		{"blog", "title", "ft_blog_title"},
		{"training", "title", "ft_training_title"},
		{"user", "username", "ft_user_username"},
		{"user", "nickname", "ft_user_nickname"},
	}
	for _, idx := range indexes {
		// 先检查索引是否已存在
		var count int64
		dao.MysqlClient.Raw(
			"SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?",
			idx.table, idx.name,
		).Scan(&count)
		if count == 0 {
			sql := fmt.Sprintf("ALTER TABLE %s ADD FULLTEXT INDEX %s(%s)", idx.table, idx.name, idx.column)
			dao.MysqlClient.Exec(sql)
		}
	}
}

// backfillContestProblems 将已有 contest_problem 的 problem_id 对应的题目数据拷贝到新字段
// 仅在 title 列为空的记录上执行（幂等）
func backfillContestProblems() {
	// 检查是否还有旧的 problem_id 列
	var colCount int64
	dao.MysqlClient.Raw(
		"SELECT COUNT(*) FROM information_schema.columns WHERE table_schema = DATABASE() AND table_name = 'contest_problem' AND column_name = 'problem_id'",
	).Scan(&colCount)
	if colCount == 0 {
		return // 旧列已不存在，无需回填
	}

	// 检查是否有需要回填的记录
	var needBackfill int64
	dao.MysqlClient.Raw(
		"SELECT COUNT(*) FROM contest_problem WHERE title = '' OR title IS NULL",
	).Scan(&needBackfill)
	if needBackfill == 0 {
		return
	}

	dao.MysqlClient.Exec(`UPDATE contest_problem cp
		INNER JOIN problem p ON p.id = cp.problem_id
		SET
			cp.title = p.title,
			cp.context = p.context,
			cp.input_description = p.input_description,
			cp.output_description = p.output_description,
			cp.tips = p.tips,
			cp.difficulty = p.difficulty,
			cp.judge_case = p.judge_case,
			cp.judge_config = p.judge_config,
			cp.judge_sample = p.judge_sample,
			cp.tags = p.tags,
			cp.submission = COALESCE(p.submission, 0),
			cp.accept = COALESCE(p.accept, 0),
			cp.source_problem_id = CAST(cp.problem_id AS SIGNED)
		WHERE cp.title = '' OR cp.title IS NULL`)
}
