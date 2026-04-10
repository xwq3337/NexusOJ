<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { contestApi, problemApi } from '@nexusoj/server'
import { ElMessage } from 'element-plus'
import type { Problem, ContestProblemItem } from '@nexusoj/type'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)

// 表单数据
const form = ref({
  title: '',
  introduction: '',
  contest_type: 'ACM',
  begin_at: '',
  end_at: '',
  duration: 90,
  is_private: false,
  password: '',
  seal_rank: false,
})

const timeRange = ref<[dayjs.Dayjs, dayjs.Dayjs]>()
const onTimeChange = (val: [dayjs.Dayjs, dayjs.Dayjs] | null) => {
  if (val) {
    form.value.begin_at = val[0].toISOString()
    form.value.end_at = val[1].toISOString()
    form.value.duration = val[1].diff(val[0], 'minute')
  }
}

// 题目编辑
interface ProblemEditItem {
  label: string
  score: number
  title: string
  context: string
  input_description: string
  output_description: string
  tips: string
  difficulty: number
  judge_case: { input: string; expected: string }[]
  judge_config: { time_limit: number; memory_limit: number }
  judge_sample: { input: string; expected: string }[]
  tags: string[]
  source_problem_id?: number
}

const problems = ref<ProblemEditItem[]>([])
const addMode = ref<'create' | 'import'>('import')
const showImportDialog = ref(false)

// ====== 从题库导入 ======
const allProblems = ref<Problem[]>([])
const problemSearchKeyword = ref('')
const importingIndex = ref(-1) // 正在导入的题目位置

const fetchProblems = async () => {
  try {
    const res = await problemApi.getProblemList()
    if (res.code === 200 && res.info) {
      allProblems.value = res.info
    }
  } catch (e) {
    console.error(e)
  }
}

const filteredProblems = computed(() => {
  if (!problemSearchKeyword.value) return allProblems.value
  const keyword = problemSearchKeyword.value.toLowerCase()
  return allProblems.value.filter(
    (p) => p.title.toLowerCase().includes(keyword) || String(p.id).includes(keyword)
  )
})

const importProblem = (problem: Problem) => {
  const label = String.fromCharCode(65 + problems.value.length)
  problems.value.push({
    label,
    score: 100,
    title: problem.title,
    context: problem.context || '',
    input_description: problem.input_description || '',
    output_description: problem.output_description || '',
    tips: problem.tips || '',
    difficulty: problem.difficulty,
    judge_case: (problem as any).judge_case || [],
    judge_config: (problem as any).judge_config || { time_limit: 1, memory_limit: 64 },
    judge_sample: (problem as any).judge_sample || [],
    tags: problem.tags || [],
    source_problem_id: problem.id,
  })
  showImportDialog.value = false
}

// ====== 从零创建 ======
const addEmptyProblem = () => {
  const label = String.fromCharCode(65 + problems.value.length)
  problems.value.push({
    label,
    score: 100,
    title: '',
    context: '',
    input_description: '',
    output_description: '',
    tips: '',
    difficulty: 0,
    judge_case: [],
    judge_config: { time_limit: 1, memory_limit: 64 },
    judge_sample: [],
    tags: [],
  })
}

const removeProblem = (index: number) => {
  problems.value.splice(index, 1)
  problems.value.forEach((p, i) => {
    p.label = String.fromCharCode(65 + i)
  })
}

// 新增测试用例
const addTestCase = (item: ProblemEditItem) => {
  item.judge_case.push({ input: '', expected: '' })
}
const removeTestCase = (item: ProblemEditItem, index: number) => {
  item.judge_case.splice(index, 1)
}

// 新增样例
const addSample = (item: ProblemEditItem) => {
  item.judge_sample.push({ input: '', expected: '' })
}
const removeSample = (item: ProblemEditItem, index: number) => {
  item.judge_sample.splice(index, 1)
}

// 提交
const handleSubmit = async () => {
  if (!form.value.title) return ElMessage.warning('请输入标题')
  if (!form.value.begin_at || !form.value.end_at) return ElMessage.warning('请设置比赛时间')

  loading.value = true
  try {
    const res = await contestApi.createContest(form.value)
    if (res.code === 200 && res.info) {
      const contestId = res.info.id
      if (problems.value.length > 0) {
        await contestApi.setContestProblems(contestId, problems.value)
      }
      ElMessage.success('创建成功')
      router.push({ name: 'contest-detail', params: { id: contestId } })
    } else {
      ElMessage.error((res.msg as string) || '创建失败')
    }
  } catch (e: unknown) {
    ElMessage.error('创建失败: ' + (e as Error)?.message)
  }
  loading.value = false
}

fetchProblems()
</script>

