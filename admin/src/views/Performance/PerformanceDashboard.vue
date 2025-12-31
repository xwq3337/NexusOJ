<template>
  <div class="performance-dashboard">
    <el-row :gutter="20">
      <!-- Web Vitals 指标 -->
      <el-col :span="24">
        <el-card class="mb-4">
          <template #header>
            <div class="card-header">
              <h3>Core Web Vitals</h3>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="8">
              <div class="metric-card">
                <h4>CLS (累积布局偏移)</h4>
                <div class="metric-value" :class="clsStatus">{{ performanceStore.webVitals.cls.toFixed(3) }}</div>
                <div class="metric-status">{{ getMetricStatus(performanceStore.webVitals.cls, 'CLS') }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="metric-card">
                <h4>FID (首次输入延迟)</h4>
                <div class="metric-value" :class="fidStatus">{{ performanceStore.webVitals.fid.toFixed(2) }}ms</div>
                <div class="metric-status">{{ getMetricStatus(performanceStore.webVitals.fid, 'FID') }}</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="metric-card">
                <h4>LCP (最大内容绘制)</h4>
                <div class="metric-value" :class="lcpStatus">{{ performanceStore.webVitals.lcp.toFixed(2) }}ms</div>
                <div class="metric-status">{{ getMetricStatus(performanceStore.webVitals.lcp, 'LCP') }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <!-- API 请求监控 -->
      <el-col :span="24">
        <el-card class="mb-4">
          <template #header>
            <div class="card-header">
              <h3>API 请求监控</h3>
              <el-select v-model="timeRange" placeholder="选择时间范围" @change="updateApiMetrics">
                <el-option label="最近1小时" value="1h" />
                <el-option label="最近24小时" value="24h" />
                <el-option label="最近7天" value="7d" />
              </el-select>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <div ref="apiRequestsChart" style="height: 300px"></div>
            </el-col>
            <el-col :span="12">
              <div ref="apiLatencyChart" style="height: 300px"></div>
            </el-col>
          </el-row>
          <el-row :gutter="20" class="mt-4">
            <el-col :span="8">
              <el-statistic title="总请求数" :value="performanceStore.apiMetrics.requestCount">
                <template #suffix>
                  <el-tag size="small" type="info">次</el-tag>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic
                title="平均响应时间"
                :value="performanceStore.apiMetrics.averageLatency"
                :precision="2"
              >
                <template #suffix>
                  <el-tag size="small" type="info">ms</el-tag>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic
                title="错误率"
                :value="performanceStore.apiMetrics.errorRate * 100"
                :precision="2"
              >
                <template #suffix>
                  <el-tag size="small" type="info">%</el-tag>
                </template>
              </el-statistic>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <!-- 用户行为分析 -->
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <h3>用户行为分析</h3>
              <el-button type="primary" @click="showUserSession">
                查看用户会话回放
              </el-button>
            </div>
          </template>
          <el-row :gutter="20" class="mb-4">
            <el-col :span="8">
              <el-statistic title="会话数" :value="performanceStore.userBehavior.sessionCount">
                <template #suffix>
                  <el-tag size="small" type="info">次</el-tag>
                </template>
              </el-statistic>
            </el-col>
            <el-col :span="8">
              <el-statistic
                title="平均会话时长"
                :value="formatDuration(performanceStore.userBehavior.averageSessionDuration)"
              />
            </el-col>
            <el-col :span="8">
              <el-statistic
                title="跳出率"
                :value="performanceStore.userBehavior.bounceRate * 100"
                :precision="2"
              >
                <template #suffix>
                  <el-tag size="small" type="info">%</el-tag>
                </template>
              </el-statistic>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户会话回放对话框 -->
    <el-dialog
      v-model="sessionReplayVisible"
      title="用户会话回放"
      width="80%"
      destroy-on-close
    >
      <div ref="rrwebPlayer" style="height: 600px"></div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { usePerformanceStore } from '@/stores/performance'
import RRWebPlayer from 'rrweb-player'
import type { RRwebPlayerOptions } from 'rrweb-player'
import 'rrweb-player/dist/style.css'

const performanceStore = usePerformanceStore()
const timeRange = ref('24h')
const apiRequestsChart = ref<HTMLElement | null>(null)
const apiLatencyChart = ref<HTMLElement | null>(null)
const sessionReplayVisible = ref(false)
const rrwebPlayer = ref<HTMLElement | null>(null)
let requestsChartInstance: echarts.ECharts | null = null
let latencyChartInstance: echarts.ECharts | null = null
let player: RRWebPlayer | null = null

// 计算状态类
const clsStatus = computed(() => getStatusClass(performanceStore.webVitals.cls, 'CLS'))
const fidStatus = computed(() => getStatusClass(performanceStore.webVitals.fid, 'FID'))
const lcpStatus = computed(() => getStatusClass(performanceStore.webVitals.lcp, 'LCP'))

// 性能指标状态判断
function getMetricStatus(value: number, metric: string): string {
  switch (metric) {
    case 'CLS':
      return value <= 0.1 ? '良好' : value <= 0.25 ? '需要改进' : '较差'
    case 'FID':
      return value <= 100 ? '良好' : value <= 300 ? '需要改进' : '较差'
    case 'LCP':
      return value <= 2500 ? '良好' : value <= 4000 ? '需要改进' : '较差'
    default:
      return '未知'
  }
}

function getStatusClass(value: number, metric: string): string {
  const status = getMetricStatus(value, metric)
  return {
    '良好': 'status-good',
    '需要改进': 'status-needs-improvement',
    '较差': 'status-poor'
  }[status] || ''
}

// 格式化持续时间
function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}秒`
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}分${remainingSeconds}秒`
}

