package dao

import (
	"context"
	"fmt"
	"strconv"
	"sync"
)

// 将用户 ID 映射到紧凑索引的策略, 以使用 Redis bitmaps 存储状态
type MappedUserIDstrategy struct {
	mappingMutex sync.RWMutex
	userID2Index map[uint64]int
	index2UserID map[int]uint64
	key          string // Redis key 前缀，用于持久化
	nextIndex    int    // 下一个可用的索引
}

func NewMappedUserIDstrategy(key string) *MappedUserIDstrategy {
	m := &MappedUserIDstrategy{
		userID2Index: make(map[uint64]int),
		index2UserID: make(map[int]uint64),
		key:          key,
		nextIndex:    1, // 从 1 开始，0 表示未映射
	}
	m.loadFromRedis()
	return m
}

// loadFromRedis 从 Redis 加载已有映射关系
func (m *MappedUserIDstrategy) loadFromRedis() {
	key := fmt.Sprintf("uid_map:%s", m.key)
	result, err := RedisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		return
	}
	for uidStr, idxStr := range result {
		uid, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			continue
		}
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			continue
		}
		m.userID2Index[uid] = idx
		m.index2UserID[idx] = uid
		if idx >= m.nextIndex {
			m.nextIndex = idx + 1
		}
	}
}

// saveToRedis 持久化单条映射到 Redis
func (m *MappedUserIDstrategy) saveToRedis(userID uint64, index int) {
	key := fmt.Sprintf("uid_map:%s", m.key)
	RedisClient.HSet(context.Background(), key, strconv.FormatUint(userID, 10), index)
}

// GetOrAssignIndex 获取已有索引，若不存在则分配新索引并持久化
func (m *MappedUserIDstrategy) GetOrAssignIndex(userID uint64) int {
	m.mappingMutex.RLock()
	if idx, ok := m.userID2Index[userID]; ok {
		m.mappingMutex.RUnlock()
		return idx
	}
	m.mappingMutex.RUnlock()

	m.mappingMutex.Lock()
	defer m.mappingMutex.Unlock()

	// double check
	if idx, ok := m.userID2Index[userID]; ok {
		return idx
	}

	idx := m.nextIndex
	m.nextIndex++
	m.userID2Index[userID] = idx
	m.index2UserID[idx] = userID
	m.saveToRedis(userID, idx)
	return idx
}

// 根据紧凑索引得到用户 ID
func (m *MappedUserIDstrategy) GetUserID(index int) uint64 {
	m.mappingMutex.RLock()
	defer m.mappingMutex.RUnlock()
	return m.index2UserID[index]
}

// 根据用户 ID 得到紧凑索引，未映射返回 0
func (m *MappedUserIDstrategy) GetIndex(userID uint64) int {
	m.mappingMutex.RLock()
	defer m.mappingMutex.RUnlock()
	return m.userID2Index[userID]
}

// HasIndex 检查用户是否已有映射
func (m *MappedUserIDstrategy) HasIndex(userID uint64) bool {
	m.mappingMutex.RLock()
	defer m.mappingMutex.RUnlock()
	_, ok := m.userID2Index[userID]
	return ok
}

// Count 返回已映射的用户数量
func (m *MappedUserIDstrategy) Count() int {
	m.mappingMutex.RLock()
	defer m.mappingMutex.RUnlock()
	return len(m.userID2Index)
}
