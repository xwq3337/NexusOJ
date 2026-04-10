<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Download, Search } from '@element-plus/icons-vue'
import { logApi } from '@nexusoj/server'
import { type Log } from '@nexusoj/type'
import dayjs from 'dayjs'

// 日期列表
const dateList = ref<string[]>([])
const selectedDate = ref<string>('')
const loading = ref(false)

// 日志数据
const logData = ref<Log[]>([])
const logLoading = ref(false)

// 筛选
const levelFilter = ref<string>('')
const ipFilter = ref<string>('')

// 获取日期列表
const fetchDateList = async () => {
  loading.value = true
  try {
    const res = await logApi.getDate()
    const { code, info } = res
    if (code === 200 && info) {
      const { dates } = info
      dateList.value = dates
      const today = dayjs().format('YYYY-MM-DD')
      if (dates.includes(today)) {
        selectedDate.value = today
      } else if (dates.length > 0) {
        selectedDate.value = dates[0]
      }
    } else {
      ElMessage.error('获取日期列表失败')
    }
  } catch (e) {
    ElMessage.error('获取日期列表失败')
    console.error(e)
  } finally {
    loading.value = false
  }
}

// 获取日志内容
const fetchLogContent = async (date: string) => {
  if (!date) return

  logLoading.value = true
  try {
    // API 直接返回纯文本内容
    const content = await logApi.getLogList(date)

    // 解析日志内容
    const parsedLogs = parseLogContent(content)
    logData.value = parsedLogs
  } catch (e) {
    ElMessage.error('获取日志内容失败')
    console.error(e)
    logData.value = []
  } finally {
    logLoading.value = false
  }
}

// 解析日志内容
const parseLogContent = (content: string): Log[] => {
  if (!content) return []

  try {
    // 尝试按行解析 JSON
    const lines = content.split('\n').filter(line => line.trim())
    const logs: Log[] = []

    for (const line of lines) {
      try {
        const log = JSON.parse(line) as Log
        logs.push(log)
      } catch {
        // 如果解析失败，创建一个基本的日志对象
        logs.push({
          level: 'INFO',
          ts: dayjs().format(),
          caller: 'unknown',
          msg: line,
          status: 0,
          method: '',
          path: '',
          ip: '',
          latency: '',
          user_agent: '',
          headers: {},
          response_body: {}
        })
      }
    }

    // 倒序数组，让最新的日志（文件末尾）显示在最上面
    return logs.reverse()
  } catch {
    return []
  }
}

// 过滤后的日志
const filteredLogs = ref<Log[]>([])

// 监听日期变化，重置数据并重新加载
watch(selectedDate, () => {
  if (selectedDate.value) {
    fetchLogContent(selectedDate.value)
  }
}, { immediate: false })

// 监听日志数据变化，应用筛选
watch([logData, levelFilter, ipFilter], () => {
  let result = logData.value

  // 按级别筛选
  if (levelFilter.value) {
    result = result.filter(log => log.level === levelFilter.value)
  }

  // 按 IP 筛选
  if (ipFilter.value) {
    result = result.filter(log => log.ip && log.ip.includes(ipFilter.value))
  }

  filteredLogs.value = result
})

// 刷新日志
const refreshLogs = () => {
  if (selectedDate.value) {
    fetchLogContent(selectedDate.value)
  }
}