<template>
  <div style="max-width: 900px">
    <h2 style="margin-bottom: 20px">创建比赛</h2>

    <el-form label-width="100px">
      <el-form-item label="比赛标题" required>
        <el-input v-model="form.title" placeholder="请输入比赛标题" />
      </el-form-item>

      <el-form-item label="比赛介绍">
        <el-input v-model="form.introduction" type="textarea" :rows="3" placeholder="比赛介绍(可选)" />
      </el-form-item>

      <el-form-item label="赛制" required>
        <el-radio-group v-model="form.contest_type">
          <el-radio value="ACM">ACM 赛制</el-radio>
          <el-radio value="OI">OI 赛制</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="比赛时间" required>
        <el-date-picker
          v-model="timeRange"
          type="datetimerange"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          @change="onTimeChange"
        />
      </el-form-item>

      <el-form-item label="比赛时长">
        <span v-if="form.duration">{{ Math.floor(form.duration / 60) }}小时{{ form.duration % 60 }}分钟</span>
        <span v-else style="color: #999">请先选择比赛时间</span>
      </el-form-item>

      <el-form-item label="私密比赛">
        <el-switch v-model="form.is_private" />
      </el-form-item>

      <el-form-item v-if="form.is_private" label="比赛密码">
        <el-input v-model="form.password" placeholder="设置比赛密码" />
      </el-form-item>

      <el-form-item label="封榜">
        <el-switch v-model="form.seal_rank" />
        <span style="margin-left: 8px; color: #999; font-size: 12px">封榜后用户无法查看实时排名</span>
      </el-form-item>

      <el-divider>比赛题目</el-divider>

      <el-form-item label="添加方式">
        <el-radio-group v-model="addMode">
          <el-radio value="import">从题库导入</el-radio>
          <el-radio value="create">从零创建</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item>
        <div style="width: 100%">
          <div style="display: flex; gap: 8px; margin-bottom: 12px">
            <el-button v-if="addMode === 'import'" @click="showImportDialog = true" type="primary" size="small">从题库导入</el-button>
            <el-button v-if="addMode === 'create'" @click="addEmptyProblem" type="primary" size="small">添加空题目</el-button>
          </div>

          <el-collapse v-if="problems.length > 0">
            <el-collapse-item v-for="(p, index) in problems" :key="index" :name="index">
              <template #title>
                <div style="display: flex; align-items: center; gap: 12px; width: 100%">
                  <el-tag>{{ p.label }}</el-tag>
                  <span>{{ p.title || '(未命名)' }}</span>
                  <el-tag v-if="p.source_problem_id" type="info" size="small">来自题库 #{{ p.source_problem_id }}</el-tag>
                  <el-button type="danger" size="small" style="margin-left: auto" @click.stop="removeProblem(index)">移除</el-button>
                </div>
              </template>

              <el-form label-width="80px" style="padding: 0 20px">
                <el-form-item label="标签">
                  <el-input v-model="p.label" style="width: 80px" />
                </el-form-item>
                <el-form-item label="分值">
                  <el-input-number v-model="p.score" :min="0" :max="1000" />
                </el-form-item>
                <el-form-item label="标题" required>
                  <el-input v-model="p.title" placeholder="题目标题" />
                </el-form-item>
                <el-form-item label="题目描述">
                  <el-input v-model="p.context" type="textarea" :rows="4" placeholder="题目描述" />
                </el-form-item>
                <el-form-item label="输入描述">
                  <el-input v-model="p.input_description" type="textarea" :rows="2" placeholder="输入描述" />
                </el-form-item>
                <el-form-item label="输出描述">
                  <el-input v-model="p.output_description" type="textarea" :rows="2" placeholder="输出描述" />
                </el-form-item>
                <el-form-item label="提示">
                  <el-input v-model="p.tips" type="textarea" :rows="2" placeholder="提示(可选)" />
                </el-form-item>
                <el-form-item label="难度">
                  <el-input-number v-model="p.difficulty" :min="0" :max="10" :step="0.1" />
                </el-form-item>
                <el-form-item label="时间限制">
                  <el-input-number v-model="p.judge_config.time_limit" :min="1" :max="60" /> 秒
                </el-form-item>
                <el-form-item label="内存限制">
                  <el-input-number v-model="p.judge_config.memory_limit" :min="4" :max="1024" /> MB
                </el-form-item>
                <el-form-item label="标签">
                  <el-select v-model="p.tags" multiple filterable allow-create default-first-option placeholder="添加标签">
                    <el-option v-for="tag in p.tags" :key="tag" :label="tag" :value="tag" />
                  </el-select>
                </el-form-item>

                <el-divider>样例</el-divider>
                <div v-for="(s, si) in p.judge_sample" :key="si" style="margin-bottom: 8px; display: flex; gap: 8px">
                  <el-input v-model="s.input" placeholder="输入" type="textarea" :rows="2" style="flex: 1" />
                  <el-input v-model="s.expected" placeholder="输出" type="textarea" :rows="2" style="flex: 1" />
                  <el-button type="danger" size="small" @click="removeSample(p, si)">删除</el-button>
                </div>
                <el-button size="small" @click="addSample(p)">添加样例</el-button>

                <el-divider>测试用例</el-divider>
                <div v-for="(tc, ti) in p.judge_case" :key="ti" style="margin-bottom: 8px; display: flex; gap: 8px">
                  <el-input v-model="tc.input" placeholder="输入" type="textarea" :rows="2" style="flex: 1" />
                  <el-input v-model="tc.expected" placeholder="期望输出" type="textarea" :rows="2" style="flex: 1" />
                  <el-button type="danger" size="small" @click="removeTestCase(p, ti)">删除</el-button>
                </div>
                <el-button size="small" @click="addTestCase(p)">添加测试用例</el-button>
              </el-form>
            </el-collapse-item>
          </el-collapse>
        </div>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :loading="loading">创建比赛</el-button>
        <el-button @click="router.back()">取消</el-button>
      </el-form-item>
    </el-form>

    <!-- 从题库导入对话框 -->
    <el-dialog v-model="showImportDialog" title="从题库导入题目" width="600px">
      <el-input v-model="problemSearchKeyword" placeholder="搜索题目" style="margin-bottom: 10px" />
      <el-table :data="filteredProblems" max-height="400" highlight-current-row
        @row-click="(row: Problem) => importProblem(row)">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="difficulty" label="难度" width="80" />
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button size="small" @click.stop="importProblem(row)">导入</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>
