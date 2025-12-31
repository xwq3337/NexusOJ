<script setup lang="ts">
import Request from '@/api/axios'
import { onMounted, ref } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import type { Problem } from '@/types/problem'

const problemList = ref<Problem[]>([])
onMounted(async () => {
  try {
    const res = await Request.get<{ info: Problem[] }>({ url: '/problem/list' })
    problemList.value = res.info
  } catch (err) {
    console.error(err)
    ElMessage.error('获取题目列表失败')
  }
})

</script>

<template>
  <el-card>
    <template #header>
      <el-button type="primary" @click="$router.push({ name: 'problem-create' })">创建题目</el-button>
    </template>

    <el-table stripe :data="problemList" style="width: 100%" :scroll-x="true">
      <el-table-column prop="title" label="题目名称" width="150" fixed="left" show-overflow-tooltip />
      <el-table-column label="题目描述" width="100" fixed="left" show-overflow-tooltip>
        <template #default="{ row }">
          {{ row.context }}
        </template>
      </el-table-column>
      <el-table-column prop="difficulty" label="难度" width="60" />
      <el-table-column prop="tags" label="标签" width="150">
        <template #default="{ row }">
          <div>
          <el-tag v-for="tag in row.tags" :key="tag" :type="'success'">{{ tag }}</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="submission" label="提交" width="80" />
      <el-table-column prop="accept" label="通过" width="80" />
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ dayjs(row.created_at).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
      </el-table-column>
      <el-table-column label="最近修改" width="180">
        <template #default="{ row }">
          {{ dayjs(row.updated_at).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
      </el-table-column>
      <el-table-column label="操作" >
        <template #default="{ row }">
          <el-button type="text" @click="() => ElMessage('删除题目')">
            <el-icon>
              <Delete />
            </el-icon>
          </el-button>
          <el-button type="text" @click="$router.push({ name: 'problem-edit', query: { id: row.id } })">
            <el-icon>
              <Edit />
            </el-icon>
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<style scoped>
.mx-1 {
  margin: 0 4px;
}

.problem-detail {
  padding: 20px;
}

.section {
  margin-bottom: 24px;
}

.section h3 {
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
}
</style>