// 下载日志
const downloadLogs = () => {
  if (!logData.value || logData.value.length === 0) {
    ElMessage.warning('暂无日志内容')
    return
  }

  const content = logData.value.map(log => JSON.stringify(log, null, 2)).join('\n')
  const blob = new Blob([content], { type: 'application/json;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `app_${selectedDate.value}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  ElMessage.success('下载成功')
}

// 获取日志级别标签类型
const getLevelTagType = (level: string) => {
  switch (level.toUpperCase()) {
    case 'ERROR':
      return 'danger'
    case 'WARN':
      return 'warning'
    case 'INFO':
      return 'info'
    case 'DEBUG':
      return 'success'
    default:
      return ''
  }
}

// 获取状态码标签类型
const getStatusTagType = (status: number) => {
  if (status >= 500) return 'danger'
  if (status >= 400) return 'warning'
  if (status >= 300) return 'info'
  if (status >= 200) return 'success'
  return ''
}

// 格式化时间戳
const formatTimestamp = (ts: string) => {
  return dayjs(ts).format('YYYY-MM-DD HH:mm:ss')
}

// 格式化延迟
const formatLatency = (latency: string) => {
  if (!latency) return '-'
  // 保持原始单位显示
  return latency
}

// 展开/收起详情
const expandedRows = ref<Set<number>>(new Set())

// 格式化 JSON
const formatJson = (data: string | object | null | undefined): string => {
  if (data === null || data === undefined) {
    return ''
  }

  if (typeof data === 'string') {
    try {
      const parsed = JSON.parse(data)
      return JSON.stringify(parsed, null, 2)
    } catch {
      return data
    }
  }

  return JSON.stringify(data, null, 2)
}


onMounted(() => {
  fetchDateList()
})
</script>

<template>
  <div class="log-list-container">
    <!-- 顶部操作栏 -->
    <el-card class="header-card">
      <div class="header-content">
        <div class="date-selector">
          <span class="label" >选择日期：</span>
          <el-select
            v-model="selectedDate"
            placeholder="请选择日期"
            :loading="loading"
            style="width: 200px"
            @change="() => { fetchLogContent(selectedDate) }"
          >
            <el-option
              v-for="date in dateList"
              :key="date"
              :label="date"
              :value="date"
            />
          </el-select>
        </div>

        <div class="actions">
          <el-input
            v-model="ipFilter"
            placeholder="筛选 IP"
            clearable
            style="width: 150px"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>

          <el-select
            v-model="levelFilter"
            placeholder="日志级别"
            clearable
            style="width: 120px"
          >
            <el-option label="ERROR" value="ERROR" />
            <el-option label="WARN" value="WARN" />
            <el-option label="INFO" value="INFO" />
            <el-option label="DEBUG" value="DEBUG" />
          </el-select>

          <el-button type="primary" :icon="Refresh" @click="refreshLogs" :loading="logLoading">
            刷新
          </el-button>
          <el-button type="success" :icon="Download" @click="downloadLogs">
            下载
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 日志表格 -->
    <el-card class="log-table-card">
      <template #header>
        <div class="card-header">
          <span>日志列表 - {{ selectedDate || '未选择日期' }}</span>
          <el-tag v-if="selectedDate" type="info">
            共 {{ filteredLogs.length }} 条
          </el-tag>
        </div>
      </template>

      <el-table
        v-loading="logLoading"
        :data="filteredLogs"
        stripe
        style="width: 100%"
        :expand-row-keys="Array.from(expandedRows)"
        row-key="ts"
      >
        <el-table-column type="expand">
          <template #default="{ row }">
            <div class="log-detail">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="时间戳">
                  {{ formatTimestamp(row.ts) }}
                </el-descriptions-item>
                <el-descriptions-item label="调用者">
                  {{ row.caller || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="客户端 IP">
                  {{ row.ip || '-' }}
                </el-descriptions-item>
                <el-descriptions-item label="延迟">
                  {{ formatLatency(row.latency) }}
                </el-descriptions-item>
                <el-descriptions-item label="ID">
                  <div class="user-id">{{ row.user_id || '-' }}</div>
                </el-descriptions-item>
                <el-descriptions-item label="参数">
                  <div class="user-id">{{ row.parameters || '-' }}</div>
                </el-descriptions-item>
                <el-descriptions-item label="UA" :span="2">
                  <div class="user-agent">{{ row.user_agent || '-' }}</div>
                </el-descriptions-item>
                <el-descriptions-item label="消息" :span="2">
                  <pre class="log-message">{{ row.msg }}</pre>
                </el-descriptions-item>
                <el-descriptions-item v-if="row.headers" label="Headers" :span="2">
                  <pre class="json-content">{{ JSON.stringify(row.headers, null, 2) }}</pre>
                </el-descriptions-item>
                <el-descriptions-item v-if="row.response_body" label="Response" :span="2">
                  <pre class="json-content">{{ formatJson(row.response_body) }}</pre>
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="level" label="级别" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getLevelTagType(row.level)" size="small">
              {{ row.level }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="ts" label="时间" width="170">
          <template #default="{ row }">
            {{ formatTimestamp(row.ts) }}
          </template>
        </el-table-column>

        <el-table-column prop="method" label="方法" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.method" :type="row.method === 'GET' ? 'info' : 'warning'" size="small">
              {{ row.method }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="path" label="路径" width="150" show-overflow-tooltip />

        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.status" :type="getStatusTagType(row.status)" size="small">
              {{ row.status }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="ip" label="IP" width="140" show-overflow-tooltip />

        <el-table-column prop="latency" label="延迟" width="130">
          <template #default="{ row }">
            {{ formatLatency(row.latency) }}
          </template>
        </el-table-column>

        <el-table-column prop="caller" label="调用者" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.log-list-container {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.header-card {
  flex-shrink: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.date-selector {
  display: flex;
  align-items: center;
  gap: 8px;
}

.label {
  font-weight: 500;
  color: #606266;
}

.actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.log-table-card {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.log-table-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.log-detail {
  padding: 20px;
  background-color: #f5f7fa;
}

.user-agent {
  word-break: break-all;
  font-size: 13px;
  color: #606266;
}
.user-id {
  word-break: break-all;
  font-size: 13px;
  color: #606266;
}
.log-message {
  margin: 0;
  padding: 10px;
  background-color: #fff;
  border-radius: 4px;
  font-size: 13px;
  line-height: 1.6;
  max-height: 200px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

.json-content {
  margin: 0;
  padding: 10px;
  background-color: #fff;
  border-radius: 4px;
  font-size: 12px;
  line-height: 1.5;
  max-height: 300px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
  word-wrap: break-word;
}

.pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

/* 响应式 */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: stretch;
  }

  .actions {
    flex-wrap: wrap;
  }
}

/* 滚动条样式 */
.json-content::-webkit-scrollbar,
.log-message::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.json-content::-webkit-scrollbar-track,
.log-message::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.json-content::-webkit-scrollbar-thumb,
.log-message::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.json-content::-webkit-scrollbar-thumb:hover,
.log-message::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
