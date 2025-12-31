<script setup lang="ts">

import Request from '@/api/axios';
import { dayjs, ElMessage } from 'element-plus';
import { onMounted, ref } from 'vue';
const blogList = ref([])
onMounted(async () => {
  await Request.get({ url: '/blog/verify-list' }).then(res => {
    blogList.value = res.info
  }).catch(err => {
    console.log(err)
  })
})
const handleVerify = (index: number) => {
  ElMessage("Verify")
}
const handleReject = (index: number) => {
  ElMessage("Reject")
}
</script>
<template>
  <el-table :data="blogList" style="width: 100%">
    <el-table-column prop="id" label="ID" width="400"></el-table-column>
    <el-table-column prop="title" label="标题" width="200"></el-table-column>
    <el-table-column prop="username" label="作者" width="100"></el-table-column>
    <el-table-column prop="created_at" label="创建时间" width="200">
      <template #default="{row }">{{ dayjs(row.created_at).format('YYYY-MM-DD HH:mm:ss') }}</template>
    </el-table-column>
    <el-table-column label="操作">
      <template #default="scope">
        <el-button type="text" @click="handleVerify(scope.$index)">通过</el-button>
        <el-button type="text" @click="handleReject(scope.$index)">拒绝</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>
