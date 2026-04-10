<script setup lang="ts">
import { blogApi } from '@nexusoj/server'
import { onMounted, ref } from 'vue'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import { safeJsonParse } from '@nexusoj/utils'
import { ElMessage } from 'element-plus'
import type { Blog } from '@nexusoj/type'
const RecyclingBlogList = ref<Blog[]>([])
onMounted(async () => {
  await blogApi.getRecycleList().then((res) => {
    const { code, info} = res
    if (code === 200 && info){
      RecyclingBlogList.value = info
    }
  }).finally(() => {})

})

const confirm = (id: string) => {
  ElMessage(id)
  // return Request.get({ url: '/blog/delete?id=' + id })
  //   .then(async () => {
  //     await Request.get({ url: "/blog/recycling" }).then((res) => {
  //       RecyclingBlogList.value = res.info
  //     }).finally(() => { })
  //   })
  //   .catch(() => { })
  //   .finally(() => { })
  return new Promise((resolve) => setTimeout(resolve, 2000))
}

const cancel = (e: MouseEvent) => {
  message.error(`'Click on No'${e.clientX}${e.clientY}`)
}
</script>

<template>
  <el-table :data="RecyclingBlogList" width="80%" fixed>
    <el-table-column prop="title" label="标题" align="center"></el-table-column>

    <el-table-column prop="username" label="作者" align="center"></el-table-column>

    <el-table-column prop="tags" label="标签" align="center" width="200">
      <template #default="{ row }">
        <el-tag v-for="(tag, index) in safeJsonParse(row.tags).data" :key="index" effect="dark">
          {{ tag }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="created_at" label="创建时间" align="center" width="200">
      <template #default="{ row }">
        {{ dayjs(row.created_at).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
    </el-table-column>
    <el-table-column label="操作" align="center" width="150">
      <template #default="{ row }">
        <el-button @click="confirm(row.id)">
          <el-icon>
            <Delete />
          </el-icon>
        </el-button>
      </template>
    </el-table-column>
  </el-table>
</template>
