package services

import (
	"log"
	"nexus/models"
	"time"
)

// InitContestStatusWorker 初始化比赛状态自动检查协程
func InitContestStatusWorker() {
	go contestStatusWorker()
}

func contestStatusWorker() {
	log.Println("Contest status worker started")
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		checkAllContestStatus()
	}
}

func checkAllContestStatus() {
	// 检查所有未结束的比赛
	statuses := []models.ContestStatus{models.ContestStatusUpcoming, models.ContestStatusLive}
	for _, status := range statuses {
		contests, err := models.Contest{}.GetContestsByStatus(status)
		if err != nil {
			continue
		}

		for _, contest := range contests {
			newStatus := CheckContestStatus(&contest)
			if newStatus != contest.Status {
				log.Printf("Contest %s status changed: %s -> %s", contest.ID, contest.Status, newStatus)

				_ = models.Contest{}.UpdateContestStatus(contest.ID, newStatus)
				InvalidateContestInfo(contest.ID)

				// 比赛结束时生成报告
				if newStatus == models.ContestStatusEnded {
					go func(c models.Contest) {
						log.Printf("Generating report for contest %s", c.ID)
						if err := GenerateContestReport(c.ID); err != nil {
							log.Printf("Failed to generate report for contest %s: %v", c.ID, err)
						}
					}(contest)
				}

				// 比赛开始时初始化缓存
				if newStatus == models.ContestStatusLive {
					go InitContestCache(&contest)
				}
			}
		}
	}
}
