<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { contestApi } from '@nexusoj/server'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import type { Contest, ContestRankItem } from '@nexusoj/type'

const route = useRoute()
const contestId = route.params.id as string

const contest = ref<Contest | null>(null)
const problems = ref<any[]>([])
const participantCount = ref(0)
const rankings = ref<ContestRankItem[]>([])
const activeTab = ref('info')

// 导入相关
const importPreview = ref<any[]>([])
const selectedImportIds = ref<number[]>([])
const importing = ref(false)

const statusMap: Record<string, { label: string; type: string }> = {
  Upcoming: { label: '未开始', type: 'info' },
  Live: { label: '进行中', type: 'success' },
  Ended: { label: '已结束', type: 'warning' },
}

const typeMap: Record<string, string> = { ACM: 'ACM 赛制', OI: 'OI 赛制' }

const isContestEnded = () => contest.value?.status === 'Ended'

const fetchDetail = async () => {
  try {
    const { code, info } =  await contestApi.getAdminContestDetail(contestId)
    if (code === 200 && info) {
      contest.value = info.contest
      problems.value = info.problems || []
      participantCount.value = info.contest.participants|| 0
    }
  } catch (e) {
    console.error(e)
  }
}

const fetchRanking = async () => {
  try {
    const res = await contestApi.getContestRanking(contestId)
    if (res.code === 200 && res.info) {
      rankings.value = res.info
    }
  } catch (e) {
    console.error(e)
  }
}

const fetchImportPreview = async () => {
  try {
    const res = await contestApi.getImportPreview(contestId)
    if (res.code === 200 && res.info) {
      importPreview.value = res.info
      selectedImportIds.value = res.info
        .filter((p: any) => !p.imported && !p.source_problem_id)
        .map((p: any) => p.id)
    }
  } catch (e) {
    console.error(e)
  }
}

const handleImport = async () => {
  if (selectedImportIds.value.length === 0) {
    return ElMessage.warning('请选择要导入的题目')
  }
  importing.value = true
  try {
    const res = await contestApi.importProblems(contestId, selectedImportIds.value)
    if (res.code === 200) {
      ElMessage.success('导入成功')
      fetchImportPreview()
    } else {
      ElMessage.error(res.msg || '导入失败')
    }
  } catch (e) {
    ElMessage.error('导入失败' + e)
  }
  importing.value = false
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const handleTabChange = (tab: string) => {
  if (tab === 'ranking' && rankings.value.length === 0) {
    fetchRanking()
  }
  if (tab === 'import' && importPreview.value.length === 0) {
    fetchImportPreview()
  }
}

// ====== 密码管理 ======
const passwordDialog = reactive({ visible: false, loading: false })
const passwordForm = reactive({ newPassword: '', confirmPassword: '' })
const passwordFormRef = ref<FormInstance>()

const passwordRules = reactive<FormRules>({
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur',
    },
  ],
})

const openPasswordDialog = () => {
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  passwordDialog.visible = true
}

const submitPasswordChange = async () => {
  if (!passwordFormRef.value) return
  try {
    await passwordFormRef.value.validate()
    passwordDialog.loading = true
    const res = await contestApi.updateContest({ id: contestId, password: passwordForm.newPassword } as any)
    if (res.code === 200) {
      ElMessage.success('密码修改成功')
      passwordDialog.visible = false
    } else {
      ElMessage.error(res.msg || '密码修改失败')
    }
  } catch (e) {
    console.error(e)
  } finally {
    passwordDialog.loading = false
  }
}

