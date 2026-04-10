<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { contestApi } from '@nexusoj/server'
import type { Contest } from '@nexusoj/type'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const contests = ref<Contest[]>([])
const total = ref(0)
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const search = ref('')

const statusMap: Record<string, { label: string; type: string }> = {
  Upcoming: { label: '未开始', type: 'info' },
  Live: { label: '进行中', type: 'success' },
  Ended: { label: '已结束', type: 'warning' },
}

const typeMap: Record<string, string> = {
  ACM: 'ACM 赛制',
  OI: 'OI 赛制',
}

const fetchContests = async () => {
  loading.value = true
  try {
    const res = await contestApi.getAdminContestList(page.value, pageSize.value, search.value)
    if (res.code === 200 && res.info) {
      contests.value = res.info.list || []
      total.value = res.info.total || 0
    }
  } catch (e) {
    console.error(e)
  }
  loading.value = false
}

const handleDelete = async (id: string) => {
  try {
    await ElMessageBox.confirm('确定删除该比赛？', '警告', { type: 'warning' })
    const res = await contestApi.deleteContest(id)
    if (res.code === 200) {
      ElMessage.success('删除成功')
      fetchContests()
    }
  } catch {}
}

const handlePageChange = (val: number) => {
  page.value = val
  fetchContests()
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

onMounted(fetchContests)
</script>

<template>
  <div>
    <div style="display: flex; justify-content: space-between; margin-bottom: 16px">
      <el-input
        v-model="search"
        placeholder="搜索比赛标题"
        style="width: 300px"
        clearable
        @clear="fetchContests"
        @keyup.enter="fetchContests"
      />
      <el-button type="primary" @click="router.push({ name: 'contest-create' })">创建比赛</el-button>
    </div>

    <el-table :data="contests" v-loading="loading" stripe>
      <el-table-column prop="id" label="ID" width="200" show-overflow-tooltip />
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column label="赛制" width="100">
        <template #default="{ row }">{{ typeMap[row.contest_type] || row.contest_type }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusMap[row.status]?.type" size="small">
            {{ statusMap[row.status]?.label || row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="开始时间" width="180">
        <template #default="{ row }">{{ formatTime(row.begin_at) }}</template>
      </el-table-column>
      <el-table-column label="时长(分)" width="100" prop="duration" />
      <el-table-column label="私密" width="80">
        <template #default="{ row }">
          <el-tag :type="row.is_private ? 'danger' : 'info'" size="small">
            {{ row.is_private ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="提交/通过" width="120">
        <template #default="{ row }">{{ row.submission }} / {{ row.accept }}</template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="router.push({ name: 'contest-detail', params: { id: row.id } })">
            详情
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 16px; justify-content: center"
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="prev, pager, next, total"
      @current-change="handlePageChange"
    />
  </div>
</template>