// 初始化图表
function initCharts() {
  if (apiRequestsChart.value && apiLatencyChart.value) {
    requestsChartInstance = echarts.init(apiRequestsChart.value)
    latencyChartInstance = echarts.init(apiLatencyChart.value)

    // API请求数量图表配置
    requestsChartInstance.setOption({
      title: { text: 'API请求数量' },
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['00:00', '06:00', '12:00', '18:00', '24:00'] },
      yAxis: { type: 'value' },
      series: [{
        data: [150, 230, 224, 218, 135],
        type: 'line',
        smooth: true
      }]
    })

    // API响应时间图表配置
    latencyChartInstance.setOption({
      title: { text: 'API响应时间' },
      tooltip: { trigger: 'axis' },
      xAxis: { type: 'category', data: ['00:00', '06:00', '12:00', '18:00', '24:00'] },
      yAxis: { type: 'value', name: 'ms' },
      series: [{
        data: [120, 200, 150, 80, 70],
        type: 'line',
        smooth: true
      }]
    })
  }
}

// 更新API指标
function updateApiMetrics() {
  // 这里实现根据timeRange更新图表数据的逻辑
  console.log('Updating metrics for timeRange:', timeRange.value)
}

// 显示用户会话回放
function showUserSession() {
  sessionReplayVisible.value = true
  if (performanceStore.rrwebEvents.length > 0) {
    // 等待 DOM 更新后初始化播放器
    setTimeout(() => {
      if (rrwebPlayer.value && !player) {
        const options: RRwebPlayerOptions = {
          target: rrwebPlayer.value,
          props: {
            events: performanceStore.rrwebEvents
          }
        }
        player = new RRWebPlayer(options)
      }
    })
  }
}

// 监听窗口大小变化
function handleResize() {
  requestsChartInstance?.resize()
  latencyChartInstance?.resize()
}

onMounted(() => {
  performanceStore.initializeMonitoring()
  initCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  requestsChartInstance?.dispose()
  latencyChartInstance?.dispose()
  if (player) {
    // RRWebPlayer 实例没有提供 destroy 方法，但我们可以移除 DOM 元素
    if (rrwebPlayer.value) {
      rrwebPlayer.value.innerHTML = ''
    }
    player = null
  }
})
</script>

<style scoped lang="scss">
.performance-dashboard {
  padding: 20px;

  .mb-4 {
    margin-bottom: 20px;
  }

  .mt-4 {
    margin-top: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }

  .metric-card {
    text-align: center;
    padding: 20px;
    border-radius: 8px;
    background-color: var(--el-bg-color);
    box-shadow: var(--el-box-shadow-lighter);
    transition: all 0.3s ease;
    h4 {
      margin: 0 0 10px;
      color: var(--el-text-color-regular);
      font-size: 14px;
    }

    .metric-value {
      font-size: 24px;
      font-weight: bold;
      margin: 10px 0;

      &.status-good {
        color: var(--el-color-success);
      }

      &.status-needs-improvement {
        color: var(--el-color-warning);
      }

      &.status-poor {
        color: var(--el-color-danger);
      }
    }

    .metric-status {
      font-size: 14px;
      color: var(--el-text-color-secondary);
    }
  }
}
</style>
