<script setup lang="ts">
import Request from '@/api/axios';
import { onMounted, ref } from 'vue'
const userCount = ref(0)
const problemCount = ref(0)
const blogCount = ref(0)
onMounted(async () => {
  await Request.get({ url: '/user/count' }).then((res) => {
    userCount.value = res.info
  })
  await Request.get({ url: '/problem/count' }).then((res) => {
    problemCount.value = res.info
  })
  await Request.get({ url: '/blog/count' }).then((res) => {
    blogCount.value = res.info
  })

})

</script>

<template>
  <el-row :gutter="20">
    <el-col :span="6">
      <el-card>
        <el-statistic title="用户总数" :value="userCount" />
        <el-statistic title="题目总数" :value="problemCount" />
        <el-statistic title="博客总数" :value="blogCount" />
      </el-card>
    </el-col>
    <el-col :span="6">
      <el-card>
        <el-statistic title="日活跃用户" :value="1000" />
        <el-statistic title="周活跃用户" :value="10000" />
        <el-statistic title="月活跃用户" :value="100000" />
      </el-card>
    </el-col>
    <el-col :span="6">
      博客排行榜
    </el-col>
    <el-col :span="6">
      用户排行榜
    </el-col>
  </el-row>
</template>

<style scoped></style>