const handlePrivateChange = async (val: boolean) => {
  try {
    const res = await contestApi.updateContest({
      id: contestId,
      is_private: val,
      password: val ? undefined : '',
    } as any)
    if (res.code === 200) {
      ElMessage.success(val ? '已设为私密比赛' : '已设为公开比赛')
    } else {
      // 恢复开关状态
      if (contest.value) contest.value.is_private = !val
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (e) {
    if (contest.value) contest.value.is_private = !val
    ElMessage.error('操作失败')
  }
}

onMounted(fetchDetail)
</script>

<template>
  <div v-if="contest">
    <h2 style="margin-bottom: 16px">{{ contest.title }}</h2>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="基本信息" name="info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ contest.id }}</el-descriptions-item>
          <el-descriptions-item label="标题">{{ contest.title }}</el-descriptions-item>
          <el-descriptions-item label="赛制">{{ typeMap[contest.contest_type] }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusMap[contest.status]?.type" size="small">
              {{ statusMap[contest.status]?.label }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="开始时间">{{ formatTime(contest.begin_at) }}</el-descriptions-item>
          <el-descriptions-item label="结束时间">{{ formatTime(contest.end_at) }}</el-descriptions-item>
          <el-descriptions-item label="时长">{{ contest.duration }} 分钟</el-descriptions-item>
          <el-descriptions-item label="参赛人数">{{ participantCount }}</el-descriptions-item>
          <el-descriptions-item label="封榜">{{ contest.seal_rank ? '是' : '否' }}</el-descriptions-item>
          <el-descriptions-item label="提交/通过">{{ contest.submission }} / {{ contest.accept }}</el-descriptions-item>
          <el-descriptions-item label="介绍" :span="2">{{ contest.introduction || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-card style="margin-top: 20px">
          <template #header>
            <span>访问控制</span>
          </template>
          <div style="display: flex; justify-content: space-between; align-items: center">
            <div>
              <span style="font-weight: 500">私密比赛</span>
              <span style="color: #999; font-size: 12px; margin-left: 8px">
                开启后参赛者需要输入密码才能报名
              </span>
            </div>
            <el-switch v-model="contest.is_private" @change="handlePrivateChange" :disabled="isContestEnded()" />
          </div>
          <div v-if="contest.is_private" style="margin-top: 16px">
            <el-button type="primary" size="small" @click="openPasswordDialog" :disabled="isContestEnded()">
              修改密码
            </el-button>
            <span style="color: #999; font-size: 12px; margin-left: 8px">密码不会在此显示</span>
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="题目列表" name="problems">
        <el-table :data="problems" stripe>
          <el-table-column prop="label" label="标签" width="80" />
          <el-table-column prop="title" label="标题" />
          <el-table-column prop="score" label="分值" width="80" />
          <el-table-column prop="difficulty" label="难度" width="80" />
          <el-table-column label="来源" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.source_problem_id" type="info" size="small">题库 #{{ row.source_problem_id }}</el-tag>
              <span v-else style="color: #999">独立题目</span>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="排名" name="ranking">
        <el-table :data="rankings" stripe>
          <el-table-column prop="rank" label="排名" width="80" />
          <el-table-column prop="username" label="用户" width="150" />
          <el-table-column prop="solved" label="通过数" width="80" />
          <el-table-column label="罚时">
            <template #default="{ row }">
              {{ Math.floor(row.total_penalty / 3600) }}:{{ String(Math.floor((row.total_penalty % 3600) / 60)).padStart(2, '0') }}:{{ String(row.total_penalty % 60).padStart(2, '0') }}
            </template>
          </el-table-column>
          <el-table-column prop="score" label="得分" width="100" />
          <el-table-column label="各题详情">
            <template #default="{ row }">
              <div style="display: flex; gap: 8px; flex-wrap: wrap">
                <el-tag v-for="(detail, label) in row.problems" :key="label" size="small"
                  :type="detail.accepted ? 'success' : 'info'">
                  {{ label }}: {{ detail.accepted ? 'AC' : `${detail.attempts}次` }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane v-if="isContestEnded()" label="导入题库" name="import">
        <el-alert type="info" :closable="false" style="margin-bottom: 16px">
          比赛已结束，可以将比赛题目导入到题库中。已导入的题目会显示题库ID。
        </el-alert>
        <el-table :data="importPreview" stripe @selection-change="(rows: any[]) => selectedImportIds = rows.map(r => r.id)">
          <el-table-column type="selection" width="55" :selectable="(row: any) => !row.imported" />
          <el-table-column prop="label" label="标签" width="80" />
          <el-table-column prop="title" label="标题" />
          <el-table-column prop="score" label="分值" width="80" />
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.imported" type="success" size="small">已导入 #{{ row.source_problem_id }}</el-tag>
              <el-tag v-else type="info" size="small">未导入</el-tag>
            </template>
          </el-table-column>
        </el-table>
        <div style="margin-top: 16px">
          <el-button type="primary" :loading="importing" @click="handleImport">
            导入选中题目 ({{ selectedImportIds.length }})
          </el-button>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="passwordDialog.visible" title="修改比赛密码" width="500px">
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px">
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password placeholder="设置新密码" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password placeholder="再次输入新密码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="submitPasswordChange" :loading="passwordDialog.loading">
          确认
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>
