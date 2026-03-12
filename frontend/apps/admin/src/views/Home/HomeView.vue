<script setup lang="ts">
import { blogApi, problemApi, userApi } from '@nexusoj/server'
import { onMounted, ref } from 'vue'
const userCount = ref(0)
const problemCount = ref(0)
const blogCount = ref(0)

onMounted(async () => {
  await userApi.Count().then((res) => {
    const { code, info } = res
    if (code == 200 && info) {
      userCount.value = info
    }
  })
  await problemApi.Count().then((res) => {
    const { code, info } = res
    if (code == 200 && info) {
      problemCount.value = info
    }
  })
    await blogApi.Count().then((res) => {
    const { code, info } = res
    if (code == 200 && info) {
      blogCount.value = info
    }
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
